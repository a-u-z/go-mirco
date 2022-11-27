在 production 上的實作
* 實際上，不會是由 broker 發送請求給 mailer
* 應該是 broker 發送登入請求給 authentication ，authentication 驗證失敗後
* 由 authentication 發送請求給 mailer 讓 mailer 發送信件給使用者