package main

import (
	"go-final/controller"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	dsn := viper.GetString("mysql.dsn")
	dialactor := mysql.Open(dsn)

	db, err := gorm.Open(dialactor)
	if err != nil {
		panic(err)
	}
	println("Connection Success")

	// Set the DB connection in the controller package
	controller.DB = db
	controller.StratServer()
}
