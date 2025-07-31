package main

import (
	"fmt"
	"net/http"
)

func getDefintion(word string) {}
	user_agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
	req, err := http.NewRequest("GET", fmt.Sprintf("https://wordreference/fren/%s", word), nil)
	req.Header.Add("User-Agent", user_agent)

	client := &http.Client{}
	resp, err := client.Do(req)
}
func main() {
	fmt.Println("coucou")
}
