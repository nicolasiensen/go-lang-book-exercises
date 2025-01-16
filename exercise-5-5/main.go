package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// !+
func main() {
	for _, url := range os.Args[1:] {
		words, images, err := CountWordsAndImages(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CountWordsAndImages: %v\n", err)
			continue
		}

		fmt.Printf("%s\t%d\t%d\n", url, words, images)
	}
}

func visitWordsAndImages(words, images int, n *html.Node) (int, int) {
	if n.Type == html.ElementNode && n.Data == "img" {
		images++
	} else if n.Type == html.TextNode {
		words += len(strings.Split(n.Data, " "))
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if n.Data != "script" {
			words, images = visitWordsAndImages(words, images, c)
		}
	}

	return words, images
}

// CountWordsAndImages does an HTTP GET request for the html
// document url and returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}

	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	return visitWordsAndImages(0, 0, n)
}
