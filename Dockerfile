FROM golang:1.22 AS builder
WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
COPY . .
RUN go build -o /app/main

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
EXPOSE 80