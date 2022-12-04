package main

import "fmt"

// Будущий логгер для ошибок
func errorLogger(err error) {
	fmt.Println(err)
}
