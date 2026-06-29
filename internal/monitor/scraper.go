package monitor

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	ErrStatusNotOK       = errors.New("the response status is not OK")
	ErrElementNotFound   = errors.New("unable to find the required element")
	ErrAttributeNotFound = errors.New("unable to find the required attribute")
	ErrNoCategory        = errors.New("no category was set")
)

type ReleaseArticle struct {
	Title string
	Url   string
	Books []ReleaseCategory
}

type ReleaseCategory struct {
	Name  string
	Books []string
}

func GetLatestReleaseArticle() (article *ReleaseArticle, err error) {
	res, err := http.Get("https://www.tongli.com.tw/TNews_List.aspx?Type=0&Page=1")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, ErrStatusNotOK
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	articleSel := doc.FindMatcher(goquery.Single(".news_list li:first-child"))
	if articleSel.Length() == 0 {
		return nil, ErrElementNotFound
	}

	titleSel := articleSel.FindMatcher(goquery.Single(".title"))
	linkSel := articleSel.FindMatcher(goquery.Single(".title a"))
	if (titleSel.Length() == 0) || (linkSel.Length() == 0) {
		return nil, ErrElementNotFound
	}

	title := strings.TrimSpace(titleSel.Text())
	link, exists := linkSel.Attr("href")
	if !exists {
		return nil, ErrAttributeNotFound
	}

	linkRef, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	link = res.Request.URL.ResolveReference(linkRef).String()
	return &ReleaseArticle{
		Title: title,
		Url:   link,
	}, nil
}

func GetReleaseBooks(article *ReleaseArticle) (err error) {
	res, err := http.Get(article.Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ErrStatusNotOK
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	contentSel := doc.FindMatcher(goquery.Single("#ContentPlaceHolder1_TNewsContent > div"))
	if contentSel.Length() == 0 {
		return ErrElementNotFound
	}

	books := []ReleaseCategory{}
	categoryIndex := -1

	var contentErr error
	contentSel.Children().EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if len(s.Nodes) == 0 {
			return true
		}

		switch s.Nodes[0].Type {
		case html.ElementNode:
			switch s.Nodes[0].Data {
			case "div":
				categoryName := strings.TrimSpace(s.Text())
				if categoryName == "" {
					return true
				}

				books = append(books, ReleaseCategory{
					Name:  categoryName,
					Books: []string{},
				})
				categoryIndex = len(books) - 1

			case "p":
				if categoryIndex == -1 {
					contentErr = ErrNoCategory
					return false
				}

				s.Contents().Each(func(_ int, child *goquery.Selection) {
					if (len(child.Nodes) == 0) || (child.Nodes[0].Type != html.TextNode) {
						return
					}

					title := strings.TrimSpace(child.Text())
					if title != "" {
						books[categoryIndex].Books = append(books[categoryIndex].Books, title)
					}
				})
			}
		}

		return true
	})
	if contentErr != nil {
		return contentErr
	}

	article.Books = books
	return nil
}
