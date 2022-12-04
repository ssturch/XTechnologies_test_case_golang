package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ValuteData struct {
	//XMLName  xml.Name `xml:"Valute"`
	NumCode  string
	CharCode string
	Name     string
	Value    string
}
type ResultXML struct {
	XMLName xml.Name     `xml:"ValCurs"`
	Date    string       `xml:"Date,attr"`
	Valutes []ValuteData `xml:"Valute"`
}

var BTCData map[string]interface{}

// Получение данных от CBR
func getXMLfromCBR() (ResultXML, error) {
	var xmlUnmarshalled ResultXML
	var resp *http.Response
	var respBody []byte
	var out []byte
	resp, err = http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		return xmlUnmarshalled, err
	}
	if resp.StatusCode > 299 {
		resp.Body.Close()
		return xmlUnmarshalled, errors.New("ERROR: Response to CBR failed. Status " + strconv.Itoa(resp.StatusCode))
	}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return xmlUnmarshalled, err
	}

	//Мало того что CBR дает неправильно оформленный XML, так еще и в кодировке windows1251, сволочи))
	//Не знаю насколько костыльно я получил данные от CBR, прошу дать оценку))

	dec := charmap.Windows1251.NewDecoder()
	out, err = dec.Bytes(respBody)
	if err != nil {
		return xmlUnmarshalled, err
	}
	strUTF8 := string(out[:])
	correctXMLData := strings.ReplaceAll(strUTF8, "<?xml version=\"1.0\" encoding=\"windows-1251\"?>", "")
	err = xml.Unmarshal([]byte(correctXMLData), &xmlUnmarshalled)
	if err != nil {
		resp.Body.Close()
		return xmlUnmarshalled, err
	}
	resp.Body.Close()
	return xmlUnmarshalled, nil
}

// Получение данных от KUCOIN
func getJSONfromKUCOIN() (map[string]interface{}, error) {
	var resp *http.Response
	resp, err = http.Get("https://api.kucoin.com/api/v1/market/stats?symbol=BTC-USDT")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		resp.Body.Close()
		return nil, errors.New("ERROR: Response to KUKOIN failed. Status " + strconv.Itoa(resp.StatusCode))
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(respBody, &BTCData); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return BTCData, nil
}
