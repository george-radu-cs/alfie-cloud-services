import { Request, Response, NextFunction } from "express";
import { StatusCodes } from "http-status-codes";
import jwtService from "@services/jwt_service";

export const authMiddleware = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  try {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      return res
        .status(StatusCodes.UNAUTHORIZED)
        .json({ error: true, message: "No token provided" });
    }

    const parts = authHeader.split(" ");
    if (parts.length !== 2) {
      return res
        .status(StatusCodes.UNAUTHORIZED)
        .json({ error: true, message: "Token error" });
    }

    const [scheme, token] = parts;
    if (!/^Bearer$/i.test(scheme)) {
      return res
        .status(StatusCodes.UNAUTHORIZED)
        .json({ error: true, message: "Token malformatted" });
    }

    const authResult = jwtService.verifyJWT(token);
    if (!authResult.auth) {
      return res
        .status(StatusCodes.UNAUTHORIZED)
        .json({ error: true, message: "Invalid token" });
    }

    res.locals.userEmail = authResult.userEmail;
    return next();
  } catch (err) {
    return next(err);
  }
};
