package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*
		Helpful Links:
		https://golang.org/pkg/net/http/#Request.FormValue
		https://golang.org/pkg/net/http/#Error
		https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	w.Header().Add("Access-Control-Allow-Origin", "*")
	requestQuery := r.URL.Query().Get("url")
	if len(requestQuery) == 0 {
		// w.Write([]byte(strconv.Itoa(http.StatusBadRequest)))
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	stream, err := fetchHTML(requestQuery)
	if err != nil {
		// w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 400)
		return
	}
	summary, err := extractSummary(requestQuery, stream)
	if err != nil {
		// w.Write([]byte(err.Error()))
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.NewEncoder(w).Encode(summary)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, fmt.Errorf("Cannot perform request. Invalid URL")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %v", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("Response is not a valid HTML page. Is %s", resp.Header.Get("Content-Type"))
	}
	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	defer htmlStream.Close()
	tokenizer := html.NewTokenizer(htmlStream)
	summaryData := &PageSummary{}
	image := &PreviewImage{}
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			switch tag := token.Data; tag {
			case "title":
				tokenType = tokenizer.Next()
				if tokenType == html.TextToken {
					if summaryData.Title == "" {
						summaryData.Title = tokenizer.Token().Data
					}
				}
			case "meta":
				// Retrieve basic property/name to value
				mapValues := map[string]string{}
				for _, attr := range token.Attr {
					if attr.Key == "property" || attr.Key == "name" {
						mapValues[attr.Key] = attr.Val
					}
					if attr.Key == "content" {
						mapValues[attr.Key] = attr.Val
					}
				}

				key, exists := mapValues["property"]
				if exists {
					switch key {
					case "og:type":
						value, _ := mapValues["content"]
						summaryData.Type = value
					case "og:url":
						value, _ := mapValues["content"]
						summaryData.URL = value
					case "og:title":
						value, _ := mapValues["content"]
						summaryData.Title = value
					case "og:site_name":
						value, _ := mapValues["content"]
						summaryData.SiteName = value
					case "og:description":
						value, _ := mapValues["content"]
						summaryData.Description = value
					case "og:image":
						if len(image.URL) != 0 {
							summaryData.Images = append(summaryData.Images, image)
							image = &PreviewImage{}
						}
						value, _ := mapValues["content"]
						url := absoluteURL(pageURL, value)
						image.URL = url
					case "og:image:secure_url":
						value, _ := mapValues["content"]
						image.SecureURL = value
					case "og:image:type":
						value, _ := mapValues["content"]
						image.Type = value
					case "og:image:width":
						value, _ := mapValues["content"]
						val, _ := strconv.Atoi(value)
						image.Width = val
					case "og:image:height":
						value, _ := mapValues["content"]
						val, _ := strconv.Atoi(value)
						image.Height = val
					case "og:image:alt":
						value, _ := mapValues["content"]
						image.Alt = value
					}
				}

				key, exists = mapValues["name"]
				if exists {
					switch key {
					case "keywords":
						value, _ := mapValues["content"]
						keywords := strings.Split(value, ",")
						for i := range keywords {
							keywords[i] = strings.TrimSpace(keywords[i])
						}
						summaryData.Keywords = keywords
					case "description":
						value, _ := mapValues["content"]
						if summaryData.Description == "" {
							summaryData.Description = value
						}
					case "author":
						value, _ := mapValues["content"]
						summaryData.Author = value
					}
				}

			case "link":
				previewImage := &PreviewImage{}
				for _, attr := range token.Attr {
					// Make sure it is the icon
					if attr.Key == "rel" && attr.Val == "icon" {

					}
					// Href
					if attr.Key == "href" {
						href := absoluteURL(pageURL, attr.Val)
						previewImage.URL = href
					}
					// Type
					if attr.Key == "type" {
						previewImage.Type = attr.Val
					}
					// Sizes
					if attr.Key == "sizes" {
						sizeString := strings.TrimSpace(attr.Val)
						if sizeString != "any" {
							sizes := strings.Split(attr.Val, "x")
							for i := range sizes {
								sizes[i] = strings.TrimSpace(sizes[i])
							}
							if len(sizes) == 2 {
								height, _ := strconv.Atoi(sizes[0])
								width, _ := strconv.Atoi(sizes[1])
								previewImage.Height = height
								previewImage.Width = width
							}
						}
					}

				}
				summaryData.Icon = previewImage
			}
		}

		// After reaching the end of the head, stop tokenizing
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
		}
	}
	// If there is data in image, append to images array
	if len(image.URL) != 0 {
		summaryData.Images = append(summaryData.Images, image)
	}
	return summaryData, nil
}

func absoluteURL(pageURL string, href string) string {
	if strings.HasPrefix(href, "/") {
		splitURL := strings.Split(pageURL, "/")
		base := strings.Join(splitURL[0:3], "/")
		href = base + href
	}
	return href
}
