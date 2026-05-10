FROM docker.io/library/golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o api-binary .

FROM docker.io/library/alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/
RUN mkdir -p /root/data

COPY --from=builder /app/api-binary .

EXPOSE 8080

CMD ["./api-binary"]
