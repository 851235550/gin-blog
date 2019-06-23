package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Uid      int    `json: "uid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Select("uid").Where(Auth{Username: username, Password: password}).First(&auth)

	if auth.Uid > 0 {
		return true
	}

	return false
}
