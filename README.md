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
- key: 字串型別，作為尋找資料的索引。
- value: 您要儲存的內容（可以是字串、數字或序列化後的資料）。
- expiration: 使用 time.Duration 型別。例如：0：永久保存。time.Minute * 5：5 分鐘後自動刪除。

### Getting Values

使用 client.Get 方法來讀取。它只需要傳入一個 鍵（key）。由於我們在 Redis 中儲存的是字串，所以我們可以呼叫 .Result() 來同時取得 值（value） 以及可能發生的 錯誤（error）。如果該鍵不存在，Redis 會回傳一個特定的「nil」錯誤，我們在撰寫程式碼時需要針對這種情況進行處理。

## Storing Composite Values

如何將自定義的結構（Structs）序列化為 JSON 後存入 Redis

雖然單純在 Redis 中儲存基本的「鍵/值對（key/value pairs）」就能實現很多功能，但有時我們需要更進一步，在資料庫中儲存更複雜的複合式資料結構（composite data structures）。

在這種情況下，我們通常會將這些複合式資料結構**序列化（Marshal）**成 JSON 格式，隨後使用與先前相同的 Set 方法，將這些 JSON 字串存入資料庫中。

因為 Redis 的 SET 指令原生只支援字串（Strings）或二進位資料，所以當我們要儲存 Go 的 struct 時，最常見的做法就是先將其轉換為 JSON 字串。

Redis 的 String 型別其實是 「二進位安全（Binary Safe）」 的，這意味著它不僅能存純文字，也能存下任何經過 json.Marshal 產生的位元組資料（bytes）。

具體的步驟：
1. 定義結構與序列化
首先，需要定義一個結構，並使用 Go 內建的 encoding/json 套件。
範例：
```go
import (
    "encoding/json"
    "fmt"
)

type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
    Age   int    `json:"age"`
}
```

2. 將結構存入 Redis (Marshal)
我們需要先使用 json.Marshal 將結構體轉換為 Byte 切片，然後將其存入 Redis。
範例：
```go
user := User{Name: "Elliot", Email: "elliot@example.com", Age: 25}

// 序列化：將 struct 轉為 JSON
data, err := json.Marshal(user)
if err != nil {
    panic(err)
}

// 存入 Redis
err = rdb.Set(ctx, "user:1", data, 0).Err()
if err != nil {
    panic(err)
}
```

3. 從 Redis 讀取並還原 (Unmarshal) 或是直接回傳 JSON 字串
讀取時，我們先取得 JSON 字串，再使用 json.Unmarshal 將資料**還原（反序列化）**回我們的結構體變數中。
範例：
```go
// 從 Redis 取得資料
val, err := rdb.Get(ctx, "user:1").Result()
if err != nil {
    panic(err)
}

// 反序列化：將 JSON 轉回 struct
var decodedUser User
err = json.Unmarshal([]byte(val), &decodedUser)
if err != nil {
    panic(err)
}

fmt.Printf("讀取到的使用者: %+v\n", decodedUser)
```

💡 為什麼要這樣做？
* 通用性：JSON 是標準格式，即便未來有其他語言（如 Python 或 Node.js）需要讀取這份資料，也能輕鬆解析。
* 效能與彈性：雖然 Redis 有 HSET (Hash) 可以存欄位，但對於複雜或嵌套的物件，序列化成 JSON 通常開發效率更高。

## 延伸問題

Q: 可以不進行反序列化：將 JSON 轉回 struct 嗎 ?

A: 可以，但這取決於你後續要如何「使用」這筆資料。

如果你不進行反序列化（Unmarshal），你拿到的會是一個 JSON 格式的字串（或是 Byte 切片）。以下分兩種情況來說明：

1. 什麼時候「不需要」反序列化？

    如果你只是要把這筆資料直接傳遞給別人，就不需要還原成 Go 的 struct。常見場景包括：

   - 當作 API Proxy：你的 Go 服務從 Redis 讀取資料後，要直接經由 HTTP 回傳給前端（瀏覽器）。既然前端本來就收 JSON，你直接把字串丟出去，效能反而更好，因為省下了 Unmarshal（解碼）再 Marshal（重新編碼）的 CPU 消耗。
   - 純快取傳輸：資料只是暫存在 Redis，讀出來後要傳給另一個微服務處理。

2. 什麼時候「必須」反序列化？

   如果你在程式碼中需要邏輯判斷或修改欄位，就一定要反序列化。例如：

   - 你需要判斷 user.Age 是否大於 18 歲。
   - 你需要修改 user.Email 欄位並存回去。
   - 你需要存取結構體中的特定屬性。

**程式碼對比**

情況 A：直接使用字串（不反序列化）
```go
val, err := rdb.Get(ctx, "user:1").Result()
if err != nil {
    panic(err)
}

// val 此時是 `{"name":"Elliot","email":"...","age":25}`
// 直接印出或回傳給前端
fmt.Println("原始 JSON 字串:", val)
```
情況 B：需要操作資料（必須反序列化）
```go
var user User
err = json.Unmarshal([]byte(val), &user)

// 這樣才能進行邏輯操作
if user.Age >= 18 {
    fmt.Println(user.Name, "已成年")
}
```

💡 效能小建議

在處理大量資料時，反序列化是一個相對「昂貴」的操作（消耗記憶體與 CPU）。如果你發現你的程式只是把 Redis 的資料撈出來直接丟給前端，那麼跳過 Unmarshal 能讓你的 API 反應速度快上一大截！


# 依照 Standard Go Project Layout 的專案結構重構程式碼, 參考：https://github.com/golang-standards/project-layout

採用像 Standard Go Project Layout 這樣的標準結構，可以帶來許多好處：

* 關注點分離 (Separation of Concerns)：每個套件都有明確的職責（例如，database 處理資料庫連線，repository 處理資料存取，domain 定義業務模型）。
* 可維護性：當您需要修改使用者資料的存取邏輯時，您會知道要去 internal/user/repository.go 而不是在一個巨大的 main.go 中尋找。
* 可測試性：將邏輯分離到不同的函數和結構中，可以更容易地為每個部分編寫單元測試。
* 可擴展性：當需要新增功能（例如，產品 product）時，您可以輕鬆地建立新的套件（internal/product），而不會影響現有程式碼。

## 重構後的目錄結構如下：

go-redis-essentials/
├── main.go              # package main
├── database/
│   └── redis.go         # package database（Redis 連線）
└── user/
    ├── user.go          # package user（User struct）
    └── repository.go    # package user（CRUD 操作）
