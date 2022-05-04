FROM golang:1.18-alpine

WORKDIR /app/content-parser

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8083 8083

RUN go build -o ./serve content_parser/cmd/serve

CMD ["./serve"]

