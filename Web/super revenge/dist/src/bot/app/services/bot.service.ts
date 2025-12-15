import type { BrowserContext } from "playwright";
import { browserCtx } from "../utils/setup";
import dns from "dns/promises";
import { port } from "../config";
import { decodeJWT, generateJWT, generateToken } from "../utils/jwt";


export async function RunBot(urlRaw: string, user: string | undefined) {
  let context: BrowserContext | undefined;
  let browserType: string = "chromium";
  let ext: boolean = false;
  if (user) {
    const validUser = (await decodeJWT(user)) as {
      visited: string;
      browser: string;
    };
    console.log(`[${new Date().toISOString()}] User: ${JSON.stringify(validUser)}`);
    browserType = validUser.browser;
    
  } else {
    console.log(`[${new Date().toISOString()}] Run chrome`);
    browserType = "chromium";
  }
  try {
    if (browserType === "chromium") {
      ext = true;
    }
    context = await browserCtx(browserType, ext);
    const url = new URL(urlRaw);
    if (url.pathname && url.pathname === "/") {
      return Error("Invalid path, i cant scrap a whole index website lol.");
    }
    if (url.host !== "proxy") {
      return Error("Invalid host, only proxy host is allowed.");
    }
    if (
      [
        "javascript:",
        "data:",
        "blob:",
        "chrome:",
        "chrome-extension:",
        "resource:",
      ].includes(url.protocol)
    ) {
      throw new Error(
        "cannot use this protocol. please use a normal protocol like http or https"
      );
    }

    console.log(`[${new Date().toISOString()}] Running bot...`);

    
    if (browserType !== "chromium") {
      if (
        (url.protocol === "http:" || url.protocol === "https:") &&
        url.hostname !== "localhost" &&
        url.hostname !== "127.0.0.1"
      ) {
        if (url.hostname !== "localhost" && url.hostname !== "127.0.0.1") {
          await dns.lookup(url.hostname);
          return Error("hostname not found");
        }
      }
    } else {
      const token = await generateJWT({
        visited: url.href,
        browser: browserType,
      });      
      const admin_token = await generateToken({
        id: 0,
        email: process.env.ADMIN_EMAIL || "<admin_email>",
        role: "admin",
        name: process.env.ADMIN_NAME || "Admin",
      })
      console.info(`Generated admin token for bot: ${admin_token}`);
      await context.addCookies([
        {
          name: "supernote_session",
          value: admin_token,
          domain: "proxy",
          path: "/",
          httpOnly: true,
          secure: false,
        },
      ]);
      

      await context.addCookies([
        {
          name: "user",
          value: token,
          domain: "localhost",
          path: "/",
          httpOnly: false,
          secure: false,
        },
      ]);
      await context.addCookies([
        {
          name: "visited",
          value: url.href,
          domain: "localhost",
          path: "/",
          httpOnly: false,
          secure: false,
        },
      ]);
      await context.addCookies([
        {
          name: "browser",
          value: browserType,
          domain: "localhost",
          path: "/",
          httpOnly: false,
          secure: false,
        },
      ]);
    }

    const page = await context.newPage();

    page.on("dialog", async (dialog) => {
      console.log("Dialog message:", dialog.message());
    });

    await page.route("**/*", async (route, request) => {
      try {
        const response = await route.fetch();
        const headers = {
          ...response.headers(),
          "content-security-policy":
            "default-src 'self'; script-src 'self'; object-src 'none';",
        };
        await route.fulfill({
          response,
          headers,
          body: await response.body(),
        });
      } catch (err) {
        console.warn(
          `[Route Error] Skipping resource: ${request.url()} - ${err}`
        );
        await route.abort();
      }
    });

    await page.goto(url.href, {
      waitUntil: "networkidle",
    });
    const urlNow = await page.url();
    console.log(`[${new Date().toISOString()}] Bot done. ${urlNow}`);
  } catch (error) {
    return error;
  } finally {
    if (context) {
      await context.close();
    }
  }
}