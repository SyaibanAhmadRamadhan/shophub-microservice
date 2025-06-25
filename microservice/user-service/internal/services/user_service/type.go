package userservice

type RegisterInput struct {
	Name        string
	Email       string
	PhoneNumber string
	Password    string
}

type RegisterOutput struct {
	ID int64
}
