FROM golang:1.24-alpine AS builder

RUN apk add --no-cache \
    git \
    openssh \
    tzdata \
    build-base \
    python3 \
    net-tools

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install github.com/buu700/gin@latest

RUN go mod tidy

RUN make build


# runtime
FROM surnet/alpine-wkhtmltopdf:3.20.2-0.12.6-full

RUN apk add --no-cache tzdata curl

WORKDIR /app

COPY --from=builder /app /app

EXPOSE 8002

ENTRYPOINT ["/app/field-service"]