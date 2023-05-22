FROM golang:1.16-alpine

RUN apk add --no-cache tesseract-ocr
RUN addgroup -S mercari && adduser -S trainee -G mercari
# RUN chown -R trainee:mercari /path/to/db

USER trainee

CMD ["go", "version"]