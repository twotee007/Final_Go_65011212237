package controller

import (
	"go-final/dto"
	"go-final/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CustomerController(router *gin.Engine) {
	routes := router.Group("/login")
	{
		routes.GET("/", getAllCutomer)
		routes.POST("/", loginCustomer)
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

	// Compare plain text passwords directly
	if customer.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"details": "Password mismatch",
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

	// Return the converted data
	c.JSON(http.StatusOK, customerModel)
}
