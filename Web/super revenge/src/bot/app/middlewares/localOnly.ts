import type { Request, Response, NextFunction } from "express";

export const localOnly = (req: Request, res: Response, next: NextFunction) => {
  const ip = req.socket.remoteAddress;
  if (
    ip?.includes("127.0.0.1") ||
    ip?.includes("::1") ||
    ip?.includes("172.22.0.1")
  ) {
    return next();
  }

  console.warn(`[IP] ${ip} tried accessing dashboard`);
  res.status(400).json({ message: "Invalid IP" });
};
