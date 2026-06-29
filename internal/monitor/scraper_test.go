package monitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestReleaseArticle(t *testing.T) {
	article, err := GetLatestReleaseArticle()
	assert.NotNil(t, article)
	assert.Nil(t, err)
}

func TestGetReleaseBooks(t *testing.T) {
	articleExpected := &ReleaseArticle{
		Title: "6/5《百合大暴走～超心動事件發生中！～》第1集《佐橋小弟的妖怪日常》第1集，新作登場！",
		Url:   "https://www.tongli.com.tw/TNews_View.aspx?Tid=20260605083824",
		Categories: []ReleaseCategory{
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

	article := &ReleaseArticle{
		Title: "6/5《百合大暴走～超心動事件發生中！～》第1集《佐橋小弟的妖怪日常》第1集，新作登場！",
		Url:   "https://www.tongli.com.tw/TNews_View.aspx?Tid=20260605083824",
	}

	GetReleaseBooks(article)
	assert.EqualExportedValues(t, articleExpected, article)

}
