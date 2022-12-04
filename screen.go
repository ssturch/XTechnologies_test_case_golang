package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"os"
//	"os/exec"
//	"strconv"
//)
//
//func clearView() {
//	cmd := exec.Command("clear")
//	cmd.Stdout = os.Stdout
//	cmd.Run()
//}
//func btcusdtView(method string) {
//	if method == "GET" {
//		clearView()
//		fmt.Print("BTC-USDT: " + strconv.FormatFloat(APIResult["LastValue_btcusd"].(float64), 'g', -1, 64))
//	}
//	if method == "POST" {
//		clearView()
//		var unmarshalRes map[string]interface{}
//		_ = json.Unmarshal(APIResult["History_btcusd"].([]byte), &unmarshalRes)
//		fmt.Println("BTC-USDT:HISTORY")
//		for _, value := range unmarshalRes["History"].([]interface{}) {
//			tempmap := value.(map[string]interface{})
//			fmt.Println("TIME: " + strconv.FormatFloat(tempmap["timestamp"].(float64), 'g', -1, 64) + " VALUE: " + strconv.FormatFloat(tempmap["value"].(float64), 'g', -1, 64))
//		}
//	}
//}
