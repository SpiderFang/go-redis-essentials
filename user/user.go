package user

// User 定義了應用程式中的使用者模型
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
