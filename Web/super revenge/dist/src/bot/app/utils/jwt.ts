import crypto from "crypto";
import { secret } from "../config";

function validateKey(key: Buffer): Buffer {
  const mask = Buffer.from("my-static-mask-32-bytes", "utf-8");
  return Buffer.from(key.map((b, i) => b ^ mask[i % mask.length]));
}

// // GenerateToken creates a new JWT token for a user
// func GenerateToken(user models.User) (string, error) {
// 	// Get JWT secret from config
// 	jwtSecret := config.GetConfig().JWTSecret
// 	if jwtSecret == "" {
// 		return "", errors.New("JWT secret is not configured")
// 	}

// 	// Set token expiration time (e.g., 24 hours)
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	// Create claims with user information
// 	claims := &JWTClaims{
// 		UserID: user.ID,
// 		Email:  user.Email,
// 		Role:   user.Role.Name,
// 		Name:   user.Name,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 			IssuedAt:  jwt.NewNumericDate(time.Now()),
// 			NotBefore: jwt.NewNumericDate(time.Now()),
// 			Issuer:    "sr-api",
// 			Subject:   fmt.Sprintf("%d", user.ID),
// 		},
// 	}

// 	// Create token with claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Sign the token with the secret key
// 	tokenString, err := token.SignedString([]byte(jwtSecret))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

export async function generateToken(user: {
  id: number | string;
  email: string;
  role: string;
  name: string;
}): Promise<string> {
  const jwtSecret = process.env.JWT_SECRET;
  if (!jwtSecret) {
    throw new Error("JWT secret is not configured");
  }
  if (!user) {
    throw new Error("User payload is required to generate a token");
  }

  const now = Math.floor(Date.now() / 1000);
  const header = { alg: "HS256", typ: "JWT" };
  const payload = {
    userId: user.id,
    email: user.email,
    role: user.role,
    name: user.name,
    exp: now + 24 * 60 * 60,
    iat: now,
    nbf: now,
    iss: "sr-api",
    sub: `${user.id}`,
  };

  const encode = (obj: object) =>
    Buffer.from(JSON.stringify(obj)).toString("base64url");
  const data = `${encode(header)}.${encode(payload)}`;
  const signature = crypto
    .createHmac("sha256", jwtSecret)
    .update(data)
    .digest("base64url");

  return `${data}.${signature}`;
}

export async function generateJWT(data: object): Promise<string> {
  const key = crypto.createHash("sha256").update(secret).digest();
  const iv = crypto.randomBytes(12);
  const cipher = crypto.createCipheriv("aes-256-gcm", key, iv);
  const payload = Buffer.from(JSON.stringify(data), "utf8");
  const encrypted = Buffer.concat([cipher.update(payload), cipher.final()]);
  const tag = cipher.getAuthTag();
  const obfKey = validateKey(key);
  const token = Buffer.concat([iv, encrypted, tag, obfKey]);
  return token.toString("base64url");
}

export async function decodeJWT(token: string): Promise<object> {
  const buf = Buffer.from(token, "base64url");
  const iv = buf.subarray(0, 12);
  const tag = buf.subarray(buf.length - 48, buf.length - 32);
  const obfKey = buf.subarray(buf.length - 32);
  const ciphertext = buf.subarray(12, buf.length - 48);
  const key = validateKey(obfKey);
  const expectedKey = crypto.createHash("sha256").update(secret).digest();
  if (!key.equals(expectedKey)) {
    console.log("expected:", expectedKey.toString("hex"));
    console.log("got     :", key.toString("hex"));
    throw new Error("Invalid key");
  }
  const decipher = crypto.createDecipheriv("aes-256-gcm", key, iv);
  decipher.setAuthTag(tag);
  const decrypted = Buffer.concat([
    decipher.update(ciphertext),
    decipher.final(),
  ]);
  return JSON.parse(decrypted.toString("utf8"));
}
