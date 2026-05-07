package main

import (
	"IndulgenceMealPlan/middleware"
	"fmt"
)

func main() {
	token, err := middleware.GenerateToken(123, "jwt颁发服务", "your_jwt_secret_key")
	if err != nil {
		fmt.Println(err)
	}
	claim, err := middleware.ParseToken(token, "your_jwt_secret_key")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(claim)
}
