import { Router } from "express";
import { authMiddleware } from "@middlewares/auth_middleware";
import ocrController from "@controllers/ocr_controller";

export const ocrRouter = Router();

ocrRouter.post("/image-to-tex", authMiddleware, ocrController.imageToTex);
