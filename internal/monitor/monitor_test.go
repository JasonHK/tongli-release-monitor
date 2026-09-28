package monitor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type juneFifthScraper struct {
	scraper
	article    *ReleaseArticle
	categories []ReleaseCategory
}

func (s juneFifthScraper) getLatestReleaseArticle() (*ReleaseArticle, error) {
	article := *s.article
	return &article, nil
}

func (s juneFifthScraper) getReleaseBooks(article *ReleaseArticle) error {
	article.Categories = s.categories
	return nil
}

func setupJuneFifthScraper() scraper {
	return &juneFifthScraper{
		scraper: NewScraper(),
		article: &ReleaseArticle{
			Title: "6/5《百合大暴走～超心動事件發生中！～》第1集《佐橋小弟的妖怪日常》第1集，新作登場！",
			Url:   "https://www.tongli.com.tw/TNews_View.aspx?Tid=20260605083824",
		},
		categories: []ReleaseCategory{
			{
				Name: "新書",
				Books: []string{
					"《妖精森林的小不點》第12集",
					"《佐橋小弟的妖怪日常》第1集",
					"《百合大暴走～超心動事件發生中！～》第1集 首刷限定版(附-畫卡+小冊子)",
					"《WITCH WATCH魔女守護者》第24集",
					"《朱音落語》第20集(附-序號卡)",
				},
			},
			{
				Name: "電子書",
				Books: []string{
					"《啟示錄四騎士》第22集",
					"《結緣甘神神社》第9集",
					"《不時輕聲地以俄語遮羞的鄰座艾莉同學》第8集",
					"《東京外星人》第1-5集",
					"連載《中華小廚師！極》",
					"連載《啟示錄四騎士》",
				},
			},
			{
				Name: "再版",
				Books: []string{
					"《一騎當千 愛藏版》第1-2集",
					"《BLUE LOCK藍色監獄角色書 EGOIST BIBLE》第1-2集",
					"《渡良瀨明明就是渣男》第1集",
				},
			},
		},
	}
}

type juneEighthScraper struct {
	scraper
	article    *ReleaseArticle
	categories []ReleaseCategory
}

func (s juneEighthScraper) getLatestReleaseArticle() (*ReleaseArticle, error) {
	article := *s.article
	return &article, nil
}

func (s juneEighthScraper) getReleaseBooks(article *ReleaseArticle) error {
	article.Categories = s.categories
	return nil
}

func setupJuneEighthScraper() scraper {
	return &juneEighthScraper{
		scraper: NewScraper(),
		article: &ReleaseArticle{
			Title: "6/8《精靈渡邊》第1集《伊藤潤二自選傑作集III Dark Colors》全，精彩登場！",
			Url:   "https://www.tongli.com.tw/TNews_View.aspx?Tid=20260605083824",
		},
		categories: []ReleaseCategory{
			{
				Name: "新書",
				Books: []string{
					"《伊藤潤二自選傑作集III Dark Colors》全 首刷限定版(附-光柵卡+相卡2入+畫板)",
					"《魔球投手 愛藏版》第1-5集完 首刷書盒版(附-收藏卡5入)",
					"《魔物娘的同居日常 4格同人合集》第5集",
					"《BLACK CAT 黑貓 愛藏版》第10集",
					"《薰香花朵凛然綻放》第19集",
				},
			},
			{
				Name: "小說",
				Books: []string{
					"《精靈渡邊》第1集 首刷限定版(附-書卡+小冊子)",
					"《彈珠汽水瓶裡的千歲同學9.5》第9.5集 首刷限定版(附-書卡+貼紙組+小冊子)",
				},
			},
			{
				Name: "電子書",
				Books: []string{
					"《派對咖孔明》第22集",
					"《薰香花朵凛然綻放》第19集",
					"《瑠璃龍龍》第4集",
					"《彈珠汽水瓶裡的千歲同學》第3集",
					"連載《SPY×FAMILY間諜家家酒》",
					"連載《派對咖孔明》",
					"連載《八角籠社畜》",
					"連載《社長、酒與星星》",
					"連載《贈予你的魔法》",
					"連載《學生會也有洞！》",
					"連載《人偶：魔女與咒的碎片》",
				},
			},
			{
				Name: "雜誌",
				Books: []string{
					"『寶島少年』第28期",
				},
			},
			{
				Name: "再版",
				Books: []string{
					"《冰之城牆》第1-6集",
				},
			},
		},
	}
}

func TestGetReleaseArticleFirstRun(t *testing.T) {
	statesDir := t.TempDir()

	scraper := setupJuneFifthScraper()
	article, updated, err := GetReleaseArticle(statesDir, scraper)

	assert.NotNil(t, article)
	assert.True(t, updated)
	assert.Nil(t, err)
}

func TestGetReleaseArticleNoUpdates(t *testing.T) {
	statesDir := t.TempDir()

	scraper := setupJuneFifthScraper()
	GetReleaseArticle(statesDir, scraper)

	article, updated, err := GetReleaseArticle(statesDir, scraper)
	assert.NotNil(t, article)
	assert.False(t, updated)
	assert.Nil(t, err)
}

func TestGetReleaseArticleUpdatedData(t *testing.T) {
	statesDir := t.TempDir()

	originalScraper := setupJuneFifthScraper()
	_, _, err := GetReleaseArticle(statesDir, originalScraper)
	assert.Nil(t, err)

	updatedScraper := setupJuneEighthScraper()
	article, updated, err := GetReleaseArticle(statesDir, updatedScraper)

	assert.Equal(t, "6/8《精靈渡邊》第1集《伊藤潤二自選傑作集III Dark Colors》全，精彩登場！", article.Title)
	assert.True(t, updated)
	assert.Nil(t, err)
}

func TestGetReleaseArticleDirNotExist(t *testing.T) {
	statesDir := filepath.Join(t.TempDir(), "non-exist")

	scraper := setupJuneFifthScraper()
	article, updated, err := GetReleaseArticle(statesDir, scraper)

	assert.NotNil(t, article)
	assert.True(t, updated)
	assert.Nil(t, err)
}

func TestGetReleaseArticleInvalidArticleJson(t *testing.T) {
	statesDir := t.TempDir()

	file, err := os.Create(filepath.Join(statesDir, "article.json"))
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString("NOT JSON")
	if err != nil {
		panic(err)
	}

	scraper := setupJuneFifthScraper()
	article, updated, err := GetReleaseArticle(statesDir, scraper)

	assert.NotNil(t, article)
	assert.True(t, updated)
	assert.Nil(t, err)
}
