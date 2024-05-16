package main

import (
	"github.com/mdshahjahanmiah/vwap-outlier-detection/pkg/trade"
	"log/slog"
	"time"
)

func main() {
	// Initialize trade windows for BTC/USD and ETH/USD pairs
	btcWindow := trade.NewTradeWindow()
	ethWindow := trade.NewTradeWindow()

	// Simulate adding trades (this would be replaced with actual trade data in a real application)
	trades := []trade.Trade{
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "BTC/USD", Price: 50000, Volume: 1},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "BTC/USD", Price: 51000, Volume: 1.5},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "BTC/USD", Price: 52000, Volume: 2},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "BTC/USD", Price: 53000, Volume: 2.5},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "BTC/USD", Price: 54000, Volume: 3},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "ETH/USD", Price: 3000, Volume: 10},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "ETH/USD", Price: 3100, Volume: 15},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "ETH/USD", Price: 3200, Volume: 20},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "ETH/USD", Price: 3300, Volume: 25},
		{Timestamp: time.Now().Add(-1 * time.Minute), Pair: "ETH/USD", Price: 3400, Volume: 30},
	}

	// Separate trades by pairs
	btcTrades, ethTrades := trade.ProcessTradesByPair(trades)

	// Channels to receive the Volume Weighted Average Price results
	btcWeightedAveragePriceChan := make(chan float64)
	ethWeightedAveragePriceChan := make(chan float64)

	// Go routine to process BTC/USD trades
	go func() {
		btcWindow.AddTrades(btcTrades)
		btcValidTrades := btcWindow.GetValidTrades()
		btcWeightedAveragePrice := trade.CalculateVolumeWeightedAveragePrice(btcValidTrades)
		btcWeightedAveragePriceChan <- btcWeightedAveragePrice
	}()

	// Go routine to process ETH/USD trades
	go func() {
		ethWindow.AddTrades(ethTrades)
		ethValidTrades := ethWindow.GetValidTrades()
		ethWeightedAveragePrice := trade.CalculateVolumeWeightedAveragePrice(ethValidTrades)
		ethWeightedAveragePriceChan <- ethWeightedAveragePrice
	}()

	// Receive and print the Volume Weighted Average Price results
	btcWeightedAveragePrice := <-btcWeightedAveragePriceChan
	ethWeightedAveragePrice := <-ethWeightedAveragePriceChan

	slog.Info("volume weighted average price", "BTC/USD", btcWeightedAveragePrice)
	slog.Info("volume weighted average price", "ETH/USD", ethWeightedAveragePrice)
}
