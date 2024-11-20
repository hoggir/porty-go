package models

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DataLoginResponse struct {
	IdUser      string `bson:"idUser" json:"idUser"`
	FullName    string `bson:"fullName" json:"fullName"`
	Email       string `bson:"email" json:"email"`
	SetPassword bool   `bson:"setPassword" json:"setPassword"`
	Token       string `bson:"token" json:"token"`
}

type LoginResponse struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    DataLoginResponse `json:"data"`
}
