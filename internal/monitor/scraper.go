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
	Title      string            `json:"title" validate:"required"`
	Url        string            `json:"url" validate:"required,http_url"`
	Categories []ReleaseCategory `json:"categories"`
}

type ReleaseCategory struct {
	Name  string   `json:"name" validate:"required"`
	Books []string `json:"books" validate:"required"`
}

type scraper interface {
	getLatestReleaseArticle() (*ReleaseArticle, error)
	getReleaseBooks(*ReleaseArticle) error
}

type Scraper struct{}

func newScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) getLatestReleaseArticle() (article *ReleaseArticle, err error) {
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

func (s *Scraper) getReleaseBooks(article *ReleaseArticle) (err error) {
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

	categories := []ReleaseCategory{}
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

				categories = append(categories, ReleaseCategory{
					Name:  categoryName,
					Books: []string{},
				})
				categoryIndex = len(categories) - 1

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
						categories[categoryIndex].Books = append(categories[categoryIndex].Books, title)
					}
				})
			}
		}

		return true
	})
	if contentErr != nil {
		return contentErr
	}

	article.Categories = categories
	return nil
}
