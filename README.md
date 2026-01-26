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

## Adding Values to Redis

如何在這個 Redis 執行個體（Instance）中，同時進行**設定（Set）與讀取（Get）**數值。

### Setting Values

使用 client.Set 方法來設定數值。這個方法需要傳入鍵（key）、值（value）以及過期時間（expiration）。若將過期時間設定為 0，則代表將該鍵設為永不過期。

實作要點解析：
在 go-redis 套件中，這個方法的定義通常如下：
```go
client.Set(ctx, key, value, expiration)
```
key: 字串型別，作為尋找資料的索引。
value: 您要儲存的內容（可以是字串、數字或序列化後的資料）。
expiration: 使用 time.Duration 型別。例如：0：永久保存。time.Minute * 5：5 分鐘後自動刪除。

### Getting Values

使用 client.Get 方法來讀取。它只需要傳入一個 鍵（key）。由於我們在 Redis 中儲存的是字串，所以我們可以呼叫 .Result() 來同時取得 值（value） 以及可能發生的 錯誤（error）。如果該鍵不存在，Redis 會回傳一個特定的「nil」錯誤，我們在撰寫程式碼時需要針對這種情況進行處理。
