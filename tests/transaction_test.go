package tests

import (
	"log"
	"testing"

	"github.com/sayyidinside/gofiber-clean-fresh/domain/entity"
	"github.com/sayyidinside/gofiber-clean-fresh/interfaces/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction_Success(t *testing.T) {
	// Create limit
	limit := entity.Limit{
		UserID:        2,
		Tenor:         6,
		CurrentLimit:  decimal.NewFromInt(10000000),
		OriginalLimit: decimal.NewFromInt(10000000),
	}
	TestDB.Create(&limit)

	token := GenerateUserTestToken()

	// Create request
	input := model.TransactionInput{
		AssetName: "Refrigerator",
		OnTheRoad: "5000000",
		AdminFee:  "50000",
		Tenor:     6,
	}

	recorder := MakeRequest(t, "POST", "/api/v1/transactions/data", input, token)
	log.Println(token)
	log.Println(recorder.Body)

	// Verify response
	assert.Equal(t, 201, recorder.Code)
	response := ParseResponse(t, recorder)

	assert.True(t, response.Success)
	assert.Equal(t, "Transaction successfully created", response.Message)

	// Verify database
	var transaction entity.Transaction
	TestDB.Where("user_id = ?", 2).First(&transaction)

	assert.Equal(t, "Refrigerator", transaction.AssetName)
	assert.True(t, decimal.NewFromInt(5000000).Equal(transaction.OnTheRoad))

	// Verify limit updated
	var updatedLimit entity.Limit
	TestDB.Where("user_id = ? AND tenor = ?", 2, 6).First(&updatedLimit)

	expectedLimit := decimal.NewFromInt(10000000).Sub(decimal.NewFromInt(5000000))
	assert.True(t, expectedLimit.Equal(updatedLimit.CurrentLimit))
}
