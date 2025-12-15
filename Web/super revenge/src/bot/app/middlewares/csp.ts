import type { Request, Response, NextFunction } from "express";

export const setCSPHeader = (req: Request, res: Response, next: NextFunction) => {
  res.setHeader("Content-Security-Policy",
    "style-src 'self' 'unsafe-inline'; " +
    "font-src 'self' https: http:; " +
    "media-src 'self'; " +
    "img-src 'self' data: https: http:; " +
    "object-src 'none'; " +
    "base-uri 'self'; " +
    "form-action 'self';"
  );
  next();
};
