package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
		type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	// TODO: fmt.Print()
	fmt.Println("HELLO")
	fmt.Printf("Starting method %s", "test")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	// Get the url query string parameter?
	requestQuery := r.URL.Query().Get("url")
	if len(requestQuery) == 0 {
		w.Write([]byte(strconv.Itoa(http.StatusBadRequest)))
	}
	log.Print("Made it")
	stream, err := fetchHTML(requestQuery)
	if err != nil {
		w.Write([]byte("1"))
	}
	summary, err := extractSummary(requestQuery, stream)
	if err != nil {
		w.Write([]byte("2"))
	}

	err = json.NewEncoder(w).Encode(summary)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	resp, err := http.Get(pageURL)
	if err != nil {
		log.Fatal(err)
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
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.

	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary

	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/
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
			log.Fatalf("error tokenizing HTML: %v", tokenizer.Err())
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
