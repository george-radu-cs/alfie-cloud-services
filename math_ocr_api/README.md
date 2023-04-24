# Math OCR API

specify only format `text`, html doesn't render

receive the json response, take the "text" key and eliminate double escapes `\\` to `\`
for text return it and for markdown replace the start `\(` and end `\)` with `$$` and return it

<!-- curl -X POST https://api.mathpix.com/v3/text \
-H 'app_id: APP_ID' \
-H 'app_key: APP_KEY' \
--form 'file=@"cases_hw.jpg"' \
--form 'options_json="{\"math_inline_delimiters\": [\"$\", \"$\"], \"rm_spaces\": true}"' -->
<!-- here can specify delimiters to be$$ -->

<!-- Try to use images under 100KB for maximum speeds.  -->
<!-- https://mathpix.com/docs/ocr/best-practices -->
