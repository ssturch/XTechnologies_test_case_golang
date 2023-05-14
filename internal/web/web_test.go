package web

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestGetXMLFromCBR(t *testing.T) {
	// Positive test
	pckpath, _ := os.Getwd()
	path := strings.Replace(pckpath, "\\internal\\web", "", 1)
	path = strings.Replace(path, "\\", "/", -1)

	content, err := os.ReadFile(path + "/test/testdata_TestRemakeXMLfromCBR_correct.txt")
	fmt.Println(path)
	if err != nil {
		fmt.Println(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))

	got, err := GetXMLfromCBR(ts.URL)
	want := ResultXML{
		XMLName: xml.Name{"", "ValCurs"},
		Date:    "29.04.2023",
		Valutes: []ValuteData{{
			NumCode:  "036",
			CharCode: "AUD",
			Name:     "Австралийский доллар",
			Value:    "53,2166",
		}, {
			NumCode:  "944",
			CharCode: "AZN",
			Name:     "Азербайджанский манат",
			Value:    "47,3584",
		}, {
			NumCode:  "826",
			CharCode: "GBP",
			Name:     "Фунт стерлингов Соединенного королевства",
			Value:    "100,5883",
		},
		},
	}

	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) || err != nil {
		fmt.Println("Positive test (get correct XML) - FAIL")
		t.Errorf("got %v, want %v", got, want)
	} else {
		fmt.Println("Positive test (get correct XML) - OK")
	}

	ts.Close()

	//Negative test (bad XML)

	content, err = os.ReadFile(path + "/test/testdata_TestRemakeXMLfromCBR_uncorrect_XML.txt")
	if err != nil {
		fmt.Println(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))

	_, err = GetXMLfromCBR(ts.URL)
	wantErr := fmt.Sprintf("Unmarshalling XML from body by %v caused \"XML syntax error on line 12: invalid characters between </CharCode and >\"", ts.URL)

	if err != nil && err.Error() == wantErr {
		fmt.Println("Negative test (get bad XML) - OK")
	} else if err != nil {
		fmt.Println("Negative test (get bad XML) - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("Negative test (get bad XML) - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (bad data)

	content, err = os.ReadFile(path + "/test/testdata_TestRemakeXMLfromCBR_uncorrect_data.txt")
	if err != nil {
		fmt.Println(err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))
	_, err = GetXMLfromCBR(ts.URL)

	wantErr = fmt.Sprintf("Unmarshalling XML from body by %v caused \"XML сontains incorrect data\"", ts.URL)

	if err != nil && err.Error() == wantErr {
		fmt.Println("Negative test (get bad data from XML) - OK")
	} else if err != nil {
		fmt.Println("Negative test (get bad data from XML) - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("Negative test (get bad data from XML) - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (null data)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	}))

	_, err = GetXMLfromCBR(ts.URL)

	if err != nil {
		fmt.Println("Null error test (get null XML) - OK")
	} else {
		fmt.Println("Null error test (get null XML) - FAIL")
		t.Errorf("got nil, want error: Unmarshalling XML from body by %v caused \"EOF\"", ts.URL)
	}
	ts.Close()

	//Negative test (418 error)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}))

	_, err = GetXMLfromCBR(ts.URL)
	cutUrl, _ := strings.CutPrefix(ts.URL, "http://")
	wantErr = fmt.Sprintf("Host %v return status %v with text \"%v\"", cutUrl, "418", "I'm a teapot")
	if err != nil && err.Error() == wantErr {
		fmt.Println("\"I'm a teapot\" (418) error test - OK")
	} else if err != nil {
		fmt.Println("\"I'm a teapot\" (418) error test - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("\"I'm a teapot\" (418) error test - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (503 error)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))

	_, err = GetXMLfromCBR(ts.URL)
	cutUrl, _ = strings.CutPrefix(ts.URL, "http://")
	wantErr = fmt.Sprintf("Host %v return status %v with text \"%v\"", cutUrl, "503", "Service Unavailable")
	if err != nil && err.Error() == wantErr {
		fmt.Println("\"Service Unavailable\" (503) error test - OK")
	} else if err != nil {
		fmt.Println("\"Service Unavailable\" (503) error test - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("\"Service Unavailable\" (503) error test - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

}

func TestGetJSONfromKUCOIN(t *testing.T) {
	// Positive test
	pckpath, _ := os.Getwd()
	path := strings.Replace(pckpath, "\\internal\\web", "", 1)
	path = strings.Replace(path, "\\", "/", -1)

	content, err := os.ReadFile(path + "/test/testdata_TestGetJSONfromKUCOIN_correct.txt")
	if err != nil {
		fmt.Println(err)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))

	got, err := GetJSONfromKUCOIN(ts.URL)

	want := map[string]interface{}{
		"code": "200000",
		"data": map[string]interface{}{
			"time":   "12345",
			"symbol": "BTC-USDT",
			"buy":    "12345.6",
			"sell":   "30000",
		},
	}

	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) || err != nil {
		fmt.Println("Positive test (get correct JSON) - FAIL")
		t.Errorf("got %v, want %v", got, want)
	} else {
		fmt.Println("Positive test (get correct JSON) - OK")
	}
	ts.Close()

	//Negative test (bad JSON)

	content, err = os.ReadFile(path + "/test/testdata_TestGetJSONfromKUCOIN_uncorrect_JSON.txt")
	if err != nil {
		fmt.Println(err)
	}
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))

	_, err = GetJSONfromKUCOIN(ts.URL)
	wantErr := fmt.Sprintf("Unmarshalling JSON from body by %v caused \"invalid character '\"' after object key:value pair\"", ts.URL)

	if err != nil && err.Error() == wantErr {
		fmt.Println("Negative test (get bad JSON) - OK")
	} else if err != nil {
		fmt.Println("Negative test (get bad JSON) - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("Negative test (get bad JSON) - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (bad data)

	content, err = os.ReadFile(path + "/test/testdata_TestGetJSONfromKUCOIN_uncorrect_data.txt")
	if err != nil {
		fmt.Println(err)
	}

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, string(content))
	}))
	_, err = GetJSONfromKUCOIN(ts.URL)
	wantErr = fmt.Sprintf("Parsing from JSON by %v to map[string]inteface{} caused \"JSON data is incorrect! Dont find required data\"", ts.URL)

	if err != nil && err.Error() == wantErr {
		fmt.Println("Negative test (get bad data from JSON) - OK")
	} else if err != nil {
		fmt.Println("Negative test (get bad data from JSON) - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("Negative test (get bad data from JSON) - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (null data)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "")
	}))
	_, err = GetJSONfromKUCOIN(ts.URL)
	if err != nil {
		fmt.Println("Null error test (get null JSON) - OK")
	} else {
		fmt.Println("Null error test (get null JSON) - FAIL")
		t.Errorf("got nil, want error: Unmarshalling JSON from body by %v caused \"unexpected end of JSON input\"", ts.URL)
	}
	ts.Close()

	//Negative test (418 error)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(418)
	}))

	_, err = GetJSONfromKUCOIN(ts.URL)
	cutUrl, _ := strings.CutPrefix(ts.URL, "http://")
	wantErr = fmt.Sprintf("Host %v return status %v with text \"%v\"", cutUrl, "418", "I'm a teapot")
	if err != nil && err.Error() == wantErr {
		fmt.Println("\"I'm a teapot\" (418) error test - OK")
	} else if err != nil {
		fmt.Println("\"I'm a teapot\" (418) error test - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("\"I'm a teapot\" (418) error test - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()

	//Negative test (503 error)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))

	_, err = GetJSONfromKUCOIN(ts.URL)
	cutUrl, _ = strings.CutPrefix(ts.URL, "http://")
	wantErr = fmt.Sprintf("Host %v return status %v with text \"%v\"", cutUrl, "503", "Service Unavailable")
	if err != nil && err.Error() == wantErr {
		fmt.Println("\"Service Unavailable\" (503) error test - OK")
	} else if err != nil {
		fmt.Println("\"Service Unavailable\" (503) error test - FAIL")
		t.Errorf("got %v, want error: %v", err, wantErr)
	} else {
		fmt.Println("\"Service Unavailable\" (503) error test - FAIL")
		t.Errorf("got nil err, want error: %v", wantErr)
	}
	ts.Close()
}

func TestRemakeXMLfromCBR(t *testing.T) {

	testdata := ResultXML{
		XMLName: xml.Name{},
		Date:    "10.01.2001",
		Valutes: []ValuteData{
			{
				NumCode:  "123",
				CharCode: "ABC",
				Name:     "Test_name",
				Value:    "120,4",
			},
			{
				NumCode:  "456",
				CharCode: "DEF",
				Name:     "Test_name_2",
				Value:    "400",
			},
		},
	}

	//Positive test

	want := make(map[string]map[string]interface{})
	want["ABC"] = make(map[string]interface{})
	want["ABC"]["time"] = any(979084800)
	want["ABC"]["name"] = any("Test_name")
	want["ABC"]["value"] = any(120.4)
	want["DEF"] = make(map[string]interface{})
	want["DEF"]["time"] = any(979084800)
	want["DEF"]["name"] = any("Test_name_2")
	want["DEF"]["value"] = any(400)

	got, err := RemakeXMLfromCBR(testdata)

	if err == nil && fmt.Sprint(want) == fmt.Sprint(got) {
		fmt.Println("Positive test (get map from CBR xml) - OK")
	} else {
		fmt.Println("Positive test (get map from CBR xml) - FAIL")
		t.Errorf("got %v, want %v", got, want)
	}

	//Negative test

	testdata = ResultXML{}
	wantErr := fmt.Sprintf("Parsing from %v to %v caused \"%v\"", "CBR XML", "map[string]map[string]interface{}", "Null result")
	got, err = RemakeXMLfromCBR(testdata)

	if err != nil {
		if err.Error() == wantErr {
			fmt.Println("Negative test (get error from nil CBR xml) - OK")
		} else {
			fmt.Println("Negative test (get error from nil CBR xml) - FAIL")
			t.Errorf("got %v, want %v", err, wantErr)
		}
	} else {
		fmt.Println("Negative test (get error from nil CBR xml) - FAIL")
		t.Errorf("got %v, want %v", got, nil)
	}

}

func TestRemakeJSONfromKUCOIN(t *testing.T) {

	//Positive test

	testdata := make(map[string]interface{})
	testdata["data"] = make(map[string]interface{})
	testdata["data"].(map[string]interface{})["time"] = float64(1684053378029)
	testdata["data"].(map[string]interface{})["buy"] = "26863.7"

	want := make(map[string]interface{})
	want["time"] = 1684053378
	want["value"] = 26863.7

	got, err := RemakeJSONfromKUCOIN(testdata)

	if err == nil && fmt.Sprint(want) == fmt.Sprint(got) {
		fmt.Println("Positive test (get converted JSON from KUCOIN) - OK")
	} else {
		fmt.Println("Positive test (get converted JSON from KUCOIN) - FAIL")
		t.Errorf("got %v, want %v", got, want)
	}

	//Negative test

	testdata = make(map[string]interface{})

	wantErr := fmt.Sprintf("Parsing from %v to %v caused \"%v\"", "KUCOIN JSON", "map[string]interface{}", "Null result")

	got, err = RemakeJSONfromKUCOIN(testdata)

	if err != nil {
		if err.Error() == wantErr {
			fmt.Println("Negative test (get error from nil KUCOIN json) - OK")
		} else {
			fmt.Println("Negative test (get error from nil KUCOIN json) - FAIL")
			t.Errorf("got %v, want %v", err, wantErr)
		}
	} else {
		fmt.Println("Negative test (get error from nil KUCOIN json) - FAIL")
		t.Errorf("got %v, want %v", got, nil)
	}

}

func TestRemakeJSONtoSample(t *testing.T) {

	//Positive test

	test := "{\"Total\":128,\"History\":[{\"timestamp\":1682362118,\"value\":27306.3},{\"timestamp\":1682362128,\"value\":27312.5}]}"
	want := "{\"total\":128,\"history\":[{\"timestamp\":1682362118,\"value\":27306.3},{\"timestamp\":1682362128,\"value\":27312.5}]}"

	got := RemakeJSONtoSample([]byte(test))

	if want == string(got) {
		fmt.Println("Positive test (remake JSON to downcase) - OK")
	} else {
		fmt.Println("Positive test (remake JSON to downcase) - FAIL")
		t.Errorf("got %v, want %v", got, want)
	}

	//Negative test
	test = ""
	got = RemakeJSONtoSample([]byte(test))
	if got == nil {
		fmt.Println("Negative test (remake JSON to downcase) - OK")
	} else {
		fmt.Println("Negative test (remake JSON to downcase) - FAIL")
		t.Errorf("got %v, want %v", got, nil)
	}

}
