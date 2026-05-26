FROM golang:1.22-alpine AS builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOFLAGS=-mod=mod go build -ldflags="-s -w" -o server ./cmd/server/

FROM alpine:latest
RUN apk add --no-cache python3 bash ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]