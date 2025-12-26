package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"errors"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// A translation is collection of interpretations
type translation struct {
	interpretations []*interpretation
}

type interpretation struct {
	from     []*word
	to       []*word
	examples []*example
	kind     string
}

type word struct {
	name       string
	language   string
	pos        string
	definition string
}

type example struct {
	language string
	phrase   string
}

type WRParser interface {
	parseRoot()
	parseWRDTable()
}

func newTranslation() *translation {
	t := translation{
		interpretations: make([]*interpretation, 0),
	}
	return &t
}

func newInterpretation() *interpretation {
	i := interpretation{
		from:     make([]*word, 0),
		to:       make([]*word, 0),
		examples: make([]*example, 0),
	}
	return &i
}

func parseTables(tables *goquery.Selection) *translation {
	t := newTranslation()
	tables.Each(func(i int, table *goquery.Selection) {
		kind, err := selectTranslationKind(table)
		if err != nil {
			log.Fatal(err)
		}
		groups := GroupEvenOddRows(table)
		for _, group := range groups {
			interp := newInterpretation()
			interp.kind = kind
			err := parseTRs(group, interp)
			if err != nil {
				log.Fatal(err)
			}
			t.interpretations = append(t.interpretations, interp)
		}
	})
	return t
}

func GroupEvenOddRows(table *goquery.Selection) []*goquery.Selection {
	var groups []*goquery.Selection
	currentGroup := new(goquery.Selection)

	table.Find("tr").Each(func(i int, tr *goquery.Selection) {
		class, _ := tr.Attr("class")

		if goquery.NodeName(tr) != "tr" {
			return
		}

		if !slices.Contains([]string{"odd", "even"}, class) {
			return
		}
		// Start a new group if none is active
		if currentGroup.Length() == 0 {
			currentGroup = currentGroup.AddSelection(tr)
		} else {
			prevClass, _ := currentGroup.Last().Attr("class")
			if !containsClass(class, prevClass) {
				groups = append(groups, currentGroup)
				currentGroup = new(goquery.Selection)
			}
			currentGroup = currentGroup.AddSelection(tr)
		}
	})
	return groups
}

func containsClass(classAttr, className string) bool {
	for c := range strings.FieldsSeq(classAttr) {
		if c == className {
			return true
		}
	}
	return false
}

func selectTranslationKind(s *goquery.Selection) (string, error) {
	row := s.Find("tr.wrtopsection")
	fmt.Printf("%v\n\n", row.Text())
	kind, exists := row.Find("span.ph").Attr("data-ph")
	row.Remove()
	if !exists {
		return "", errors.New("tr element with class=wrtopsection and attribute data-ph not found")
	}
	switch kind {
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

// Parses even/odd group rows from table
func parseTRs(group *goquery.Selection, interp *interpretation) error {
	group.Each(func(i int, tr *goquery.Selection) {
		html, _ := goquery.OuterHtml(tr)
		fmt.Println(html)
	})
	//group.Find("td[class]")
	re := regexp.MustCompile(`\(([^)]*)\)`)

	group.Each(func(i int, tr *goquery.Selection) {
		tr.Find("tr[id]").Each(func(i int, td *goquery.Selection) {
			em := tr.RemoveFiltered("em")
			def := tr.RemoveFiltered("td:not([class])")
			FrWd = tr.Find("td.FrWd")
			for _, phrase := range strings.Split(FrWd.Text(), ",") {
				interp.from = append(interp.from, &word{
					name:       phrase,
					language:   FrLang,
					pos:        FrPosEM.Text(),
					definition: definitions[0],
				})
			}
		})
		tr.Find("td.FrEx").Each(func(i int, td *goquery.Selection) {
			return
		})
		tr.Find("td.ToEx").Each(func(i int, td *goquery.Selection) {
			return
		})
	})
	FrWdSelection := group.Find("td.FrWd")
	ToWdSelection := group.Find("td.ToWd")
	DefFrSelection := group.Find("td:not([class])")
	FrExSelection := group.Find("td.FrEx")
	ToExSelection := group.Find("td.ToEx")
	FrPosEM := FrWdSelection.RemoveFiltered("em")
	ToPosEM := ToWdSelection.RemoveFiltered("em")
	FrLang := FrPosEM.AttrOr("data-lang", "fr")
	ToLang := ToPosEM.AttrOr("data-lang", "en")

	fmt.Printf("to: %v\nfrom: %v\n", FrWdSelection.Text(), ToWdSelection.Text())

	re := regexp.MustCompile(`\(([^)]*)\)`)

	fromMatches := re.FindAllStringSubmatch(DefFrSelection.Text(), -1)
	//toMatches := re.FindAllStringSubmatch(DefToSelection.Text(), -1)

	var definitions []string
	for _, match := range fromMatches {
		if len(match) > 1 {
			definitions = append(definitions, match[1])
		}
	}

	if len(definitions) < 1 {
		return errors.New("no definitions found")
	}

	if FrWdSelection.Length() > 1 {
		return errors.New("more than one tr.FrWd found")
	}
	{
		// add word to "from" slice
		for _, phrase := range strings.Split(FrWdSelection.Text(), ",") {
			interp.from = append(interp.from, &word{
				name:       phrase,
				language:   FrLang,
				pos:        FrPosEM.Text(),
				definition: definitions[0],
			})
		}
		//add words to "to" slice
		ToWdSelection.Each(func(i int, s *goquery.Selection) {
			for _, phrase := range strings.Split(s.Text(), ",") {
				interp.to = append(interp.to, &word{
					name:       phrase,
					language:   ToLang,
					pos:        ToPosEM.Text(),
					definition: definitions[0],
				})
			}
		})
		//add "from" expressions to examples
		FrExSelection.Each(func(i int, s *goquery.Selection) {
			interp.examples = append(interp.examples, &example{
				language: FrLang,
				phrase:   s.Text(),
			})
		})
		//add "to" expressions to examples
		ToExSelection.Each(func(i int, s *goquery.Selection) {
			interp.examples = append(interp.examples, &example{
				language: ToLang,
				phrase:   s.Text(),
			})
		})

	}
	return nil
}

// id=articleWRD
func getWRFTranslationHTTP(word string) *translation {
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

	return parseTables(doc.Find("table.WRD"))
}

func getWRTRanslationFile(fileName string) {
	// create from a file
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		log.Fatal(err)
	}
	// use the goquery document...
	parseTables(doc.Find("table.WRD"))
}
func main() {

	getWRTRanslationFile("cuire.html")
	//translation := getWRFTranslationHTTP("coucou")
	/*
		for _, interp := range translation.interpretations {
			fmt.Printf("%s translations", interp.kind)
			fmt.Printf("--------------------------------------------")
			for _, from := range interp.from {
				fmt.Printf("%s (%s): %s", from.name, from.pos, from.definition)
			}
			for _, to := range interp.to {
				fmt.Printf("%s (%s): %s", to.name, to.pos, to.definition)
			}
			for _, example := range interp.examples {
				fmt.Printf("%s", example)
			}
		}
	*/
}
