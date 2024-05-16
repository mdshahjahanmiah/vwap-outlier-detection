package trade

import (
	"testing"
	"time"
)

// Test_NewTradeWindow checks if a new TradeWindow is initialized correctly
func Test_NewTradeWindow(t *testing.T) {
	tw := NewTradeWindow()
	if len(tw.Trades) != 0 {
		t.Errorf("Expected TradeWindow to be initialized with 0 trades, got %d", len(tw.Trades))
	}
}

// Test_AddTrades ensures that trades are added correctly and old trades are removed
func Test_AddTrades(t *testing.T) {
	tw := NewTradeWindow()

	trade1 := createTrade(1, "BTC/USD", 50000, 1)
	trade2 := createTrade(3, "BTC/USD", 51000, 1.5) // Older than 2 minutes

	tw.AddTrades([]Trade{trade1, trade2})

	if len(tw.Trades) != 1 {
		t.Errorf("Expected 1 trade in TradeWindow, got %d", len(tw.Trades))
	}

	if tw.Trades[0] != trade1 {
		t.Errorf("Expected trade %+v, got %+v", trade1, tw.Trades[0])
	}
}

// Test_GetValidTrades checks if outliers are filtered out correctly using the IQR method
func Test_GetValidTrades(t *testing.T) {
	tw := NewTradeWindow()

	trades := []Trade{
		createTrade(1, "BTC/USD", 50000, 1),
		createTrade(1, "BTC/USD", 50010, 1.5),
		createTrade(1, "BTC/USD", 50020, 2),
		createTrade(1, "BTC/USD", 50030, 2.5),
		createTrade(1, "BTC/USD", 100000, 3), // Outlier
	}

	tw.AddTrades(trades)

	validTrades := tw.GetValidTrades()

	if len(validTrades) != 4 {
		t.Errorf("Expected 4 valid trades, got %d", len(validTrades))
	}

	for _, trade := range validTrades {
		if trade.Price == 100000 {
			t.Error("Outlier trade was not filtered out")
		}
	}
}

// Test_CalculateVolumeWeightedAveragePrice checks if the VWAP calculation is correct
func Test_CalculateVolumeWeightedAveragePrice(t *testing.T) {
	trades := []Trade{
		createTrade(1, "BTC/USD", 50000, 1),
		createTrade(1, "BTC/USD", 51000, 1.5),
		createTrade(1, "BTC/USD", 52000, 2),
	}

	expectedVWAP := (50000*1 + 51000*1.5 + 52000*2) / (1 + 1.5 + 2)
	vwap := CalculateVolumeWeightedAveragePrice(trades)

	if vwap != expectedVWAP {
		t.Errorf("Expected VWAP %.2f, got %.2f", expectedVWAP, vwap)
	}
}

// Test_CalculateVWAPWithNoTrades checks if VWAP is zero when there are no trades
func Test_CalculateVWAPWithNoTrades(t *testing.T) {
	trades := []Trade{}

	vwap := CalculateVolumeWeightedAveragePrice(trades)

	if vwap != 0 {
		t.Errorf("Expected VWAP 0, got %.2f", vwap)
	}
}

// Helper function to create a Trade
func createTrade(minutesAgo int, pair string, price, volume float64) Trade {
	return Trade{
		Timestamp: time.Now().Add(time.Duration(-minutesAgo) * time.Minute),
		Pair:      pair,
		Price:     price,
		Volume:    volume,
	}
}
