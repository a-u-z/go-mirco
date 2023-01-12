1. 設置 env
   1. 在 .env 檔案中加入
   2. 在 docker compose 檔案中加入，在該服務名稱的下一階。
      1. environment: ENV_TEST: ${ENV_TEST}
   3. 將此服務 down(stop, rm) 再 up ，讓他吃到新的設定