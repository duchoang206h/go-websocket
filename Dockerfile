FROM golang:1.22.6-alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o app ./cmd/main.go

FROM alpine:3.19 as production

WORKDIR /app

COPY --from=builder /app/app .
#dont's do that
COPY --from=builder /app/app.env .

EXPOSE 8080
CMD ["./app"]