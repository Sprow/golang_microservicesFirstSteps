FROM golang:1.18-alpine

WORKDIR /app/content-task

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8081 8081

RUN go build -o ./serve ContentTask/cmd/serve

CMD ["./serve"]