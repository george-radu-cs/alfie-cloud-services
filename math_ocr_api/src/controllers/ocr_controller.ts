import { NextFunction, Request, Response } from "express";
import { UploadedFile } from "express-fileupload";
import { StatusCodes } from "http-status-codes";
import ocrService from "@services/ocr_service";

const imageToTex = async (req: Request, res: Response, _: NextFunction) => {
  try {
    if (!req.files) {
      return res
        .status(StatusCodes.BAD_REQUEST)
        .json({ error: true, message: "image not provided" });
    }

    const image = req.files.file as UploadedFile;
    const outputText = await ocrService.convertImageToTex(image.data);

    return res
      .status(StatusCodes.OK)
      .json({ error: false, message: outputText });
  } catch (err) {
    console.error(`error processing image: ${err}`);
    return res
      .status(StatusCodes.BAD_REQUEST)
      .json({ error: true, message: "couldn't process image" });
  }
};

export default { imageToTex };
