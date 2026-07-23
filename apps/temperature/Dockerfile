FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o temperature .

FROM alpine:3.23

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/temperature ./temperature

EXPOSE 8081

CMD ["./temperature"]
