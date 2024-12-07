package main

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type ReceiptPoints struct {
	ID     string
	Points int
}

var receipts = make(map[string]ReceiptPoints)
var logger = logrus.New()

func main() {
	router := gin.New()
	router.Use(gin.Logger())
	router.RemoveExtraSlash = true

	base := router.Group("")
	receiptsGroup := base.Group("/receipts")

	receiptsGroup.POST("/process", processReceipt)

	receiptsGroup.GET("/:id/points", retrievePoints)

	router.Run(":8080")
}

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		logger.Errorf("error binding receipt: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	retailerPoints := calculateRetailerPoints(receipt.Retailer)
	totalPoints, err := calculateTotalPoints(receipt.Total)
	if err != nil {
		logger.Errorf("error calculating total points: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid total"})
		return
	}

	itemPairPoints := calculateItemPairPoints(receipt.Items)
	itemDescriptionPoints, err := calculateItemDescriptionPoints(receipt.Items)
	if err != nil {
		logger.Errorf("error calculating item description points: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item description"})
	}

	oddDayPoints, err := calculateOddDayPoints(receipt.PurchaseDate)
	if err != nil {
		logger.Errorf("error calculating odd day points: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase date"})
	}

	timePoints, err := calculateTimePoints(receipt.PurchaseTime)
	if err != nil {
		logger.Errorf("error calculating time points: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid purchase time"})
	}

	total := retailerPoints + totalPoints + itemPairPoints + itemDescriptionPoints + oddDayPoints + timePoints

	receiptID := uuid.New().String()

	receipts[receiptID] = ReceiptPoints{
		ID:     receiptID,
		Points: total,
	}

	logger.Infof("receipt processed: %+v, Points: %d", receipt, total)
	c.JSON(http.StatusOK, gin.H{"id": receiptID})
}

func calculateTotalPoints(total string) (int, error) {
	points := 0
	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		logger.Errorf("error parsing total '%s': %v", total, err)
		return 0, err
	}

	// 50 points for round total
	if totalFloat == float64(int(totalFloat)) {
		points += 50
	}

	// 25 points for multiple of 0.25
	if int(totalFloat*100)%25 == 0 {
		points += 25
	}

	return points, nil
}

func calculateRetailerPoints(retailer string) int {
	count := 0
	for _, char := range retailer {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			count++
		}
	}
	return count
}

func calculateItemPairPoints(items []Item) int {
	points := (len(items) / 2) * 5
	return points
}

func calculateItemDescriptionPoints(items []Item) (int, error) {
	points := 0
	for _, item := range items {
		trimmedItemDescriptionLength := len(strings.TrimSpace(item.ShortDescription))
		if trimmedItemDescriptionLength%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				logger.Errorf("error parsing item price '%s': %v", item.Price, err)
				return 0, err
			}
			itemPoints := int(math.Ceil(itemPrice * 0.2))
			points += itemPoints
		}
	}
	return points, nil
}

func calculateOddDayPoints(purchaseDate string) (int, error) {
	day := strings.Split(purchaseDate, "-")[2]
	dayInt, err := strconv.Atoi(day)
	if err != nil {
		logger.Errorf("error parsing purchase day '%s': %v", day, err)
		return 0, err
	}
	if dayInt%2 == 1 {
		return 6, nil
	}
	return 0, nil
}

func calculateTimePoints(purchaseTime string) (int, error) {
	time := strings.Split(purchaseTime, ":")
	hour, err := strconv.Atoi(time[0])
	if err != nil {
		logger.Errorf("error parsing purchase time '%s': %v", purchaseTime, err)
		return 0, err
	}
	mins, err := strconv.Atoi(time[1])
	if err != nil {
		logger.Errorf("error parsing purchase time '%s': %v", purchaseTime, err)
		return 0, err
	}

	if (hour == 14 && mins != 0) || hour == 15 {
		return 10, nil
	}
	return 0, err
}

func retrievePoints(c *gin.Context) {
	id := c.Param("id")

	receiptPoints, exists := receipts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "receipt id not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"points": receiptPoints.Points})
}
