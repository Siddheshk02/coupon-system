# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/server
RUN go build -o /server .

# Final stage: use distroless for minimal attack surface
FROM gcr.io/distroless/base-debian12

EXPOSE 8080

COPY --from=builder /server /server

CMD ["/server"]