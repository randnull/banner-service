package models


type Token struct {
	Token string `json:"token"`
}

type Register struct {
	Username 	string `json:"username"`
	Password  	string `json:"password"`
}


type User struct {
	ID			int		`db:"id"`
	Username 	string	`db:"username"`
	Password  	string	`db:"password"`
	IsAdmin		bool	`db:"is_admin"`
}
