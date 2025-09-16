package modules

type Users struct {
	UserId    int    `json:"user_id,omitempty"`
	Name      string `json:"name,omitempty" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Role      string `json:"role" validate:"required"`
	Address   string `json:"address" validate:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
