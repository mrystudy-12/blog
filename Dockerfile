# 阶段一：编译
FROM golang:1.21-alpine AS builder
WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY backend ./backend
RUN cd backend && go build -o main cmd/main.go

# 阶段二：运行
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai
COPY --from=builder /app/backend/main .
COPY --from=builder /app/backend/internal/config/config.yaml ./config.yaml
# 创建前端静态资源目录（对应 file_util.go 的路径）
RUN mkdir -p /app/frontend/static/images/avatars
RUN mkdir -p /app/frontend/static/images/articles
EXPOSE 8080
CMD ["./main"]