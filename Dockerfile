FROM golang:alpine AS builder
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /tmp/tonflow cmd/main.go

FROM alpine:latest AS tonflow
COPY --from=builder /tmp/tonflow /app/tonflow
# COPY --from=builder /build/assets/logo.png /app/assets/
CMD ["/app/tonflow"]