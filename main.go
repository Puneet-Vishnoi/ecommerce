package main

import (
	"ecommerce-project/constant"
	"ecommerce-project/database"
	"ecommerce-project/helper"
	"ecommerce-project/router"
	"ecommerce-project/types"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// init func runs before main func

func init() {
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading the config from .env file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env config file")
		}
		log.Println("Successfully loaded the config file")
	}
	database.ConnectDb()

	// creating system admin
	hashPassword := helper.GenPassHash("1234")
	user := types.User{
		Name:     "Admin",
		Email:    "admin@gmail.com",
		Password: hashPassword,
		UserType: constant.AdminUser,
	}

	u := database.Mgr.GetSingleRecordByEmailForUser(user.Email, constant.UserCollection)
	
	if u.Email == "" {
		// insertion query to db
		_, err := database.Mgr.Insert(user, constant.UserCollection)

		if err != nil {
			log.Fatal(err)
		}
	}

}

func main() {
	router.ClientRoutes()
}
