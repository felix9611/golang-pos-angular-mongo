package dto

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	IpAddress string `form:"ipAddress" json:"ipAddress" binding:"required"`
}
