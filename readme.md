一個小型、獨立、鬆散耦合的微服務，這些微服務將相互通信，並開發一個前端應用程序，使用 REST API、RPC、gRPC 以及通過發送和消費 使用 AMQP（高級消息隊列協議）的消息。 包括以下功能：
- 前端服務：僅顯示網頁
- Broker 服務：它是微服務集群的可選單點入口
- Listener 服務：它從 RabbitMQ 接收消息並對其進行操作
- 郵件服務：接受一個 JSON 有效負載，將其轉換為格式化的電子郵件，並將其發送
- Postgres 數據庫：帶有身份驗證服務
- MongoDB 數據庫：帶有日誌記錄服務

上述服務使用 Golang 編寫

下載 minikube https://minikube.sigs.k8s.io/docs/start/
minikube start --nodes=2
如果 minikube 啟動失敗，使用以下指令
minikube delete
minikube status
minikube stop
minikube start
minikube dashboard
