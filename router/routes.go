package router

import (
	"ecommerce-project/constant"
	"ecommerce-project/controller"
	"net/http"
)

var userRoutes = Routes{
	Route{"VerifyEmail", http.MethodPost, constant.VerifyEmailRoute, controller.VerifyEmail},
	Route{"VerifyOtp", http.MethodPost, constant.VerifyOtpRoute, controller.VerifyOtp},
	Route{"Email", http.MethodPost, constant.ResendEmailRoute, controller.VerifyEmail},

	// Resister User
	Route{"RegisterUser", http.MethodPost, constant.UserRegisterRoute, controller.RegisterUser},
	Route{"LoginUser", http.MethodPost, constant.UserLoginRoute, controller.UserLogin},
}

var productGlobalRoutes = Routes{
	Route{"List Product", http.MethodGet, constant.ListProductRoute, controller.ListProductsController},
	Route{"Search Product", http.MethodPost, constant.SearchProductRoute, controller.SearchProduct},
}

var productRoutes = Routes{
	Route{"Register Product", http.MethodPost, constant.RegisterProductRoute, controller.RegisterProduct},
	Route{"Update Product", http.MethodPut, constant.UpdateProductRoute, controller.UpdateProduct },
	Route{"Delete PRoduct", http.MethodDelete, constant.DeleteProductRoute, controller.DeleteProduct},
}


var userAuthRoutes = Routes{
	Route{"Add to cart", http.MethodPost, constant.AddToCartRoute, controller.AddToCart},
	Route{"AddAddress", http.MethodPost, constant.AddAddressRoute, controller.AddAddressOfUser},
	Route{"Get Single User", http.MethodPost, constant.GetSingleUserRoute, controller.GetSingleUser},
	Route{"Update User", http.MethodPut, constant.UpdateUser, controller.UpdateUser},
	Route{"Checkout Order", http.MethodPut, constant.CheckoutRoute, controller.CheckoutOrder},
}