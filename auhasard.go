package main

import (
	"fmt"
	//"io"
	"errors"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// A translation is collection of interpretations
type translation struct {
	interpretation map[string]interpretation
}

type interpretation struct {
	from    word
	to      []word
	context []context
	kind    string
}

type word struct {
	name     string
	language string
	pos      string
}

type context struct {
	language string
	phrase   string
}

type WRParser interface {
	parseRoot()
	parseWRDTable()
}

func parseRoot(root *goquery.Selection) *translation {
	root.Each(func(i int, s *goquery.Selection) {
		tType, err := selectTranslationType(s, t)
		if err != nil {
			log.Fatal(err)
		}
		t.interpretation.parseWRDTable(s)
	})

	return t
}

func selectTranslationType(s *goquery.Selection, t *translation) (string, error) {
	tType, exists := s.Find("tr.wrtopsection").Attr("data-ph")
	if exists != true {
		return "", errors.New("tr element with class=wrtopsection and attribute data-ph not found")
	}
	switch tType {
	case "sMainMeanings":
		return "main", nil
	case "sCmpdForms":
		return "compound", nil
	case "sAddTrans":
		return "supplement", nil
	default:
		return "", errors.New("data-ph attribute did not hold one of expected values")
	}
}

func parseWRDTable(root *goquery.Selection, t translation) {
	var interp []interpretation

	//loop through siblings and construct interpretations

	return &interp
}

// id=articleWRD
func getWRTranslation(word string) translation {
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
