package main

import (
	"strings"
	"testing"
)

var receiptExample1, receiptExample2, receiptExample3 Receipt

var err1 = LoadJSON("example1.json", &receiptExample1)
var err2 = LoadJSON("example2.json", &receiptExample2)
var err3 = LoadJSON("example3.json", &receiptExample3)

func TestItemValidate1(t *testing.T) {
	item := Item{
		ShortDescription: "Skittles",
		Price:            "1.50",
	}
	err := item.Validate()
	if err != nil {
		t.Errorf("Validation errors for item %v", item)
	}
}

func TestItemValidateBadPrice(t *testing.T) {
	item := Item{
		ShortDescription: "Skittles",
		Price:            "1.50.2",
	}
	err := item.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for item %v", item)
	}
	if !strings.Contains(err.Error(), "invalid format for price") {
		t.Errorf("Validation error should contain 'invalid format for price' ... %s", err.Error())
	}
}

func TestItemValidateEmptyPrice(t *testing.T) {
	item := Item{
		ShortDescription: "Skittles",
		Price:            "",
	}
	err := item.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for item %v", item)
	}
	if !strings.Contains(err.Error(), "price cannot be empty") {
		t.Errorf("Validation error should contain 'price cannot be empty' ... %s", err.Error())
	}
}

func TestItemValidateBadDescription(t *testing.T) {
	item := Item{
		ShortDescription: "Skittles #5",
		Price:            "1.50",
	}
	err := item.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for item %v", item)
	}
	if !strings.Contains(err.Error(), "invalid format for shortDescription") {
		t.Errorf("Validation error should contain 'invalid format for shortDescription' ... %s", err.Error())
	}
}

func TestItemValidateEmptyDescription(t *testing.T) {
	item := Item{
		ShortDescription: "",
		Price:            "1.50",
	}
	err := item.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for item %v", item)
	}
	if !strings.Contains(err.Error(), "shortDescription cannot be empty") {
		t.Errorf("Validation error should contain 'shortDescription cannot be empty' ... %s", err.Error())
	}
}

func TestItemPoints1(t *testing.T) {
	item := Item{
		ShortDescription: "Skittles",
		Price:            "1.50",
	}
	points, message := item.PointsForItem()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestItemPoints2(t *testing.T) {
	item := Item{
		ShortDescription: "Dasani",
		Price:            "1.40",
	}
	points, message := item.PointsForItem()
	if points != 1 {
		t.Errorf("Should have 1 point ... %s", message)
	}
}

func TestItemPoints3(t *testing.T) {
	item := Item{
		ShortDescription: "Emils Cheese Pizza",
		Price:            "12.25",
	}
	points, message := item.PointsForItem()
	if points != 3 {
		t.Errorf("Should have 3 points ... %s", message)
	}
}

func TestReceiptValidate1(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err != nil {
		t.Errorf("Validation errors for item %v", receipt)
	}
}

func TestReceiptValidateBadRetailer(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store!!",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "invalid format for retailer") {
		t.Errorf("Validation error should contain 'invalid format for retailer' ... %s", err.Error())
	}
}

func TestReceiptValidateEmptyRetailer(t *testing.T) {
	receipt := Receipt{
		Retailer:     "",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "retailer cannot be empty") {
		t.Errorf("Validation error should contain 'retailer cannot be empty' ... %s", err.Error())
	}
}

func TestReceiptRetailerNamePoints1(t *testing.T) {
	receipt := Receipt{
		Retailer: "Corner Store",
	}
	points, message := receipt.PointsForRetailerName()
	if points != 11 {
		t.Errorf("Should have 11 points ... %s", message)
	}
}

func TestReceiptRetailerNamePoints2(t *testing.T) {
	receipt := Receipt{
		Retailer: "Walgreens",
	}
	points, message := receipt.PointsForRetailerName()
	if points != 9 {
		t.Errorf("Should have 9 points ... %s", message)
	}
}

func TestReceiptValidateBadPurchaseDate1(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-110",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "invalid format for purchaseDate") {
		t.Errorf("Validation error should contain 'invalid format for purchaseDate' ... %s", err.Error())
	}
}

func TestReceiptValidateBadPurchaseDate2(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-22-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "purchaseDate cannot be parsed") {
		t.Errorf("Validation error should contain 'purchaseDate cannot be parsed' ... %s", err.Error())
	}
}

func TestReceiptValidateEmptyPurchaseDate(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "purchaseDate cannot be empty") {
		t.Errorf("Validation error should contain 'purchaseDate cannot be empty' ... %s", err.Error())
	}
}

func TestReceiptPurchaseDatePoints1(t *testing.T) {
	receipt := Receipt{
		PurchaseDate: "2024-12-11",
	}
	points, message := receipt.PointsForPurchaseDate()
	if points != 6 {
		t.Errorf("Should have 6 points ... %s", message)
	}
}

func TestReceiptPurchaseDatePoints2(t *testing.T) {
	receipt := Receipt{
		PurchaseDate: "2024-12-12",
	}
	points, message := receipt.PointsForPurchaseDate()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptValidateBadPurchaseTime1(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:051",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "invalid format for purchaseTime") {
		t.Errorf("Validation error should contain 'invalid format for purchaseTime' ... %s", err.Error())
	}
}

func TestReceiptValidateBadPurchaseTime2(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "25:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "purchaseTime cannot be parsed") {
		t.Errorf("Validation error should contain 'purchaseTime cannot be parsed' ... %s", err.Error())
	}
}

func TestReceiptValidateEmptyPurchaseTime(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "purchaseTime cannot be empty") {
		t.Errorf("Validation error should contain 'purchaseTime cannot be empty' ... %s", err.Error())
	}
}

