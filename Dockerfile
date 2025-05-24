FROM golang:1.23-alpine

WORKDIR /app

ARG SERVICE_DIR

COPY . .

WORKDIR /app/${SERVICE_DIR}

RUN go mod tidy

RUN go build -o main .

CMD ["./main"]