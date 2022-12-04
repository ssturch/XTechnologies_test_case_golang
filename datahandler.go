package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Переработка данных от CBR
func remakeXMLfromCBR(input ResultXML) (map[string]map[string]interface{}, error) {
	var tempIntTimeByXML [3]int
	cbrData := make(map[string]map[string]interface{})

	tempStrTimeByXML := strings.Split(input.Date, ".")
	for i := range tempStrTimeByXML {
		tempIntTimeByXML[i], _ = strconv.Atoi(tempStrTimeByXML[i])
	}
	tempTime := time.Date(tempIntTimeByXML[2], time.Month(tempIntTimeByXML[1]), tempIntTimeByXML[0], 24, 0, 0, 0, time.UTC).Unix()

	for _, elem := range input.Valutes {
		cbrData[elem.CharCode] = make(map[string]interface{})
		cbrData[elem.CharCode]["time"] = tempTime
		cbrData[elem.CharCode]["name"] = elem.Name
		tempStrElemValue := strings.ReplaceAll(elem.Value, ",", ".")
		tempFloatElemValue, _ := strconv.ParseFloat(tempStrElemValue, 8)
		cbrData[elem.CharCode]["value"] = tempFloatElemValue
	}
	if len(cbrData) == 0 {
		err = errors.New("ERROR: Conversion from CBR XML is failed!")
		return cbrData, err
	}
	return cbrData, err
}

// Переработка данных от KUCOIN
func remakeJSONfromKUCOIN(input map[string]interface{}) (map[string]interface{}, error) {
	kucoinData := make(map[string]interface{})
	BTCUSDTData := input["data"].(map[string]interface{})
	kucoinData["time"] = int(BTCUSDTData["time"].(float64)) / 1000
	kucoinData["value"], _ = strconv.ParseFloat(BTCUSDTData["buy"].(string), 64)
	if len(kucoinData) == 0 {
		err := errors.New("ERROR: Conversion from KUCOIN JSON is failed!")
		return kucoinData, err
	}
	return kucoinData, nil
}

// Приведение JSON к виду в ТЗ (соблюдение регистра для переменных)
func remakeJSONtoSample(input []byte) []byte {
	inputStr := string(input)
	re, _ := regexp.Compile("[A-Z][a-z]+")
	res := re.FindAllString(inputStr, -1)
	if len(res) > 0 {
		for i := range res {
			tempStr := strings.ToLower(res[i])
			inputStr = strings.Replace(inputStr, res[i], tempStr, -1)
		}
	} else {
		return nil
	}
	return []byte(inputStr)
}
