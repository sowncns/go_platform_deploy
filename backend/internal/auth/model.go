package auth

type GitHubUser struct {
	ID        int64
	Login     string
	Name      string
	Email     string
	AvatarURL string
}