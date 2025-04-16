FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o textprotocol ./src

FROM alpine

RUN adduser -D appuser

WORKDIR /home/appuser

COPY --from=builder /app/textprotocol .

COPY ./config.json .

RUN chown appuser:appuser textprotocol

EXPOSE 4242

CMD ["./textprotocol", "./config.json"]