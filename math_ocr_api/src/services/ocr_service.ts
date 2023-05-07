import FormData from "form-data";
import imageCompressionService from "@services/image_compression_service";
import { GetStringEnvVar } from "@config/env";
import fetch from "node-fetch";
import { Readable } from "stream";

type MathpixOcrResponse = {
  request_id: string;
  version: string;
  image_width: number;
  image_height: number;
  is_printed: boolean;
  is_handwritten: number;
  auto_rotate_confidence: number;
  auto_rotate_degrees: number;
  confidence: number;
  confidence_rate: number;
  latex_styled: string;
  error: string;
  text: string;
};

const MATHPIX_API_URI = GetStringEnvVar("MATHPIX_API_URI");
const MATHPIX_APP_ID = GetStringEnvVar("MATHPIX_APP_ID");
const MATHPIX_APP_KEY = GetStringEnvVar("MATHPIX_APP_KEY");
const MATHPIX_MIN_CONFIDENCE_THRESHOLD = 0.85;

const convertImageToTex = async (imageBuffer: Buffer): Promise<string> => {
  try {
    const compressedImageBuffer =
      await imageCompressionService.compressImageFromBuffer(imageBuffer);

    const outputText = await callMathpixOCRApiToGetTextFromImage(
      compressedImageBuffer
    );

    return outputText;
  } catch (err) {
    throw new Error(`error getting text from image: ${err}`);
  }
};

const callMathpixOCRApiToGetTextFromImage = async (
  imageBuffer: Buffer
): Promise<string> => {
  try {
    const formData = new FormData();
    formData.append("file", Readable.from(imageBuffer), {
      filename: "image_buffer.jpeg",
      contentType: "image/jpeg",
    });
    formData.append(
      "options_json",
      JSON.stringify({
        math_inline_delimiters: ["$$ ", " $$"],
        rm_spaces: true,
        confidence_threshold: MATHPIX_MIN_CONFIDENCE_THRESHOLD,
      })
    );

    const response = await fetch(MATHPIX_API_URI, {
      method: "POST",
      headers: {
        app_id: MATHPIX_APP_ID,
        app_key: MATHPIX_APP_KEY,
        contentType: "multipart/form-data",
      },
      body: formData,
    });

    const responseJson: MathpixOcrResponse = await response.json();
    if (responseJson.error) {
      throw new Error(`error calling Mathpix OCR API: ${responseJson.error}`);
    }

    return responseJson.text;
  } catch (err) {
    throw new Error(`error calling Mathpix OCR API: ${err}`);
  }
};

export default { convertImageToTex };
