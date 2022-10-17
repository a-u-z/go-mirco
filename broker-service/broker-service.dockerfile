# base go image
FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

# 進入 container 的資料夾
WORKDIR /app

# 打包成二進制檔案 : build -o <檔名> <位置>
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api

RUN chmod +x /app/brokerApp

# build a tiny docker image
FROM alpine:latest

# 創見一個叫做 app 的資料夾
RUN mkdir /app

# 把上面的 builder 的路徑的檔案（/app/brokerApp）複製到這個新的空間的 /app 路徑中
COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]