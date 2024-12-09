package controller

import (
	"ecommerce-project/auth"
	"ecommerce-project/constant"
	"ecommerce-project/database"
	"ecommerce-project/helper"
	"ecommerce-project/types"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// VerifyEmail validates an email address and handles OTP generation/expiration
func VerifyEmail(c *gin.Context) {
	var req types.Verification

	// Parse and bind the incoming JSON payload into the 'req' struct
	postBodyErr := c.BindJSON(&req)
	if postBodyErr != nil {
		// Return a 400 error if the request payload is invalid
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": postBodyErr.Error()})
		return
	}

	// Validate if the email field is not empty
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}

	// Fetch the existing OTP record from the database
	resp := database.Mgr.GetSingleRecordByEmail(req.Email, constant.VerificationsCollection)

	if resp.Otp != 0 {
		// If an OTP already exists, check if it is expired
		expirationTime := resp.CreatedAt + constant.OtpValidation
		if expirationTime < time.Now().Unix() {
			// Generate and send a new OTP using a helper function
			req, checkEmail := helper.SendEmailSendGrid(req)
			if checkEmail != nil {
				log.Panicln(checkEmail)
				c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
				return
			}
			// Update the record with the new OTP creation time
			req.CreatedAt = time.Now().Unix()
			database.Mgr.UpdateVerification(req, constant.VerificationsCollection)
			c.JSON(http.StatusOK, gin.H{"error": false, "message": "OTP sent successfully"})
			return
		}
		// Inform the user that the OTP is still valid
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OptAlreadySentError})
		return
	}

	// Generate a new OTP as no prior OTP exists
	req, checkEmail := helper.SendEmailSendGrid(req)
	if checkEmail != nil {
		log.Println(checkEmail)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}

	// Record the OTP creation time and save it in the database
	req.CreatedAt = time.Now().Unix()
	database.Mgr.Insert(req, constant.VerificationsCollection)
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "OTP sent successfully"})
}

// VerifyOtp validates the OTP provided by the user for email verification
func VerifyOtp(c *gin.Context) {
	var req types.Verification

	// Parse and bind the incoming JSON payload into the 'req' struct
	postBodyErr := c.BindJSON(&req)
	if postBodyErr != nil {
		log.Println(postBodyErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": postBodyErr.Error()})
		return
	}

	// Check if email and OTP fields are provided in the request
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailValidationError})
		return
	}
	if req.Otp <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}

	// Fetch the OTP record associated with the given email
	resp := database.Mgr.GetSingleRecordByEmail(req.Email, constant.VerificationsCollection)

	// Check if the email has already been verified
	if resp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AlreadyVerifiedError})
		return
	}

	// Verify the OTP and check if it is expired
	expirationTime := resp.CreatedAt + constant.OtpValidation
	if resp.Otp != req.Otp {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}
	if expirationTime < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpExpiredValidationError})
		return
	}

	// Update the verification record to mark the email as verified
	req.Status = true
	req.CreatedAt = time.Now().Unix()
	err := database.Mgr.UpdateEmailVerifiedStatus(req, constant.VerificationsCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.OtpValidationError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Email verified successfully"})
}

// RegisterUser handles the registration of a new user post-email verification
func RegisterUser(c *gin.Context) {
	var userClient types.UserClient
	var dbUser types.User

	// Parse and bind the incoming JSON payload into the 'userClient' struct
	regErr := c.BindJSON(&userClient)
	if regErr != nil {
		log.Println(regErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": regErr.Error()})
		return
	}

	// Validate the user input fields
	err := helper.CheckUserValidation(userClient)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Check if the email is verified in the verification records
	verificationResp := database.Mgr.GetSingleRecordByEmail(userClient.Email, constant.VerificationsCollection)
	if !verificationResp.Status {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.EmailIsNotVerified})
		return
	}

	// Ensure that the email is not already registered with another user
	userResp := database.Mgr.GetSingleRecordByEmailForUser(userClient.Email, constant.UserCollection)
	if userResp.Email != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AlreadyRegisterWithThisEmail})
		return
	}

	// Prepare the user data for insertion into the database
	dbUser.Email = userClient.Email
	dbUser.Name = userClient.Name
	dbUser.Phone = userClient.Phone
	dbUser.UserType = constant.NormalUser
	dbUser.Password = helper.GenPassHash(userClient.Password) // Hash the user's password
	dbUser.CreatedAt = time.Now().Unix()
	dbUser.UpdatedAt = time.Now().Unix()

	// Insert the new user record into the database
	InsertedID, err := database.Mgr.Insert(dbUser, constant.UserCollection)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Generate a JWT token for the newly registered user
	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48, // Token expiration time in hours
	}
	userID := InsertedID.(primitive.ObjectID)
	token, err := jwtWrapper.GenrateToken(userID, userClient.Email, constant.NormalUser)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Send a success response with the user data and token
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Registration successful", "data": dbUser, "token": token})
}

