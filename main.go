package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
)

var rxDescription = regexp.MustCompile(`^[\w\s\-]+$`)
var rxPrice = regexp.MustCompile(`^(\d+)\.(\d{2})$`)
var rxRetailer = regexp.MustCompile(`^[\w\s\-&]+$`)
var rxDate = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
var rxTime = regexp.MustCompile(`^(\d{2}):(\d{2})$`)
var twoPM, _ = time.Parse("15:04", "14:00")
var fourPM, _ = time.Parse("15:04", "16:00")
var dataStore = make(map[string]Receipt)
var mu sync.RWMutex

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

func (item *Item) Validate() error {
	var errors []string

	if strings.TrimSpace(item.ShortDescription) == "" {
		errors = append(errors, "shortDescription cannot be empty")
	} else if !rxDescription.MatchString(item.ShortDescription) {
		errors = append(errors, fmt.Sprintf("invalid format for shortDescription (%s)", item.ShortDescription))
	}

	if strings.TrimSpace(item.Price) == "" {
		errors = append(errors, "price cannot be empty")
	} else if !rxPrice.MatchString(item.Price) {
		errors = append(errors, fmt.Sprintf("invalid format for price (%s)", item.Price))
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, ", "))
	}
	return nil
}

func (item *Item) PointsForItem() (int, string) {
	points := 0
	trimmedDescription := strings.TrimSpace(item.ShortDescription)
	if len(trimmedDescription)%3 == 0 {
		itemPriceFloat, _ := strconv.ParseFloat(item.Price, 64)
		points = int(math.Ceil(itemPriceFloat * 0.2))
	}
	message := fmt.Sprintf("%d point(s) for item (%s | %s)", points, item.ShortDescription, item.Price)
	return points, message
}

func (receipt *Receipt) Validate() error {
	var errors []string

	if strings.TrimSpace(receipt.Retailer) == "" {
		errors = append(errors, "retailer cannot be empty")
	} else if !rxRetailer.MatchString(receipt.Retailer) {
		errors = append(errors, fmt.Sprintf("invalid format for retailer (%s)", receipt.Retailer))
	}

	if strings.TrimSpace(receipt.PurchaseDate) == "" {
		errors = append(errors, "purchaseDate cannot be empty")
	} else if !rxDate.MatchString(receipt.PurchaseDate) {
		errors = append(errors, fmt.Sprintf("invalid format for purchaseDate (%s)", receipt.PurchaseDate))
	} else {
		_, err := time.Parse("2006-01-02", receipt.PurchaseDate)
		if err != nil {
			errors = append(errors, fmt.Sprintf("purchaseDate cannot be parsed (%s)", receipt.PurchaseDate))
		}
	}

	if strings.TrimSpace(receipt.PurchaseTime) == "" {
		errors = append(errors, "purchaseTime cannot be empty")
	} else if !rxTime.MatchString(receipt.PurchaseTime) {
		errors = append(errors, fmt.Sprintf("invalid format for purchaseTime (%s)", receipt.PurchaseTime))
	} else {
		_, err := time.Parse("15:04", receipt.PurchaseTime)
		if err != nil {
			errors = append(errors, fmt.Sprintf("purchaseTime cannot be parsed (%s)", receipt.PurchaseTime))
		}
	}

	if strings.TrimSpace(receipt.Total) == "" {
		errors = append(errors, "total cannot be empty")
	} else if !rxPrice.MatchString(receipt.Total) {
		errors = append(errors, fmt.Sprintf("invalid format for total (%s)", receipt.Total))
	}

	calculatedTotal := 0.0
	for index, item := range receipt.Items {
		itemPriceFloat, _ := strconv.ParseFloat(item.Price, 64)
		calculatedTotal += itemPriceFloat
		err := item.Validate()
		if err != nil {
			errors = append(errors, fmt.Sprintf("item %d errors ... %s", index, err))
		}
	}

	calculatedTotalString := fmt.Sprintf("%.2f", calculatedTotal)
	if calculatedTotalString != receipt.Total {
		errors = append(errors, fmt.Sprintf("sum of item prices (%s) != given total (%s)", calculatedTotalString, receipt.Total))
	}

	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, " | "))
	}
	return nil
}

