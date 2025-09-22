package main

import (
	"fmt"
	//"io"
	"errors"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type translation struct {
	kind map[string]interpretation
}

type interpretation struct {
	from    string
	to      string
	context []string
}

func parseRoot(root *goquery.Selection) *translation {
	t := new(translation)

	root.Each(func(i int, s *goquery.Selection) {
		translationType, err := selectTranslationType(s, t)
		if err != nil {
			log.Fatal(err)
		}

		
	})

	return t
}

func selectTranslationType(s *goquery.Selection, t *translation) (string, error) {
	tType, exists := s.Find("tr.wrtopsection").Attr("data-ph")
	if exists != true {
		return "", errors.New("tr element with class=wrtopsection and attribute data-ph found")
	}
	switch tType {
	case "sMainMeanings":
		return "main", nil
	case "sCmpdForms":
		return "compound", nil
	case "sAddTrans":
		return "supplement", nil
	}


func parseWRDTable(root *goquery.Selection) []*interpretation {
	return root.Find("tr").Each(func(i int, s *goquery.Selection) {

	})
}

// id=articleWRD
func getWRTranslation(word string) {
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
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d", resp.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return parseRoot(doc.Find("table.WRD"))
}
func main() {
	getWRTranslation("coucou")
}
