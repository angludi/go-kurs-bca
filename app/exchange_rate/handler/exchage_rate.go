package handler

import (
	"github.com/gin-gonic/gin"
	"go-kurs-bca/app/base"
	RateInterface "go-kurs-bca/app/exchange_rate"
	"go-kurs-bca/models"
	"go-kurs-bca/transformers"
	"net/http"
	"time"
)

type RateHandler struct {
	RateUsecase RateInterface.IRateUsecase
}

func (h *RateHandler) Indexing(c *gin.Context) {
	err := h.RateUsecase.Indexing()
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	base.RespondCreated(c, "Success")
	return
}

func (h *RateHandler) GetRates(c *gin.Context) {
	symbol := c.Param("symbol")
	startDateStr := c.Query("startdate")
	endDateStr := c.Query("enddate")

	// Validate start date
	_, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		base.RespondError(c, "Invalid startdate", http.StatusUnprocessableEntity)
		return
	}

	// Validate end date
	_, err = time.Parse("2006-01-02", endDateStr)
	if err != nil {
		base.RespondError(c, "Invalid enddate", http.StatusUnprocessableEntity)
		return
	}

	// Get rates
	rates, err := h.RateUsecase.GetRates(symbol, startDateStr, endDateStr)
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusInternalServerError)
	}

	res := new(transformers.CollectionTransformer)
	res.RateCollection(rates)

	base.RespondJSON(c, res)
}

func (h *RateHandler) DeleteByDate(c *gin.Context) {
	date := c.Param("date")

	// Validate date
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		base.RespondError(c, "Invalid date", http.StatusBadRequest)
		return
	}

	// Delete
	err = h.RateUsecase.DeleteByDate(date)
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusInternalServerError)
		return
	}

	base.RespondDeleted(c, "Success")
}

func (h *RateHandler) Add(c *gin.Context) {
	var req models.RateRequest
	err := c.BindJSON(&req)
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusBadRequest)
		return
	}

	rate := mappingRequest(req)
	err = h.RateUsecase.Create(rate)

	res := new(transformers.Transformer)
	res.RateTransformer(rate)

	base.RespondJSON(c, res)
}

func (h *RateHandler) Update(c *gin.Context) {
	var req models.RateRequest
	err := c.BindJSON(&req)
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusBadRequest)
		return
	}

	// Get rates
	rates, err := h.RateUsecase.GetRates(req.Symbol, req.Date, req.Date)
	if err != nil {
		base.RespondError(c, err.Error(), http.StatusInternalServerError)
	}

	if len(rates) == 0 {
		base.RespondNotFound(c, "Data not found")
		return
	}

	rate := mappingRequest(req)
	err = h.RateUsecase.Update(rate)

	res := new(transformers.Transformer)
	res.RateTransformer(rate)

	base.RespondJSON(c, res)
}

func mappingRequest(req models.RateRequest) (rate models.ExchangeRate) {
	rate.Date, _ = time.Parse("2006-01-02", req.Date)
	rate.Symbol = req.Symbol
	rate.ERateBuy = req.ERate.Beli
	rate.ERateSell = req.ERate.Jual
	rate.TTBuy = req.TTCounter.Beli
	rate.TTSell = req.TTCounter.Jual
	rate.BNBuy = req.BankNotes.Beli
	rate.BNSell = req.BankNotes.Jual

	return
}