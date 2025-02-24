package entities

type Account struct {
	ID        int64  `db:"id"`
	Email     string `db:"email"`
	FullName  string `db:"fullname"`
	AvatarURL string `db:"avatar_url"`
}

type AccountPassword struct {
	ID       int64  `db:"id"`
	Password string `db:"password"`
}
