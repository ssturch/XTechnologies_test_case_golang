package calc

import (
	"math"
)

// Пересчет валюты (рубль)
func ClcRUBBTC(rubusdValue float64, btcusdValue float64) float64 {
	res := round(rubusdValue*btcusdValue, 2)
	return res
}

// Пересчет валюты (все валюты)
func ClcAllValutesBTC(valutes map[string]map[string]interface{}, btcusdValue float64) map[string]float64 {
	res := make(map[string]float64)
	for key, _ := range valutes {
		res[key] = round(ClcRUBBTC(valutes["USD"]["value"].(float64), btcusdValue)/valutes[key]["value"].(float64), 2)
	}
	res["RUB"] = ClcRUBBTC(valutes["USD"]["value"].(float64), btcusdValue)
	return res
}

// Округление веществ. числа с заданной точностью
func round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}
	return rounder / pow
}
