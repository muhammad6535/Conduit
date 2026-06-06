package cost

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type ModelPrice struct {
	InputPricePer1K  float64
	OutputPricePer1K float64
}

var pricingDB = map[string]ModelPrice{
	"gpt-4o":               {InputPricePer1K: 0.0025, OutputPricePer1K: 0.01},
	"gpt-4o-2024-08-06":    {InputPricePer1K: 0.0025, OutputPricePer1K: 0.01},
	"gpt-4o-mini":          {InputPricePer1K: 0.00015, OutputPricePer1K: 0.0006},
	"gpt-4o-mini-2024-07-18": {InputPricePer1K: 0.00015, OutputPricePer1K: 0.0006},
	"gpt-4-turbo":          {InputPricePer1K: 0.01, OutputPricePer1K: 0.03},
	"gpt-4":                {InputPricePer1K: 0.03, OutputPricePer1K: 0.06},
	"gpt-3.5-turbo":        {InputPricePer1K: 0.0005, OutputPricePer1K: 0.0015},
	"claude-sonnet-4":      {InputPricePer1K: 0.003, OutputPricePer1K: 0.015},
	"claude-sonnet-4-20250514": {InputPricePer1K: 0.003, OutputPricePer1K: 0.015},
	"claude-haiku-3":       {InputPricePer1K: 0.00025, OutputPricePer1K: 0.00125},
	"claude-haiku-3-20250313": {InputPricePer1K: 0.00025, OutputPricePer1K: 0.00125},
	"claude-opus-4":        {InputPricePer1K: 0.015, OutputPricePer1K: 0.075},
}

func GetPrice(model string) (ModelPrice, bool) {
	p, ok := pricingDB[model]
	return p, ok
}

func CalculateCost(model string, promptTokens, completionTokens int) float64 {
	price, ok := pricingDB[model]
	if !ok {
		return 0
	}
	cost := (float64(promptTokens) / 1000.0 * price.InputPricePer1K) +
		(float64(completionTokens) / 1000.0 * price.OutputPricePer1K)
	return math.Round(cost*100000) / 100000
}

type UsageRecord struct {
	Model            string
	Provider         string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
	Cost             float64
	Timestamp        time.Time
}

type Tracker struct {
	monthlyCost atomic.Value
	records     []UsageRecord
	mu          sync.RWMutex
}

func NewTracker() *Tracker {
	t := &Tracker{}
	t.monthlyCost.Store(0.0)
	return t
}

func (t *Tracker) Record(model, provider string, promptTokens, completionTokens int) float64 {
	cost := CalculateCost(model, promptTokens, completionTokens)

	record := UsageRecord{
		Model:            model,
		Provider:         provider,
		PromptTokens:     promptTokens,
		CompletionTokens: completionTokens,
		TotalTokens:      promptTokens + completionTokens,
		Cost:             cost,
		Timestamp:        time.Now(),
	}

	t.mu.Lock()
	t.records = append(t.records, record)
	t.mu.Unlock()

	for {
		current := t.monthlyCost.Load().(float64)
		newTotal := current + cost
		if t.monthlyCost.CompareAndSwap(current, newTotal) {
			break
		}
	}

	return cost
}

func (t *Tracker) MonthlyCost() float64 {
	return t.monthlyCost.Load().(float64)
}

func (t *Tracker) Records() []UsageRecord {
	t.mu.RLock()
	defer t.mu.RUnlock()
	cp := make([]UsageRecord, len(t.records))
	copy(cp, t.records)
	return cp
}

func (t *Tracker) Summary() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	byModel := make(map[string]int)
	byProvider := make(map[string]float64)
	totalTokens := 0
	totalCost := 0.0

	for _, r := range t.records {
		byModel[r.Model] += r.TotalTokens
		byProvider[r.Provider] += r.Cost
		totalTokens += r.TotalTokens
		totalCost += r.Cost
	}

	s := fmt.Sprintf("Total requests: %d\n", len(t.records))
	s += fmt.Sprintf("Total tokens: %d\n", totalTokens)
	s += fmt.Sprintf("Total cost: $%.4f\n", totalCost)
	s += "By provider:\n"
	for p, c := range byProvider {
		s += fmt.Sprintf("  %s: $%.4f\n", p, c)
	}
	s += "By model:\n"
	for m, tok := range byModel {
		s += fmt.Sprintf("  %s: %d tokens\n", m, tok)
	}

	return s
}