func (receipt *Receipt) PointsForRetailerName() (int, string) {
	points := 0
	for _, char := range receipt.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points += 1
		}
	}
	message := fmt.Sprintf("%d points for retailer name (%s)", points, receipt.Retailer)
	return points, message
}

func (receipt *Receipt) centsString() string {
	return receipt.Total[len(receipt.Total)-2:]
}

func (receipt *Receipt) PointsForRoundDollarAmount() (int, string) {
	points := 0
	if receipt.centsString() == "00" {
		points = 50
	}
	message := fmt.Sprintf("%d points for round dollar amount (%s)", points, receipt.Total)
	return points, message
}

func (receipt *Receipt) PointsForCentsMultiple25() (int, string) {
	points := 0
	cents, _ := strconv.Atoi(receipt.centsString())
	if cents%25 == 0 {
		points = 25
	}
	message := fmt.Sprintf("%d points for being multiple of 0.25 (%s)", points, receipt.Total)
	return points, message
}

func (receipt *Receipt) PointsForNumItems() (int, string) {
	pairsOfItems := int(math.Floor(float64(len(receipt.Items)) / 2))
	points := pairsOfItems * 5
	message := fmt.Sprintf("%d points for number of items (%d)", points, len(receipt.Items))
	return points, message
}

func (receipt *Receipt) PointsForPurchaseDate() (int, string) {
	points := 0
	match := rxDate.FindStringSubmatch(receipt.PurchaseDate)
	dayInt, _ := strconv.Atoi(match[3])
	if !(dayInt%2 == 0) {
		points = 6
	}
	message := fmt.Sprintf("%d points for purchase day being odd (%s)", points, receipt.PurchaseDate)
	return points, message
}

func (receipt *Receipt) PointsForPurchaseTime() (int, string) {
	points := 0
	timeObj, _ := time.Parse("15:04", receipt.PurchaseTime)
	if timeObj.After(twoPM) && timeObj.Before(fourPM) {
		points = 10
	}
	message := fmt.Sprintf("%d points for time of purchase between 2pm and 4pm (%s)", points, receipt.PurchaseTime)
	return points, message
}

func (receipt *Receipt) GetTotalPointsAndBreakdown() (int, []string) {
	var breakdown []string
	totalPoints := 0

	points, message := receipt.PointsForRetailerName()
	totalPoints += points
	breakdown = append(breakdown, message)

	points, message = receipt.PointsForRoundDollarAmount()
	totalPoints += points
	breakdown = append(breakdown, message)

	points, message = receipt.PointsForCentsMultiple25()
	totalPoints += points
	breakdown = append(breakdown, message)

	points, message = receipt.PointsForNumItems()
	totalPoints += points
	breakdown = append(breakdown, message)

	for _, item := range receipt.Items {
		points, message = item.PointsForItem()
		totalPoints += points
		breakdown = append(breakdown, message)
	}

	points, message = receipt.PointsForPurchaseDate()
	totalPoints += points
	breakdown = append(breakdown, message)

	points, message = receipt.PointsForPurchaseTime()
	totalPoints += points
	breakdown = append(breakdown, message)

	return totalPoints, breakdown
}

func printDelimiter() {
	fmt.Printf("\n\n" + strings.Repeat("-", 80) + "\n\n")
}

func LoadJSON(filename string, v interface{}) error {
	fp, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Error opening file: %w", err)
	}
	defer fp.Close()

	data, err2 := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err2)
	}

	err3 := json.Unmarshal(data, v)
	if err3 != nil {
		return fmt.Errorf("Error unmarshaling JSON: %w", err3)
	}

	return nil
}

