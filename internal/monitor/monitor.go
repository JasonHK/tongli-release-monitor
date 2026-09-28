package monitor

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
)

var (
	ErrFailedToGetArticle = errors.New("failed to get release article")
)

type Monitor interface {
	GetReleaseArticle(statesDir string, scraper scraper) (article *ReleaseArticle, updated bool, err error)
}

func GetReleaseArticle(statesDir string, scraper scraper) (article *ReleaseArticle, updated bool, err error) {
	var lastArticle *ReleaseArticle
	lastArticlePath := filepath.Join(statesDir, "article.json")
	if _, err := os.Stat(lastArticlePath); err != nil {
		err := os.MkdirAll(filepath.Dir(statesDir), os.ModePerm)
		if err != nil {
			return nil, false, err
		}

		if _, err := os.Stat(statesDir); err != nil {
			err = os.Mkdir(statesDir, 0o700)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		lastArticleFile, err := os.Open(lastArticlePath)
		if err != nil {
			return nil, false, err
		}
		defer lastArticleFile.Close()

		lastArticleBytes, err := io.ReadAll(lastArticleFile)
		if err != nil {
			return nil, false, err
		}

		err = json.Unmarshal(lastArticleBytes, &lastArticle)
		if err != nil {
			var syntaxErr *json.SyntaxError
			if !errors.As(err, &syntaxErr) {
				return nil, false, err
			}
		}
	}

	article, err = scraper.getLatestReleaseArticle()
	if err != nil {
		return nil, false, err
	}

	err = scraper.getReleaseBooks(article)
	if err != nil {
		return nil, false, err
	}

	lastArticleFile, err := os.Create(lastArticlePath)
	if err != nil {
		return nil, false, err
	}
	defer lastArticleFile.Close()

	lastArticleBytes, err := json.Marshal(article)
	if err != nil {
		return nil, false, err
	}

	_, err = lastArticleFile.Write(lastArticleBytes)
	if err != nil {
		return nil, false, err
	}

	updated = !reflect.DeepEqual(article, lastArticle)
	return article, updated, nil
}
