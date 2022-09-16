package auth

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type SigninResponse struct {
	Token string `json:"token"`
	User  UserResponse
}

type UserResponse struct {
	ID       int64  `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
	Role     string `json:"role"`
}

type UserID struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type PasswordUpdate struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
	User        string `json:"user" swaggerignore:"true"`
	UserID      int64  `json:"user_id" swaggerignore:"true"`
}
