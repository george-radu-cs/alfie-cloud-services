import jwt from "jsonwebtoken";
import { GetStringEnvVar } from "@config/env";

const JWT_SECRET = GetStringEnvVar("JWT_SECRET");

const verifyJWT = (token: string) => {
  try {
    const verified = jwt.verify(token, JWT_SECRET);
    if (!verified) {
      throw new Error("JWT verification failed");
    }

    const jwtPayload = jwt.decode(token, { json: true });
    if (!jwtPayload) {
      throw new Error("JWT payload is empty");
    }

    const { email } = jwtPayload;
    return { auth: true, userEmail: email };
  } catch (err) {
    console.log(`JWT verification error: ${err}`);
    return { auth: false, email: "" };
  }
};

export default { verifyJWT };
