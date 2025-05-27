# E-commerce API Project

A comprehensive e-commerce REST API built with Go, Gin framework, and MongoDB. This project provides complete user authentication, product management, shopping cart functionality, and order processing capabilities.

## 🚀 Features

### Authentication & User Management
- Email verification with OTP
- User registration and login
- JWT-based authentication
- User profile management
- Admin and regular user roles

### Product Management
- Product CRUD operations (Admin only)
- Product listing with pagination
- Product search functionality
- Image URL support
- Product metadata management

### Shopping Cart & Orders
- Add products to cart
- User address management
- Order checkout functionality
- Cart management

## 🛠️ Tech Stack

- **Backend Framework**: Go with Gin
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Email Service**: SendGrid
- **Password Hashing**: bcrypt
- **Environment Management**: godotenv

## 📋 Prerequisites

Before running this application, make sure you have the following installed:

- Go 1.19 or higher
- MongoDB (local or cloud instance)
- SendGrid account (for email verification)

## 🔧 Installation

1. **Clone the repository**
   ```bash
   git clone <your-repo-url>
   cd ecommerce-project
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Create environment file**
   Create a `.env` file in the root directory:
   ```env
   # Server Configuration
   PORT=8080
   API_VERSION=/api/v1

   # Database Configuration
   BD_HOST=localhost:27017
   DATABASE_NAME=ecommerce_db

   # JWT Configuration
   JwtSecrets=your-secret-key-here
   JwtIssuer=ecommerce-api

   # SendGrid Configuration
   SENDGRID_API_KEY=your-sendgrid-api-key
   FROM_EMAIL=noreply@yourdomain.com
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### 1. Verify Email (Send OTP)
```http
POST /ecommerce/verify-email
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### 2. Verify OTP
```http
POST /ecommerce/verify-otp
Content-Type: application/json

{
  "email": "user@example.com",
  "otp": 123456
}
```

#### 3. Register User
```http
POST /ecommerce/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "user@example.com",
  "phone": "+1234567890",
  "password": "securepassword"
}
```

#### 4. Login User
```http
POST /ecommerce/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

### Product Endpoints (Public)

#### 1. List Products
```http
GET /ecommerce-product/products?page=1&limit=10&offset=0
```

#### 2. Search Products
```http
POST /ecommerce-product/search
Content-Type: application/json

{
  "search": "laptop"
}
```

### Product Management (Admin Only)

#### 1. Create Product
```http
POST /ecommerce/products
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "Product Name",
  "description": "Product description",
  "price": 99.99,
  "imageUrl": "https://example.com/image.jpg",
  "metaInfo": {
    "category": "electronics",
    "brand": "BrandName"
  }
}
```

#### 2. Update Product
```http
PUT /ecommerce/products
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "id": "product-id",
  "name": "Updated Product Name",
  "description": "Updated description",
  "price": 149.99
}
```

#### 3. Delete Product
```http
DELETE /ecommerce/products?id=<product-id>
Authorization: Bearer <jwt-token>
```

### User Operations (Authenticated)

#### 1. Add to Cart
```http
POST /ecommerce/cart
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "productId": "product-id"
}
```

#### 2. Add Address
```http
POST /ecommerce/address
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "address1": "123 Main St",
  "city": "New York",
  "country": "USA"
}
```

#### 3. Get User Profile
```http
GET /ecommerce/user/:id
Authorization: Bearer <jwt-token>
```

#### 4. Update User Profile
```http
PUT /ecommerce/user
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "id": "user-id",
  "name": "Updated Name",
  "phone": "+1234567890"
}
```

#### 5. Checkout Order
```http
PUT /ecommerce/checkout
Authorization: Bearer <jwt-token>
```

## 🏗️ Project Structure

```
ecommerce-project/
├── main.go                 # Application entry point
├── .env                    # Environment variables
├── go.mod                  # Go modules file
├── go.sum                  # Go dependencies checksum
├── auth/                   # JWT authentication logic
├── constant/               # Application constants
├── controller/             # HTTP handlers
│   ├── product.go         # Product-related handlers
│   └── user.go           # User-related handlers
├── database/              # Database connection and queries
│   └── manager.go        # Database manager
│   └── connection.go     # Database connection
├── helper/               # Utility functions
├── middleware/           # Custom middleware
├── router/              # Route definitions
│   └── routes.go       # Application routes
|   └── router.go
└── types/              # Data structures and models
```

## 🔐 Authentication Flow

1. **Email Verification**: User provides email → System sends OTP → User verifies OTP
2. **Registration**: After email verification → User registers with details → JWT token issued
3. **Login**: User provides credentials → System validates → JWT token issued
4. **Protected Routes**: JWT token required in Authorization header

## 👥 User Roles

### Regular User
- Register and login
- View and search products
- Manage cart and addresses
- Place orders
- Update profile

### Admin User
- All regular user permissions
- Create, update, delete products
- Manage product inventory

## 🗄️ Database Collections

### Users Collection
- User profile information
- Authentication credentials
- User roles and permissions

### Products Collection
- Product details and pricing
- Product metadata and categories
- Creation and update timestamps

### Verifications Collection
- Email verification records
- OTP management and expiration

### Cart Collection
- User shopping cart items
- Product references
- Checkout status

### Address Collection
- User delivery addresses
- Address validation data

## 🚦 Error Handling

The API uses consistent error response format:

```json
{
  "error": true,
  "message": "Error description"
}
```

Success responses format:
```json
{
  "error": false,
  "message": "success",
  "data": { ... }
}
```

## 🔒 Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- Email verification for registration
- Role-based access control
- CORS middleware configuration
- Input validation and sanitization

## 🧪 Testing

To test the API endpoints, you can use tools like:
- Postman
- cURL
- Thunder Client (VS Code extension)
- Insomnia

## 🚀 Deployment

### Environment Variables for Production
Ensure the following environment variables are set:

```env
PORT=8080
API_VERSION=/api/v1
BD_HOST=your-mongodb-connection-string
JwtSecrets=your-production-jwt-secret
JwtIssuer=your-app-name
SENDGRID_API_KEY=your-sendgrid-key
FROM_EMAIL=your-verified-sender-email
```

### Docker Deployment (Optional)
Create a `Dockerfile`:

```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```


## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Happy Coding! 🎉**