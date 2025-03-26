package controller

import "github.com/gin-gonic/gin"

func StratServer() {
	//set Realase Mode
	gin.SetMode(gin.ReleaseMode)

	// Load Controller
	router := gin.Default()
	CustomerController(router) // Add this to register the customer routes
	router.Run()
}