// UserLogin authenticates a user and generates a JWT token
func UserLogin(c *gin.Context) {
	var loginReq types.Login

	// Parse and bind the incoming JSON payload into the 'loginReq' struct
	err := c.BindJSON(&loginReq)
	if err != nil {
		log.Panicln(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Fetch the user record from the database using the email
	userResp := database.Mgr.GetSingleRecordByEmailForUser(loginReq.Email, constant.UserCollection)
	if userResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return
	}

	// Validate the user's password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(userResp.Password), []byte(loginReq.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.PasswordNotMatchedError})
		return
	}

	// Generate a JWT token for the authenticated user
	jwtWrapper := auth.JwtWrapper{
		SecretKey:      os.Getenv("JwtSecrets"),
		Issuer:         os.Getenv("JwtIssuer"),
		ExpirationTime: 48, // Token expiration time in hours
	}
	token, err := jwtWrapper.GenrateToken(userResp.Id, userResp.Email, userResp.UserType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// Send a success response with the generated token
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Login successful", "token": token})
}

func AddToCart(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotAuthorizedUserError})
		return
	}

	userDBResp := database.Mgr.GetSingleRecordByEmail(email.(string), constant.UserCollection)

	if userDBResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return
	}

	address, err := database.Mgr.GetSingleAddress(userDBResp.ID, constant.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	if address.Address1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.AddressNotExists})
		return
	}
	var cart types.CartClient
	var cartDb types.Cart
	err = c.BindJSON(&cart)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	productId, _ := primitive.ObjectIDFromHex(cart.ProductID)
	// userId, _ := primitive.ObjectIDFromHex(cart.UserId)

	cartDb.ProductID = productId
	cartDb.UserId = userDBResp.ID

	_, err = database.Mgr.Insert(cartDb, constant.CartCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "message": "successful"})
}

func AddAddressOfUser(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotAuthorizedUserError})
		return
	}

	userDBResp := database.Mgr.GetSingleRecordByEmail(email.(string), constant.UserCollection)
	if userDBResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return
	}

	var addressReq types.AddressClient
	err := c.BindJSON(&addressReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	// userId, err := primitive.ObjectIDFromHex(addressReq.UserId)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	// 	return
	// }

	var addressDB types.Address

	addressDB.Address1 = addressReq.Address1
	addressDB.UserId = userDBResp.ID
	addressDB.City = addressReq.City
	addressDB.Country = addressReq.Country

	_, err = database.Mgr.Insert(addressDB, constant.AddressCollection)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "error": false})

}

func GetSingleUser(c *gin.Context) {
	userIdStr := c.Param("id")
	userId, err := primitive.ObjectIDFromHex(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	}

	user, _ := database.Mgr.GetSingleUserByUserId(userId, constant.UserCollection)

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.NotRegisteredUser})
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{"message": "success", "error": true, "data": user})

}

func UpdateUser(c *gin.Context) {
	var userUpdate types.UserUpdateClient
	err := c.BindJSON(&userUpdate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}
	userId, err := primitive.ObjectIDFromHex(userUpdate.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
		return
	}

	userResp, _ := database.Mgr.GetSingleUserByUserId(userId, constant.UserCollection)

	if userResp.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.UserDoesNotExists})
		return
	}

	var user types.User

	user.Id = userId
	user.Name = userResp.Name
	user.Email = userResp.Email
	user.Password = userResp.Password
	user.Phone = userResp.Phone
	user.UserType = userResp.UserType
	user.UpdatedAt = time.Now().Unix()
	user.CreatedAt = userResp.CreatedAt

	if userUpdate.Email != "" {
		user.Email = userUpdate.Email
	}

	if userUpdate.Password != "" {
		user.Password = userUpdate.Password
	}

	if userUpdate.Phone != "" {
		user.Phone = userUpdate.Phone
	}

	if userUpdate.Name != "" {
		user.Name = userUpdate.Name
	}

	err = database.Mgr.UpdateUser(user, constant.UserCollection)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": constant.UserDoesNotExists})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "error": true, "data": user})
}

