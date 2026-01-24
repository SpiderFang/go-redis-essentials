# Getting Started with Redis in Go

探索如何將 Redis 整合進 Go 應用程式中並加以使用。

## Running Redis With Docker Locally

```bash
$ docker pull redis
$ docker run --name redis-test-instance -p 6379:6379 -d redis
```
第一個 pull 指令的作用是從 DockerHub 取得 Redis 映像檔（image），以便我們後續能透過第二個指令將其作為容器（container）執行。

在第二個指令中，我們指定了 Redis 容器的名稱，並使用 -p 標籤將本地端的 6379 連接埠對應到容器內 Redis 運行的連接埠。

## Connecting to our Redis Instance

1. 匯入廣泛使用的 github.com/go-redis/redis 套件。

2. 定義一個客戶端（Client），其參數包含各項設定選項，例如：我們想要連線的 Redis 執行個體（Instance）位址、密碼，以及在該執行個體中預計使用的資料庫。
在目前的案例中，由於我們本地執行的執行個體並未設定密碼，因此可以將該欄位留白；此外，我們暫時會使用預設的資料庫（以數值 0 表示）。

3. 在定義好這個新的 Redis 客戶端後，接著我們會嘗試執行 Ping 指令來測試執行個體，以確保所有設定都正確無誤，並印出測試結果。
