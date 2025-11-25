package user

type UserReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (u *UserReq) ToUser() *User {
	return &User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (ur *User) ToUserRes() *UserRes {
	return &UserRes{
		ID:       ur.ID,
		Username: ur.Username,
		Email:    ur.Email,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	User         *UserRes `json:"user"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
