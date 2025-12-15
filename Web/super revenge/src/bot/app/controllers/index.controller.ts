import type { Request, Response, RequestHandler } from "express";
import { hashIP, unique } from "../utils/hash";
import { RunBot } from "../services/bot.service";

type TurnstileVerifyResponse = {
  success: boolean;
  ["error-codes"]?: string[];
};

export const renderIndex: RequestHandler = async (req, res) => {
  const ip = req.socket.remoteAddress || req.ip;
  if (!ip) {
    res.status(400).json({ message: "Invalid IP" });
    return;
  }

  const hashed = hashIP(ip);
  unique[ip] = hashed;

  res.cookie("hash", hashed);
  res.render("index", {
    csrfToken: req.csrfToken(),
    turnstileSiteKey: process.env.CF_TURNSTILE_SITE_KEY ?? "",
  });
};

// export const getToken: RequestHandler = async (req, res) => {
//   res.json({ csrfToken: req.csrfToken() });
// };

export const startBot: RequestHandler = async (req, res) => {
  const ip = req.socket.remoteAddress || req.ip;
  if (!ip) {
    res.status(400).json({ message: "Invalid IP" });
    return;
  }

  const expected = req.cookies.hash;
  if (unique[ip] !== expected) {
    res.status(403).json({ message: "Forbidden. Visit index first." });
    return;
  }

  unique[ip] = hashIP(ip);
  const url = req.body.url;
  if (!url) {
    res.status(400).json({ message: "Missing URL" });
    return;
  }

  const captchaToken: string | undefined = req.body.turnstileToken;
  const turnstileSecret = process.env.CF_TURNSTILE_SECRET;
  if (!turnstileSecret) {
    res.status(500).json({ message: "Captcha misconfigured." });
    return;
  }
  if (!captchaToken) {
    res.status(400).json({ message: "Captcha required." });
    return;
  }

  const verifyResponse = await fetch(
    "https://challenges.cloudflare.com/turnstile/v0/siteverify",
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        secret: turnstileSecret,
        response: captchaToken,
        remoteip: ip,
      }),
    }
  );

  if (!verifyResponse.ok) {
    res.status(502).json({ message: "Captcha verification failed upstream." });
    return;
  }

  const verifyBody =
    (await verifyResponse.json()) as TurnstileVerifyResponse;
  if (!verifyBody.success) {
    res.status(400).json({
      message:
        "Captcha verification failed." +
        (verifyBody["error-codes"]
          ? ` (${verifyBody["error-codes"].join(", ")})`
          : ""),
    });
    return;
  }

  const result = await RunBot(url, req.cookies.user);
  if (result instanceof Error) {
    res.status(500).json({ message: result.message });
    return;
  }

  res.json({ message: "Crawl complete." });
};
