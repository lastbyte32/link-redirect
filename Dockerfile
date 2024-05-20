FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod  ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o server

FROM alpine:latest
COPY --from=builder /app/server /app/server
CMD ["/app/server"]

EXPOSE 80