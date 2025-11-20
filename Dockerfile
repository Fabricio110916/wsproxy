FROM golang:1.22 AS builder
WORKDIR /app

COPY go.mod ./
COPY main.go ./

RUN go build -o proxy .

FROM gcr.io/distroless/base
WORKDIR /app

COPY --from=builder /app/proxy /app/proxy

EXPOSE 8080

CMD ["/app/proxy"]