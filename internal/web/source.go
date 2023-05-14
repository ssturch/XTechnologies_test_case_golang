package web

import (
	oth "XTapi/internal"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"strings"
)

type ValuteData struct {
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

// Получение данных от CBR
func GetXMLfromCBR(url string) (ResultXML, error) {
	var xmlUnmarshalled ResultXML
	var resp *http.Response
	var respBody []byte
	var out []byte
	var err error

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MaxVersion: tls.VersionTLS13,
		},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)

	//Добавлен хэдер так как без него возвращает 403 ошибку
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)
	resp, err = client.Do(req)

	if err != nil {
		return xmlUnmarshalled, err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		resp.Body.Close()
		return xmlUnmarshalled, &oth.RequestError{Host: resp.Request.Host, Status: resp.StatusCode, Text: http.StatusText(resp.StatusCode), Err: err}
	}
	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return xmlUnmarshalled, &oth.ReaderError{ReadFrom: "Response body by " + url, Err: err}
	}

	//Мало того что CBR дает неправильно оформленный XML, так еще и в кодировке windows1251, сволочи))
	//Не знаю насколько костыльно я получил данные от CBR, прошу дать оценку))

	dec := charmap.Windows1251.NewDecoder()
	out, err = dec.Bytes(respBody)
	if err != nil {
		return xmlUnmarshalled, &oth.DecodeUnmarshallError{Do: "Decoding", Cause: "XML from body by " + url, Err: err}
	}
	strUTF8 := string(out[:])
	correctXMLData := strings.ReplaceAll(strUTF8, "<?xml version=\"1.0\" encoding=\"windows-1251\"?>", "")
	err = xml.Unmarshal([]byte(correctXMLData), &xmlUnmarshalled)

	if err != nil {
		resp.Body.Close()
		return xmlUnmarshalled, &oth.DecodeUnmarshallError{Do: "Unmarshalling", Cause: "XML from body by " + url, Err: err}
	}
	v := xmlUnmarshalled.Valutes
	for i := 0; i < len(v); i++ {
		if v[i].Name == "" || v[i].NumCode == "" || v[i].CharCode == "" || v[i].Value == "" {
			return xmlUnmarshalled, &oth.DecodeUnmarshallError{Do: "Unmarshalling", Cause: "XML from body by " + url, Err: errors.New("XML сontains incorrect data")}
		}
	}

	resp.Body.Close()
	return xmlUnmarshalled, nil
}

// Получение данных от KUCOIN
func GetJSONfromKUCOIN(url string) (map[string]interface{}, error) {
	var BTCData map[string]interface{}
	var resp *http.Response
	var err error
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		resp.Body.Close()
		return nil, &oth.RequestError{Host: resp.Request.Host, Status: resp.StatusCode, Text: http.StatusText(resp.StatusCode), Err: err}
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(respBody, &BTCData); err != nil {
		resp.Body.Close()
		return nil, &oth.DecodeUnmarshallError{Do: "Unmarshalling", Cause: "JSON from body by " + url, Err: err}
	}
	resp.Body.Close()

	chkJSONdata := BTCData["data"].(map[string]interface{})
	chkJSONtime := chkJSONdata["time"]
	chkJSONvalue := chkJSONdata["buy"]

	if chkJSONdata == nil || chkJSONtime == nil || chkJSONvalue == nil {
		return nil, &oth.ParsingError{
			From: "JSON by " + url,
			To:   "map[string]inteface{}",
			Err:  errors.New("JSON data is incorrect! Dont find required data"),
		}
	}

	return BTCData, nil
}
