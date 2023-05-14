package web

import (
	oth "XTapi/internal"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Переработка данных от CBR
func RemakeXMLfromCBR(input ResultXML) (map[string]map[string]interface{}, error) {
	var tempIntTimeByXML [3]int
	var err error
	cbrData := make(map[string]map[string]interface{})

	tempStrTimeByXML := strings.Split(input.Date, ".")
	for i := range tempStrTimeByXML {
		tempIntTimeByXML[i], _ = strconv.Atoi(tempStrTimeByXML[i])
	}
	tempTime := time.Date(tempIntTimeByXML[2], time.Month(tempIntTimeByXML[1]), tempIntTimeByXML[0], 0, 0, 0, 0, time.UTC).Unix()

	for _, elem := range input.Valutes {
		cbrData[elem.CharCode] = make(map[string]interface{})
		cbrData[elem.CharCode]["time"] = tempTime
		cbrData[elem.CharCode]["name"] = elem.Name
		tempStrElemValue := strings.ReplaceAll(elem.Value, ",", ".")
		tempFloatElemValue, _ := strconv.ParseFloat(tempStrElemValue, 8)
		cbrData[elem.CharCode]["value"] = tempFloatElemValue
	}

	if len(cbrData) == 0 {
		return cbrData, &oth.ParsingError{
			From: "CBR XML",
			To:   "map[string]map[string]interface{}",
			Err:  errors.New("Null result"),
		}
	}
	return cbrData, err
}

// Переработка данных от KUCOIN
func RemakeJSONfromKUCOIN(input map[string]interface{}) (map[string]interface{}, error) {
	if len(input) == 0 {
		return nil, &oth.ParsingError{
			From: "KUCOIN JSON",
			To:   "map[string]interface{}",
			Err:  errors.New("Null result"),
		}
	}

	kucoinData := make(map[string]interface{})
	BTCUSDTData := input["data"].(map[string]interface{})
	kucoinData["time"] = int(BTCUSDTData["time"].(float64)) / 1000
	kucoinData["value"], _ = strconv.ParseFloat(BTCUSDTData["buy"].(string), 64)

	return kucoinData, nil
}

// Приведение JSON к виду в ТЗ (соблюдение регистра для переменных)
func RemakeJSONtoSample(input []byte) []byte {

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
