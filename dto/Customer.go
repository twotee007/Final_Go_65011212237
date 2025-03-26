package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateAddressRequest struct {
	NewAddress string `json:"new_address" binding:"required"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
type SearchPro struct {
	ProductName string `form:"product_name" binding:"omitempty"`
	MinPrice    string `form:"min_price" binding:"omitempty"`
	MaxPrice    string `form:"max_price" binding:"omitempty"`
}

type AddToCartRequest struct {
	CustomerID int    `json:"customer_id" binding:"required"`
	ProductID  int    `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
	CartName   string `json:"cart_name" binding:"required"`
}