func handleError(writer http.ResponseWriter, statusCode int, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	json.NewEncoder(writer).Encode(map[string]string{"error": message})
	log.Println(fmt.Sprintf("(%d) ERROR %s", statusCode, message))
}

func handleReceiptPost(writer http.ResponseWriter, request *http.Request) {
	log.Println(fmt.Sprintf("%s %s", request.Method, request.URL.Path))
	if request.Method != http.MethodPost {
		message := "Only POST is allowed"
		handleError(writer, http.StatusMethodNotAllowed, message)
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		message := "Could not read request body"
		handleError(writer, http.StatusBadRequest, message)
		return
	}
	defer request.Body.Close()

	var receipt Receipt
	err2 := json.Unmarshal(body, &receipt)
	if err2 != nil {
		message := fmt.Sprintf("Error unmarshaling JSON: %s", err2.Error())
		handleError(writer, http.StatusBadRequest, message)
		return
	}
	err3 := receipt.Validate()
	if err3 != nil {
		message := fmt.Sprintf("Validation errors: %s", err3.Error())
		handleError(writer, http.StatusBadRequest, message)
		return
	}

	id := uuid.New().String()
	mu.Lock()
	dataStore[id] = receipt
	mu.Unlock()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(map[string]string{"id": id})
	log.Println(fmt.Sprintf("(%d) OK %s", http.StatusOK, id))
}

func handleGetPoints(writer http.ResponseWriter, request *http.Request) {
	log.Println(fmt.Sprintf("%s %s", request.Method, request.URL.Path))
	if request.Method != http.MethodGet {
		message := "Only GET is allowed"
		handleError(writer, http.StatusMethodNotAllowed, message)
		return
	}
	parts := strings.Split(request.URL.Path, "/")
	id := parts[2]
	mu.RLock()
	receipt, exists := dataStore[id]
	mu.RUnlock()
	if !exists {
		message := fmt.Sprintf("receipt %s not found", id)
		handleError(writer, http.StatusNotFound, message)
		return
	}
	points, _ := receipt.GetTotalPointsAndBreakdown()
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string]int{"points": points})
	log.Println(fmt.Sprintf("(%d) OK points for %s", http.StatusOK, id))
}

func handleGetBreakdown(writer http.ResponseWriter, request *http.Request) {
	log.Println(fmt.Sprintf("%s %s", request.Method, request.URL.Path))
	if request.Method != http.MethodGet {
		message := "Only GET is allowed"
		handleError(writer, http.StatusMethodNotAllowed, message)
		return
	}
	parts := strings.Split(request.URL.Path, "/")
	id := parts[2]
	mu.RLock()
	receipt, exists := dataStore[id]
	mu.RUnlock()
	if !exists {
		message := fmt.Sprintf("receipt %s not found", id)
		handleError(writer, http.StatusNotFound, message)
		return
	}
	points, breakdown := receipt.GetTotalPointsAndBreakdown()
	breakdown = append(breakdown, fmt.Sprintf("%d points total", points))
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(map[string][]string{"breakdown": breakdown})
	log.Println(fmt.Sprintf("(%d) OK breakdown for %s", http.StatusOK, id))
}

func main() {
	fmt.Println("This is the receipt processor!")

	// var receipt Receipt
	// err := LoadJSON("example3.json", &receipt)
	// if err != nil {
	// 	fmt.Println("Error loading JSON:", err)
	// }
	// fmt.Printf("This is the receipt data that was loaded: %v\n", receipt)
	// errReceipt := receipt.Validate()
	// fmt.Printf("Validation errors for receipt: %v\n", errReceipt)
	// totalPoints, breakdown := receipt.GetTotalPointsAndBreakdown()
	// fmt.Printf("\n%d total points\n\nbreakdown:\n%s\n", totalPoints, strings.Join(breakdown, "\n"))

	http.HandleFunc("/receipts/process", handleReceiptPost)
	http.HandleFunc("/receipts/{id}/points", handleGetPoints)
	http.HandleFunc("/receipts/{id}/breakdown", handleGetBreakdown)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
