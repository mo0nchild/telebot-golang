package htmlparser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//HTMLPackageTest is function to test package
func HTMLPackageTest() {
	var inputRequest string
	fmt.Print("Type keyword for search information on wiki page: ")
	fmt.Scanf("%s\n", &inputRequest)

	fmt.Println(WikiParser(inputRequest))
	fmt.Println("You were looking for: " + inputRequest)

	fmt.Println("Current time:", GetWorldTime())
}

//GetWorldTime is function to get world time
func GetWorldTime() string {
	html := HTMLGet("https://time.is/ru/")
	if html == nil {
		return "Не могу открыть веб страницу :("
	}
	return html.Find(".w1").Find("#clock0_bg").Find("#clock").Text()
}

// WikiParser is function for find text information in WiKi page
func WikiParser(itemName string) string {
	// Find the review items

	html := HTMLGet("https://ru.wikipedia.org/wiki/" + itemName)
	if html == nil {
		return "Не могу открыть веб страницу :("
	}

	var parsedOutput []string
	var parsedOutputCleaner []string
	var liOutput []string
	var liOutputCleaner []string

	html.Find(".mw-parser-output").Find("p").Each(func(index int, item *goquery.Selection) {
		title := item.Text()
		parsedOutput = append(parsedOutput, title)
		item.Find("span").Each(func(a int, i *goquery.Selection) {
			parsedOutputCleaner = append(parsedOutputCleaner, i.Text())
		})
		//fmt.Printf("Post #%d: %s - %s\n\n", index, title, parsedOutputCleaner)
	})

	var counter int
	var liCounter []int

	liFinder := func(index int, item *goquery.Selection) {
		item.Find("li").Each(func(a int, i *goquery.Selection) {
			liOutput = append(liOutput, i.Text())
			if a < counter {
				liCounter = append(liCounter, counter)
			}
			counter = a
			i.Find("span").Each(func(a int, i *goquery.Selection) {
				liOutputCleaner = append(liOutputCleaner, i.Text())
			})
		})
	}

	cleanerFunc := func(POC []string, str string) string {
		for _, word := range POC {
			str = strings.ReplaceAll(str, word, "")
		}
		return str
	}

	var wikiAnswer string
	if len(parsedOutput) <= 1 {
		html.Find(".mw-parser-output").Find("ul").Each(liFinder)
		for i := 0; i <= liCounter[0]; i++ {
			wikiAnswer = fmt.Sprintf("%s#%d %s\n", wikiAnswer, i, liOutput[i])
		}
	} else {
		if len(cleanerFunc(parsedOutputCleaner, parsedOutput[0])) <= 10 {
			wikiAnswer = parsedOutput[1]
		} else {
			wikiAnswer = parsedOutput[0]
		}
		wikiAnswer = cleanerFunc(parsedOutputCleaner, wikiAnswer)
	}

	var cleaner string
	for i := 1; i <= 20; i++ {
		cleaner = "[" + strconv.Itoa(i) + "]"
		wikiAnswer = strings.ReplaceAll(wikiAnswer, cleaner, "")
	}

	//S := doc.Find(".mw-parser-output").Find("P").Text()
	//parsedOutput := doc.Find(".mw-parser-output").Find("P").Html()

	//fmt.Println(doc.Find(".mw-parser-output").Find("P").Html())
	//fmt.Printf("ANSWER: %s\n", wikiAnswer)

	return wikiAnswer
}

// HTMLGet is function to get html file
func HTMLGet(requestTopic string) *goquery.Document {
	res, err := http.Get(requestTopic)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error: %d %s\n", res.StatusCode, res.Status)
		return nil
	}
	// Load the HTML document
	HTMLdocument, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	return HTMLdocument
}
