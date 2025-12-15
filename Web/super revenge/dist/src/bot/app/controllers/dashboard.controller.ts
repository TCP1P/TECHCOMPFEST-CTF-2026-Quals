import type { Request, Response } from "express";
import { DOMPurify } from "../services/sanitizer.service";
import { username } from "../config";

export const renderDashboard = (req: Request, res: Response) => {
  const rawQuery = req.query.qry || "i promise its not that complex";
  const sanitized = DOMPurify.sanitize(rawQuery.toString());

  console.log(`[MD] ${sanitized}`);

  res.render("dashboard", {
    data: sanitized,
    username
  });
};
