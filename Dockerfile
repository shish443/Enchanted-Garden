#Enchanted-Garden/Dockerfile
FROM golang:1.26-alpine AS builder
 
WORKDIR /build
 
COPY go.mod go.sum ./
RUN go mod download
 
COPY . .
 
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o garden ./cmd
 
FROM alpine:3.21
 
RUN apk add --no-cache ca-certificates tzdata
 
WORKDIR /app
 
COPY --from=builder /build/garden .
COPY --from=builder /build/db ./db
 
EXPOSE 8080
 
CMD ["./garden"]