FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/app ./app
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["./app"]
CMD ["serve"]