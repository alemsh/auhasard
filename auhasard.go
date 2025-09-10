package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// id=articleWRD
func getDefintion(word string) {
	user_agent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.wordreference.com/fren/%s", word), nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("User-Agent", user_agent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
func main() {
	getDefintion("coucou")
}
