package priceoracle_test

import (
	"context"
	hopabi "github.com/redesblock/hop/contracts/abi"
	"io"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/redesblock/hop/core/settle/swap/priceoracle"
	"github.com/redesblock/hop/core/transaction"
	transactionmock "github.com/redesblock/hop/core/transaction/mock"
	"github.com/redesblock/hop/core/util/logging"
)

var (
	priceOracleABI = transaction.ParseABIUnchecked(hopabi.PriceOracleABI)
)

func TestExchangeGetPrice(t *testing.T) {
	priceOracleAddress := common.HexToAddress("0xabcd")

	expectedPrice := big.NewInt(100)
	expectedDeduce := big.NewInt(200)

	result := make([]byte, 64)
	expectedPrice.FillBytes(result[0:32])
	expectedDeduce.FillBytes(result[32:64])

	ex := priceoracle.New(
		logging.New(io.Discard, 0),
		priceOracleAddress,
		transactionmock.New(
			transactionmock.WithABICall(
				&priceOracleABI,
				priceOracleAddress,
				result,
				"getPrice",
			),
		),
		1,
	)

	price, deduce, err := ex.GetPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if expectedPrice.Cmp(price) != 0 {
		t.Fatalf("got wrong price. wanted %d, got %d", expectedPrice, price)
	}

	if expectedDeduce.Cmp(deduce) != 0 {
		t.Fatalf("got wrong deduce. wanted %d, got %d", expectedDeduce, deduce)
	}
}
