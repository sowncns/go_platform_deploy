package dto

type CreateAuthUser struct {
	ID          int64
	Login       string
	Name        string
	Email       string
	AvatarURL   string
	AccessToken string
}
