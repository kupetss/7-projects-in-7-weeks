package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func Links(node *html.Node) []string {
	var links []string
	var fnc func(*html.Node)
	fnc = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fnc(c)
		}
	}
	fnc(node)
	return links
}

func fetchHTML(urlStr string) (*html.Node, error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return html.Parse(resp.Body)
}

func sameDomain(link *url.URL, base *url.URL) bool {
	return link.Host == base.Host
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Укажите ссылку")
		os.Exit(1)
	}
	urlStr := os.Args[1]
	baseURL, _ := url.Parse(urlStr)
	doc, err := fetchHTML(urlStr)
	if err != nil {
		fmt.Println("Ошибка загрузки страницы")
		os.Exit(1)
	}
	links := Links(doc)
	found := 0
	ch := make(chan string, len(links))
	for _, link := range links {
		go func(link string) {
			if link == "" {
				ch <- ""
				return
			}
			linkURL, err := url.Parse(link)
			if err != nil {
				ch <- ""
				return
			}
			absoluteURL := baseURL.ResolveReference(linkURL)
			if absoluteURL.Scheme == "http" || absoluteURL.Scheme == "https" {
				if sameDomain(absoluteURL, baseURL) {
					ch <- absoluteURL.String()
					return
				}
			}
			ch <- ""
		}(link)
	}

	for i := 0; i < len(links); i++ {
		if result := <-ch; result != "" {
			fmt.Println(result)
			found++
		}
	}
	close(ch)

	if found == 0 {
		fmt.Println("Ссылки не найдены")
	}
	fmt.Println("Всего ссылок:", found)
}
