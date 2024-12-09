package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JwtWrapper struct {
	SecretKey       string
	Issuer         string
	ExpirationTime int64
}

type JwtClaim struct {
	UserId   primitive.ObjectID
	Email    string
	UserType string
	jwt.StandardClaims
}


// GenrateToken generates a JWT token based on user data and signs it with the secret key
func (j *JwtWrapper) GenrateToken(id primitive.ObjectID, email, userType string) (token string, err error) {
	// Define claims struct to hold user-specific data and standard JWT claims
	claims := &JwtClaim{
			UserId:   id,         // Set the user's unique ID
			UserType: userType,   // Set the user's type (e.g., admin, regular user)
			Email:    email,      // Set the user's email
			StandardClaims: jwt.StandardClaims{
					ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.ExpirationTime)).Unix(), // Set expiration time (in hours, based on ExpirationTime)
					Issuer:    j.Issuer, // Set the issuer of the token (from JwtWrapper)
			},
	}

	// Create a new token with the specified signing method and claims
	token1 := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key and return the token string
	token, err = token1.SignedString([]byte(j.SecretKey)) 
	if err != nil {
			// If there is an error while signing, return an empty string and the error
			return "", err
	}

	// Return the signed token and a nil error
	return token, nil
}



// ValidateToken validates a signed JWT token and extracts the claims
func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaim, err error) {
	// Parse the signed token with claims of type JwtClaim
	token, err := jwt.ParseWithClaims(
			signedToken,
			&JwtClaim{}, // Empty JwtClaim struct, the real claims will be filled here
			func(token *jwt.Token) (interface{}, error) {
					// Provide the secret key for verification using the JwtWrapper's SecretKey
					return []byte(j.SecretKey), nil
			},
	)
	
	// If there was an error parsing the token, return the error
	if err != nil {
			return
	}

	// Attempt to type assert the token claims into the expected JwtClaim type
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
			// If the claims could not be parsed into JwtClaim, return an error
			return nil, errors.New("could not parse claims")
	}

	// Check if the token has expired by comparing the expiration time (ExpiresAt) with the current time
	if claims.ExpiresAt < time.Now().Local().Unix() {
			// If the token is expired, return an error
			return nil, errors.New("token is expired")
	}

	// If everything is valid, return the claims and nil error
	return
}


func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			// If no token is provided, return Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Authorization token is required"})
			c.Abort()
			return
		}

		// Split token to extract the bearer token
		extractedToken := strings.Split(token, "Bearer")
		if len(extractedToken) != 2 {
			// If token format is incorrect, return Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token format"})
			c.Abort()
			return
		}

		// Trim spaces and extract the actual token
		clientToken := strings.TrimSpace(extractedToken[1])

		// Create a JWT wrapper instance with your secret and issuer
		jwtWrapper := JwtWrapper{
			SecretKey: os.Getenv("JwtSecrets"),
			Issuer:    os.Getenv("JwtIssuer"),
		}

		// Validate the token using your JWT wrapper
		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
			// If token validation fails, return Unauthorized
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set claims data into context for further processing
		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)

		// Call the next handler in the chain
		c.Next()
	}
}
