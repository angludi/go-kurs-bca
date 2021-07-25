package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	RateHandler "go-kurs-bca/app/exchange_rate/handler"
	RateRepository "go-kurs-bca/app/exchange_rate/repository"
	RateUsecase "go-kurs-bca/app/exchange_rate/usecase"
	ScrapperRepository "go-kurs-bca/app/scrapper/repository"
	mysql "go-kurs-bca/db"
	"go-kurs-bca/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	handler *RateHandler.RateHandler
	router  *gin.Engine
)

func init() {
	gin.SetMode(gin.TestMode)

	db := mysql.MySQLConn()
	sr := ScrapperRepository.NewScrapperRepository()
	rr := RateRepository.NewRateRepository(db)
	ru := RateUsecase.NewRateUsecase(rr,sr)

	handler = &RateHandler.RateHandler{
		RateUsecase: ru,
	}

	router = gin.Default()
	router.POST("/api/v1/kurs/indexing", handler.Indexing)
	router.GET("/api/v1/kurs/:symbol", handler.GetRates)
	router.GET("/api/v1/kurs", handler.GetRates)
	router.POST("/api/v1/kurs", handler.Add)
	router.PUT("/api/v1/kurs", handler.Update)
	router.DELETE("/api/v1/kurs/:date", handler.DeleteByDate)
}

func TestAPIIndexing(t *testing.T) {
	w := httptest.NewRecorder()
	//c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("POST", "/api/v1/kurs/indexing", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t,201, w.Code)
}

func TestAPIGetRatesByDate(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/v1/kurs", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	q := req.URL.Query()
	q.Add("startdate", "2021-07-23")
	q.Add("enddate", "2021-07-23")
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	assert.Equal(t,200, w.Code)
}

func TestAPIGetRatesBySymbol(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/api/v1/kurs/USD", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	q := req.URL.Query()
	q.Add("startdate", "2021-07-23")
	q.Add("enddate", "2021-07-23")
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)

	assert.Equal(t,200, w.Code)
}

func TestAPIDeleteRatesByDate(t *testing.T) {
	w := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", "/api/v1/kurs/2021-07-21", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(w, req)

	assert.Equal(t,200, w.Code)
}

func TestAPIInsertRates(t *testing.T) {
	var params models.RateRequest

	w := httptest.NewRecorder()
	//c, _ := gin.CreateTestContext(w)

	params.Date = "2021-07-21"
	params.Symbol = "USD"
	params.ERate.Beli = 14001
	params.ERate.Jual = 14002
	params.TTCounter.Jual = 14004
	params.TTCounter.Beli = 14003
	params.BankNotes.Beli = 14005
	params.BankNotes.Jual = 14006

	paramsJSON, err := json.Marshal(params)

	req, err := http.NewRequest("POST", "/api/v1/kurs", bytes.NewBuffer(paramsJSON))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(w, req)

	b, _ := ioutil.ReadAll(w.Body)
	fmt.Println("RESULT:",string(b))

	assert.Equal(t,200, w.Code)
}

func TestAPIUpdateRates(t *testing.T) {
	var params models.RateRequest

	w := httptest.NewRecorder()
	//c, _ := gin.CreateTestContext(w)

	params.Date = "2021-07-21"
	params.Symbol = "USD"
	params.ERate.Beli = 14011
	params.ERate.Jual = 14012
	params.TTCounter.Jual = 14004
	params.TTCounter.Beli = 14003
	params.BankNotes.Beli = 14005
	params.BankNotes.Jual = 14006

	paramsJSON, err := json.Marshal(params)

	req, err := http.NewRequest("PUT", "/api/v1/kurs", bytes.NewBuffer(paramsJSON))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	router.ServeHTTP(w, req)

	b, _ := ioutil.ReadAll(w.Body)
	fmt.Println("RESULT:",string(b))

	assert.Equal(t,200, w.Code)
}