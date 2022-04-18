FROM golang:1.18-alpine

WORKDIR /app/content-scraper

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8082 8082

RUN go build -o ./serve ContentScraper/cmd/serve

CMD ["./serve"]

