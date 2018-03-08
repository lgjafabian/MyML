package app

import (
	mymlController "github.com/lgjafabian/MyML/src/api/controllers/myml"
    pingController "github.com/lgjafabian/MyML/src/api/controllers/ping"
    "github.com/gin-gonic/gin"
)

func StartApp() {
    router := gin.Default()
	router.GET("/ping", pingController.Ping)
	router.GET("/myml/order/:orderID", mymlController.GetInformation)
	router.Run(":8080")
}