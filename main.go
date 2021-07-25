package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	routes "go-kurs-bca/app"
	rateRepo "go-kurs-bca/app/exchange_rate/repository"
	rateUsecase "go-kurs-bca/app/exchange_rate/usecase"
	scrapRepo "go-kurs-bca/app/scrapper/repository"
	"go-kurs-bca/config"
	gorm "go-kurs-bca/db"
	"log"
)

var appConfig = config.Config.App

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	db := gorm.MySQLConn()

	// Repositories
	er := rateRepo.NewRateRepository(db)
	sr := scrapRepo.NewScrapperRepository()

	// Usecases
	eru := rateUsecase.NewRateUsecase(er, sr)

	// Routes
	routes.KursHTTPRequest(r, eru)

	if err := r.Run(fmt.Sprintf(":%s", appConfig.HTTPPort)); err != nil {
		log.Fatal(err)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-MA-Client, X-Platform, X-Api-Key, X-Secret-Key, Accept-Language, X-Product, X-Payment-Token, X-Request-Time")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}