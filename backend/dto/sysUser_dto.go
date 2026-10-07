package dto

import "golang-fixedasset-mongo-backend/backend/models"

type SysUserPageDto struct {
	RoleIds []string `json:"roleIds"`
	Name    string   `json:"name"`
	Page    int64    `json:"page"`
	Limit   int64    `json:"limit"`
}

type SysUserList struct {
	Lists []models.SysUsers `json:"lists"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

type SysUserUpdatePasswordDto struct {
	Username    string `json:"username"`
	NewPassword string `json:"newPassword"`
}

type SysUserAvatarUpdateDto struct {
	Username  string `json:"username"`
	PhotoBase string `json:"photo"`
}

type SysUserCreateResponseDto struct {
	ID string `json:"id"`
}
