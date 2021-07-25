package app

import (
	"github.com/gin-gonic/gin"
	RateInterface "go-kurs-bca/app/exchange_rate"
	KHandler "go-kurs-bca/app/exchange_rate/handler"
)

const (
	SERVICE = "api"
	VERSION = "v1"
)

func KursHTTPRequest(e *gin.Engine, u RateInterface.IRateUsecase) {
	handler := &KHandler.RateHandler{
		RateUsecase: u,
	}

	project := e.Group(SERVICE)
	version := project.Group(VERSION)
	route := version.Group("kurs")

	route.GET("/indexing", handler.Indexing)
	route.GET("/:symbol/", handler.GetRates)
	route.GET("/", handler.GetRates)
	route.POST("/", handler.Add)
	route.PUT("/", handler.Update)
	route.DELETE("/:date", handler.DeleteByDate)
}