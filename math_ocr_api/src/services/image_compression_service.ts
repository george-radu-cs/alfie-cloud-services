import sharp from "sharp";

const MAX_COMPRESSED_IMAGE_BYTE_SIZE = 100000;
const START_COMPRESSION_QUALITY = 95;
const MIN_COMPRESSION_QUALITY = 5;

/**
 * - compress image until obtaining a size less or equal than target size
 * or until compression quality is less than the minimum quality
 * - if the image provided is already smaller than the target size
 * return the image, since no compression is needed
 * -will throw an error if the image cannot be compressed
 *
 * @param {Buffer} imageBuffer - image buffer to compress
 * @returns {Buffer} - compressed image buffer
 */
const compressImageFromBuffer = async (
  imageBuffer: Buffer
): Promise<Buffer> => {
  try {
    // if image is already smaller than target size, return image, no need to compress
    if (imageBuffer.byteLength < MAX_COMPRESSED_IMAGE_BYTE_SIZE) {
      return imageBuffer;
    }

    let compressionQuality = START_COMPRESSION_QUALITY;
    let outputImageBuffer;
    while (compressionQuality >= MIN_COMPRESSION_QUALITY) {
      outputImageBuffer = await sharp(imageBuffer)
        .withMetadata()
        .jpeg({ quality: compressionQuality })
        .toBuffer();

      if (outputImageBuffer.byteLength <= MAX_COMPRESSED_IMAGE_BYTE_SIZE) {
        return outputImageBuffer;
      }

      compressionQuality -= 5;
    }

    if (!outputImageBuffer) {
      throw new Error("image not compressed");
    }

    return outputImageBuffer;
  } catch (err) {
    throw new Error(`Image compression failed: ${err}`);
  }
};

export default { compressImageFromBuffer };
