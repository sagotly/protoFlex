package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	actualIP, err := http.Get("https://ifconfig.me")
	if err != nil {
		fmt.Println("FAIL: Unable to fetch IP:", err)
		os.Exit(1)
	}

	// Ваше значение IP, с которым нужно сравнить (пример)

	// Проверяем, совпадает ли IP
	if actualIP != nil && actualIP.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(actualIP.Body)
		if err != nil {
			fmt.Println("FAIL: Unable to read response body:", err)
			os.Exit(1)
		}

		actualIPString := string(bodyBytes)
		fmt.Println("Actual IP:", actualIPString)
		if actualIPString == "32" {
			fmt.Println("FAIL: IP matches host")
			os.Exit(1)
		} else {
			fmt.Println("PASS: IP is isolated")
		}
	} else {
		fmt.Println("FAIL: Unable to fetch IP")
		os.Exit(1)
	}
}
