package main

import (
	"fmt"
	"strings"
	"regexp"
	"github.com/serpapi/serpapi-golang" 
	"github.com/gocolly/colly/v2"
)

const API_KEY = "SERPAPI_API_KEY" // Change with your API Key

func findSentence(rawText string, searchText string) string {
	
    // 1. Replace all whitespaces with a single space
	re := regexp.MustCompile(`\s+`) 
	fullText := re.ReplaceAllString(rawText, " ")

	// 2. Replace all backtik ’ into ' at rawText
	re1 := regexp.MustCompile(`’`)
	fullText = re1.ReplaceAllString(fullText, "'")

    // 3. Find the start index of searchText
    startIndex := strings.Index(fullText, searchText)
    if startIndex == -1 {
        return "Text not found"
    }

    // 4. Calculate the end index of the snippet
    snippetEndIndex := startIndex + len(searchText)

    // 5. Find the end of the sentence after the snippet
    endOfSentenceIndex := strings.Index(fullText[snippetEndIndex:], ".")
    if endOfSentenceIndex == -1 {
        // Return the rest of the text from snippet if not found
        return fullText[startIndex:]
    }

    // Adjust to get the correct index in the full text
    endOfSentenceIndex += snippetEndIndex + 1
    
    return fullText[startIndex:endOfSentenceIndex]
}

func scrapeFullSnippet(link string, snippet string) string{
	c := colly.NewCollector()

	// Visit a webpage
	rawText := ""
	c.OnHTML("body", func(e *colly.HTMLElement) {
		rawText = e.Text
	})

	// Handle any errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start scraping
	c.Visit(link)

	fullSnippet := findSentence(rawText, snippet)
	return fullSnippet
}

func main() {
	client_parameter := map[string]string{
		"engine": "google",
		"gl": "us",
		"hl": "en",
		"location": "Austin, Texas, United States",
		"api_key": API_KEY,
	}
	client := serpapi.NewClient(client_parameter)

	parameter := map[string]string{ 
		"q": "why the sky is blue", // is donut healthy
		"num": "5",
	}

	data, err := client.Search(parameter)

	// Print each of the link
	type OrganicResult struct {
		Title string
		Snippet string
		Link string
	}

	var organic_results []OrganicResult

	for _, result := range data["organic_results"].([]interface{}) {
		result := result.(map[string]interface{})

		organic_result := OrganicResult{
			Title: result["title"].(string),
			Snippet: result["snippet"].(string),
			Link: result["link"].(string),
		}

		// Check if snippet has ... at the end
		if strings.HasSuffix(organic_result.Snippet, "...") {
			snippet := strings.ReplaceAll(organic_result.Snippet, " ...", "")
			fmt.Println("Checking Snippet :" ,snippet)
			fullSnippet := scrapeFullSnippet(organic_result.Link, snippet)
			organic_result.Snippet = fullSnippet
		}
		
		fmt.Println("--------------")
		fmt.Println(organic_result.Link)
		fmt.Println(organic_result.Snippet)
		fmt.Println("--------------")

		organic_results = append(organic_results, organic_result)
	}
	
	if err != nil {
		fmt.Println(err)
	}
}