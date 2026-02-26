FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY *.go ./
RUN go build -o simple-echo .

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/simple-echo .
EXPOSE 8080
CMD ["./simple-echo"]