func TestReceiptPurchaseTimePoints1(t *testing.T) {
	receipt := Receipt{
		PurchaseTime: "15:05",
	}
	points, message := receipt.PointsForPurchaseTime()
	if points != 10 {
		t.Errorf("Should have 10 points ... %s", message)
	}
}

func TestReceiptPurchaseTimePoints2(t *testing.T) {
	receipt := Receipt{
		PurchaseTime: "11:05",
	}
	points, message := receipt.PointsForPurchaseTime()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptPurchaseTimePoints3(t *testing.T) {
	receipt := Receipt{
		PurchaseTime: "17:05",
	}
	points, message := receipt.PointsForPurchaseTime()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptValidateBadTotal1(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.00a",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "invalid format for total") {
		t.Errorf("Validation error should contain 'invalid format for total' ... %s", err.Error())
	}
}

func TestReceiptValidateBadTotal2(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "3.05",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "sum of item prices") {
		t.Errorf("Validation error should contain 'sum of item prices' ... %s", err.Error())
	}
}

func TestReceiptValidateEmptyTotal(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Corner Store",
		PurchaseDate: "2024-12-11",
		PurchaseTime: "15:05",
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
		Total: "",
	}
	err := receipt.Validate()
	if err == nil {
		t.Errorf("Should have a validation error for receipt %v", receipt)
	}
	if !strings.Contains(err.Error(), "total cannot be empty") {
		t.Errorf("Validation error should contain 'total cannot be empty' ... %s", err.Error())
	}
}

func TestReceiptRoundDollarPoints1(t *testing.T) {
	receipt := Receipt{
		Total: "3.00",
	}
	points, message := receipt.PointsForRoundDollarAmount()
	if points != 50 {
		t.Errorf("Should have 50 points ... %s", message)
	}
}

func TestReceiptRoundDollarPoints2(t *testing.T) {
	receipt := Receipt{
		Total: "3.50",
	}
	points, message := receipt.PointsForRoundDollarAmount()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptCentsMultiple25Points1(t *testing.T) {
	receipt := Receipt{
		Total: "3.00",
	}
	points, message := receipt.PointsForCentsMultiple25()
	if points != 25 {
		t.Errorf("Should have 25 points ... %s", message)
	}
}

func TestReceiptCentsMultiple25Points2(t *testing.T) {
	receipt := Receipt{
		Total: "3.25",
	}
	points, message := receipt.PointsForCentsMultiple25()
	if points != 25 {
		t.Errorf("Should have 25 points ... %s", message)
	}
}

func TestReceiptCentsMultiple25Points3(t *testing.T) {
	receipt := Receipt{
		Total: "3.50",
	}
	points, message := receipt.PointsForCentsMultiple25()
	if points != 25 {
		t.Errorf("Should have 25 points ... %s", message)
	}
}

func TestReceiptCentsMultiple25Points4(t *testing.T) {
	receipt := Receipt{
		Total: "3.75",
	}
	points, message := receipt.PointsForCentsMultiple25()
	if points != 25 {
		t.Errorf("Should have 25 points ... %s", message)
	}
}

func TestReceiptCentsMultiple25Points5(t *testing.T) {
	receipt := Receipt{
		Total: "3.85",
	}
	points, message := receipt.PointsForCentsMultiple25()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptNumItemsPoints1(t *testing.T) {
	receipt := Receipt{
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
	}
	points, message := receipt.PointsForNumItems()
	if points != 0 {
		t.Errorf("Should have 0 points ... %s", message)
	}
}

func TestReceiptNumItemsPoints2(t *testing.T) {
	receipt := Receipt{
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
	}
	points, message := receipt.PointsForNumItems()
	if points != 5 {
		t.Errorf("Should have 5 points ... %s", message)
	}
}

func TestReceiptNumItemsPoints3(t *testing.T) {
	receipt := Receipt{
		Items: []Item{
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
			{
				ShortDescription: "Skittles",
				Price:            "1.50",
			},
		},
	}
	points, message := receipt.PointsForNumItems()
	if points != 10 {
		t.Errorf("Should have 10 points ... %s", message)
	}
}

func TestJSONExample1(t *testing.T) {
	err := receiptExample1.Validate()
	if err != nil {
		t.Errorf("Validation errors for receipt %v", receiptExample1)
	}
	totalPoints, breakdown := receiptExample1.GetTotalPointsAndBreakdown()
	if totalPoints != 15 {
		t.Errorf("Should have 15 points ... %v", breakdown)
	}
	lenBreakdown := len(breakdown)
	if lenBreakdown != 8 {
		t.Errorf("Should have 8 items in the breakdown not %d ... %v", lenBreakdown, breakdown)
	}
}

func TestJSONExample2(t *testing.T) {
	err := receiptExample2.Validate()
	if err != nil {
		t.Errorf("Validation errors for receipt %v", receiptExample2)
	}
	totalPoints, breakdown := receiptExample2.GetTotalPointsAndBreakdown()
	if totalPoints != 28 {
		t.Errorf("Should have 28 points ... %v", breakdown)
	}
	lenBreakdown := len(breakdown)
	if lenBreakdown != 11 {
		t.Errorf("Should have 11 items in the breakdown not %d ... %v", lenBreakdown, breakdown)
	}
}

func TestJSONExample3(t *testing.T) {
	err := receiptExample3.Validate()
	if err != nil {
		t.Errorf("Validation errors for receipt %v", receiptExample3)
	}
	totalPoints, breakdown := receiptExample3.GetTotalPointsAndBreakdown()
	if totalPoints != 109 {
		t.Errorf("Should have 109 points ... %v", breakdown)
	}
	lenBreakdown := len(breakdown)
	if lenBreakdown != 10 {
		t.Errorf("Should have 10 items in the breakdown not %d ... %v", lenBreakdown, breakdown)
	}
}
