import crypto from "crypto";

export const unique: Record<string, string> = {};

export function hashIP(ip: string) {
  return crypto.createHash("sha256")
    .update(ip + crypto.randomBytes(32).toString("hex"))
    .digest("hex");
}
