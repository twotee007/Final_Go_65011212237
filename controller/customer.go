package controller

import (
	"go-final/dto"
	"go-final/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CustomerController(router *gin.Engine) {
	routes := router.Group("/login")
	{
		routes.GET("/", getAllCutomer)
		routes.POST("/", loginCustomer)
	}
	customerRoutes := router.Group("/customer")
	{
		customerRoutes.PUT("/:id/address", updateAddress)
		customerRoutes.PUT("/:id/password", changePassword)
		customerRoutes.GET("/:id/carts", getCustomerCarts) // New endpoint
	}
	productRoutes := router.Group("/products")
	{
		productRoutes.GET("/search", searchProducts)  // Search products
		productRoutes.POST("/cart", addProductToCart) // Add product to cart
	}
}

func getAllCutomer(c *gin.Context) {
	var customers []model.Customer
	result := DB.Find(&customers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func loginCustomer(c *gin.Context) {
	user := dto.LoginRequest{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	var customer model.Customer
	err := DB.Where("email = ?", user.Email).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"details": "Email not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"details": err.Error(),
			})
		}
		return
	}

	// Check if the stored password is a bcrypt hash (starts with $2a$ or $2b$)
	if len(customer.Password) > 4 && (customer.Password[:4] == "$2a$" || customer.Password[:4] == "$2b$") {
		// Verify hashed password with bcrypt
		if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(user.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"details": "Password mismatch",
			})
			return
		}
	} else {
		// Compare plain text passwords directly (for backward compatibility)
		if customer.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"details": "Password mismatch",
			})
			return
		}
	}

	// Convert model.Customer to model.CustomerModel
	customerModel := model.CustomerModel{
		CustomerID:  int64(customer.CustomerID),
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		Password:    customer.Password, // จะไม่แสดงใน JSON เพราะมี `json:"-"`
		CreatedAt:   customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// Return the converted data
	c.JSON(http.StatusOK, customerModel)
}
func updateAddress(c *gin.Context) {
	// Parse the request body to get the new address
	var request dto.UpdateAddressRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Get the customer ID from the URL parameters
	customerID := c.Param("id")

	// Find the customer by ID
	var customer model.Customer
	err := DB.Where("customer_id = ?", customerID).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Customer not found",
				"details": "No customer found with the given ID",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"details": err.Error(),
			})
		}
		return
	}

	// Update the customer's address
	customer.Address = request.NewAddress

	// Save the updated customer record
	if err := DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update address",
			"details": err.Error(),
		})
		return
	}

	// Convert model.Customer to model.CustomerModel
	customerModel := model.CustomerModel{
		CustomerID:  int64(customer.CustomerID),
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Address:     customer.Address,
		Password:    customer.Password,
		CreatedAt:   customer.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   customer.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	// Respond with the updated customer information
	c.JSON(http.StatusOK, gin.H{
		"message":  "Address updated successfully",
		"customer": customerModel,
	})
}
func changePassword(c *gin.Context) {

	var request dto.PasswordChangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	customerID := c.Param("id")
	var customer model.Customer
	err := DB.Where("customer_id = ?", customerID).First(&customer).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Customer not found",
				"details": "No customer found with the given ID",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"details": err.Error(),
			})
		}
		return
	}

	// Verify old password
	if len(customer.Password) > 4 && (customer.Password[:4] == "$2a$" || customer.Password[:4] == "$2b$") {
		if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(request.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid old password",
				"details": "Old password does not match",
			})
			return
		}
	} else {
		if customer.Password != request.OldPassword {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid old password",
				"details": "Old password does not match",
			})
			return
		}
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to encrypt password",
			"details": err.Error(),
		})
		return
	}

	// Update the password
	customer.Password = string(hashedPassword)
	if err := DB.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update password",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password updated successfully",
	})
}
func searchProducts(c *gin.Context) {
	var request dto.SearchPro
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters", "details": err.Error()})
		return
	}

	var products []model.Product
	query := DB.Model(&model.Product{})

	if request.ProductName != "" {
		query = query.Where("product_name LIKE ?", "%"+request.ProductName+"%")
	}
	if request.MinPrice != "" {
		query = query.Where("price >= ?", request.MinPrice)
	}
	if request.MaxPrice != "" {
		query = query.Where("price <= ?", request.MaxPrice)
	}

	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "products": products})
}
func addProductToCart(c *gin.Context) {
	var request dto.AddToCartRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}
	// Check if customer exists
	var customer model.Customer
	if err := DB.Where("customer_id = ?", request.CustomerID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Check if product exists and has enough stock
	var product model.Product
	if err := DB.Where("product_id = ?", request.ProductID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}
	if product.StockQuantity < request.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
		return
	}

	// Find or create cart
	var cart model.Cart
	err := DB.Where("customer_id = ? AND cart_name = ?", request.CustomerID, request.CartName).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create new cart
			cart = model.Cart{
				CustomerID: request.CustomerID,
				CartName:   request.CartName,
			}
			if err := DB.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart", "details": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
	}

	// Check if product already exists in cart
	var cartItem model.CartItem
	err = DB.Where("cart_id = ? AND product_id = ?", cart.CartID, request.ProductID).First(&cartItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Add new cart item
			cartItem = model.CartItem{
				CartID:    cart.CartID,
				ProductID: request.ProductID,
				Quantity:  request.Quantity,
			}
			if err := DB.Create(&cartItem).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart", "details": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
			return
		}
	} else {
		// Update existing cart item quantity
		newQuantity := cartItem.Quantity + request.Quantity
		if product.StockQuantity < newQuantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for updated quantity"})
			return
		}
		cartItem.Quantity = newQuantity
		if err := DB.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully", "cart_id": cart.CartID})
}
func getCustomerCarts(c *gin.Context) {
	customerIDStr := c.Param("id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID", "details": err.Error()})
		return
	}

	// Verify customer exists
	var customer model.Customer
	if err := DB.Where("customer_id = ?", customerID).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		}
		return
	}

	// Fetch all carts with preloaded items and products
	var carts []model.Cart
	if err := DB.Where("customer_id = ?", customerID).
		Preload("CartItems.Product").
		Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Convert to response format
	var cartResponses []dto.CartResponse
	for _, cart := range carts {
		cartResponse := dto.CartResponse{
			CartID:    cart.CartID,
			CartName:  cart.CartName,
			CreatedAt: cart.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: cart.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		for _, item := range cart.CartItems {
			cartItemResponse := dto.CartItemResponse{
				CartItemID:  item.CartItemID,
				ProductID:   item.ProductID,
				ProductName: item.Product.ProductName,
				Description: item.Product.Description,
				Price:       item.Product.Price,
				Quantity:    item.Quantity,
			}
			cartResponse.Items = append(cartResponse.Items, cartItemResponse)
		}

		cartResponses = append(cartResponses, cartResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"carts":  cartResponses,
	})
}
