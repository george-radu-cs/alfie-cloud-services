import { GetIntEnvVar } from "@config/env";
import express, { json } from "express";
import rateLimit from "express-rate-limit";
import cors from "cors";
import http from "http";
import fileUpload from "express-fileupload";
import { ocrRouter } from "@routes/ocr_router";

const PORT = GetIntEnvVar("PORT");
const app = express();
const server = http.createServer(app);
app.use(json({ strict: true }));
app.use(cors());
app.use(fileUpload());

const rateLimiter = rateLimit({
  windowMs: GetIntEnvVar("RATE_LIMITER_WINDOW_IN_MS"),
  max: GetIntEnvVar("RATE_LIMITER_WINDOW_MAX_REQUESTS"),
});
app.use(rateLimiter);

app.use("/ocr", ocrRouter);

(async () => {
  try {
    server.listen(PORT, () => {
      console.log(`Server started on port ${PORT}`);
    });
  } catch (err) {
    console.error(err);
  }
})();
