FROM alpine

RUN apk add --no-cache tesseract-ocr
RUN apk add --no-cache go

RUN addgroup -S mercari && adduser -S trainee -G mercari

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

USER trainee

CMD ["go", "run", "main.go"]
