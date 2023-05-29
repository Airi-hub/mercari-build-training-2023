FROM alpine

RUN apk add --no-cache tesseract-ocr
RUN apk add build-base 
RUN apk add --no-cache go
RUN addgroup -S mercari && adduser -S trainee -G mercari



WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN chown -R trainee:mercari images items.json

USER trainee

RUN go env -w GOARCH=arm64

CMD ["go", "run", "app/main.go"]
