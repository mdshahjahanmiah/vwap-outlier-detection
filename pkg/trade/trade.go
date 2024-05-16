package trade

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Trade represents a trade data structure containing timestamp, trading pair, price, and volume.
type Trade struct {
	Timestamp time.Time
	Pair      string
	Price     float64
	Volume    float64
}

// TradeWindow manages a collection of trades within a 2-minute window.
type TradeWindow struct {
	Trades []Trade
	mu     sync.Mutex
}

// NewTradeWindow initializes and returns a new TradeWindow.
func NewTradeWindow() *TradeWindow {
	return &TradeWindow{
		Trades: make([]Trade, 0),
	}
}

// AddTrades adds a batch of new trades to the window without removing old trades.
func (tw *TradeWindow) AddTrades(newTrades []Trade) {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	cutoff := time.Now().Add(-2 * time.Minute)
	for _, trade := range newTrades {
		if trade.Timestamp.After(cutoff) {
			tw.Trades = append(tw.Trades, trade)
		}
	}

	// Remove old trades
	filteredTrades := make([]Trade, 0)
	for _, t := range tw.Trades {
		if t.Timestamp.After(cutoff) {
			filteredTrades = append(filteredTrades, t)
		}
	}
	tw.Trades = filteredTrades

	fmt.Printf("Trades after adding: %+v\n", tw.Trades)
}

// GetValidTrades filters out trades that are considered outliers using the IQR method.
func (tw *TradeWindow) GetValidTrades() []Trade {
	tw.mu.Lock()
	defer tw.mu.Unlock()

	if len(tw.Trades) == 0 {
		return nil
	}

	// Extract trade prices
	prices := make([]float64, len(tw.Trades))
	for i, t := range tw.Trades {
		prices[i] = t.Price
	}

	// Sort prices
	sort.Float64s(prices)

	// Calculate Q1, Q3 and IQR
	q1 := quantile(prices, 0.25)
	q3 := quantile(prices, 0.75)
	iqr := q3 - q1

	fmt.Printf("Q1: %f, Q3: %f, IQR: %f\n", q1, q3, iqr)

	// Filter out outliers
	validTrades := make([]Trade, 0)
	for _, t := range tw.Trades {
		if t.Price >= (q1-1.5*iqr) && t.Price <= (q3+1.5*iqr) {
			validTrades = append(validTrades, t)
		}
	}

	fmt.Printf("Valid trades: %+v\n", validTrades)
	return validTrades
}

// CalculateVolumeWeightedAveragePrice computes the Volume Weighted Average Price (VWAP) excluding outliers.
func CalculateVolumeWeightedAveragePrice(trades []Trade) float64 {
	totalVolume := 0.0
	weightedPriceSum := 0.0
	for _, t := range trades {
		totalVolume += t.Volume
		weightedPriceSum += t.Price * t.Volume
	}
	if totalVolume == 0 {
		return 0
	}
	return weightedPriceSum / totalVolume
}

// Helper function to calculate quantiles
func quantile(data []float64, percentile float64) float64 {
	index := (percentile) * float64(len(data)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))
	if lower == upper {
		return data[lower]
	}
	lowerValue := data[lower]
	upperValue := data[upper]
	return lowerValue + (upperValue-lowerValue)*(index-float64(lower))
}

// ProcessTradesByPair separates trades into two slices based on the trading pair
func ProcessTradesByPair(trades []Trade) ([]Trade, []Trade) {
	var btcTrades []Trade
	var ethTrades []Trade
	for _, t := range trades {
		if t.Pair == "BTC/USD" {
			btcTrades = append(btcTrades, t)
		} else if t.Pair == "ETH/USD" {
			ethTrades = append(ethTrades, t)
		}
	}
	return btcTrades, ethTrades
}
