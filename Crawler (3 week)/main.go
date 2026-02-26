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

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func sameDomain(link *url.URL, base *url.URL) bool {
	return link.Host == base.Host
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("укажмте ссылку")
		os.Exit(1)
	}

	urlStr := os.Args[1]
	url, _ := url.Parse(urlStr)

	doc, err := fetchHTML(urlStr)
	if err != nil {
		fmt.Println("Ошибка загрузки страницы")
		os.Exit(1)
	}

	links := Links(doc)
	found := 0

	for _, link := range links {
		if link == "" {
			continue
		}

		linkUrl, _ := url.Parse(link)

		if linkUrl.Scheme == "http" || linkUrl.Scheme == "https" {
			if sameDomain(linkUrl, url) {
				fmt.Println(linkUrl.String())
				found++
			}
		}
	}

	if found == 0 {
		fmt.Println("ccылки не найдены")
	}
	fmt.Println("Всего ссылок: ", found)
}
