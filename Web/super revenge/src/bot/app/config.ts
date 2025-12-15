import crypto from "crypto";

export const port = process.env.PORT || 1337;
export const secret = process.env.SECRET || crypto.randomBytes(32).toString("hex");
export const username = process.env.USERNAME || "cowboybeboplover";

console.log(`[${new Date().toISOString()}] Secret: ${secret}`);
