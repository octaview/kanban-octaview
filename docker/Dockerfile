# Dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /kanban-backend ./cmd/server

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
COPY --from=builder /kanban-backend /kanban-backend
EXPOSE 8080
ENTRYPOINT ["/kanban-backend"]
