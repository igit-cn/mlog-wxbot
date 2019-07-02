package wxbot

import (
	"fmt"
	"github.com/mlogclub/mlog-wxbot/config"
	"image"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mlogclub/mlog-wxbot/baiduai"

	"github.com/mlogclub/simple"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"gopkg.in/resty.v1"
)

type WxArticle struct {
	Id          int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
	Title       string `gorm:"type:longtext" json:"title"`              // 标题
	Author      string `gorm:"type:longtext" json:"author"`             // 作者
	AppName     string `gorm:"type:longtext" json:"appName"`            // 公众号名称
	AppID       string `gorm:"type:longtext" json:"appId"`              // 公众号ID
	Cover       string `gorm:"type:longtext" json:"cover"`              // 文章封面
	Intro       string `gorm:"type:longtext" json:"intro"`              // 描述
	HtmlContent string `gorm:"type:longtext" json:"htmlContent"`        // 公众号文章html内容
	MdContent   string `gorm:"type:longtext" json:"mdContent"`          // 公众号文章md内容
	TextContent string `gorm:"type:longtext" json:"textContent"`        // 文本内容
	PubAt       string `gorm:"type:longtext" json:"pubAt"`              // 发布时间
	UrlMd5      string `gorm:"size:64;index:idx_url_md5" json:"urlMd5"` // 链接地址的md5
	RoundHead   string `gorm:"type:longtext" json:"roundHead"`          // 圆头像
	OriHead     string `gorm:"type:longtext" json:"oriHead"`            // 原头像
	Url         string `gorm:"type:text" json:"url"`                    // 微信文章链接地址
	SourceURL   string `gorm:"type:text" json:"sourceUrl"`              // 公众号原文地址
	ArticleId   int64  `json:"articleId"`                               // 发布线上返回的id
	Tags        string `gorm:"type:longtext" json:"tags"`               // 标签字符串
	Category    string `gorm:"type:longtext" json:"category"`           // 一级分类
	Categories  string `gorm:"type:longtext" json:"categories"`         // 二级分类
	Copyright   string `gorm:"type:longtext" json:"copyright"`          // 已经 0,1,2   微小宝那 1 标识为原创
	Video       string `gorm:"type:longtext" json:"video"`              // 视频地址
	Audio       string `gorm:"type:longtext" json:"audio"`              // 音频地址
	WxID        string `gorm:"type:longtext" json:"wxId"`               // 微信公众号ID
	WxIntro     string `gorm:"type:longtext" json:"wxIntro"`            // 微信公众号介绍
	Images      string `gorm:"type:longtext" json:"images"`             // 图片
	PublishTime int64  `json:"publishTime"`                             // 采集器发布时间
	CreateTime  int64  `json:"createTime" `
	UpdatedTime int64  `json:"updatedTime"`
}

type WxArticleCallback func(article WxArticle)

// 抓取文章
func collect(url string) *WxArticle {
	article, err := collectArticle(url)
	if err != nil {
		logrus.Error("抓取文章失败：" + url)
		return nil
	}
	if len(article.URL) == 0 {
		logrus.Error("不支持该链接")
		return nil
	}

	var wxArticle WxArticle
	wxArticle.AppID = article.AppID
	wxArticle.AppName = article.AppName
	wxArticle.RoundHead = article.RoundHead
	wxArticle.OriHead = article.OriHead
	wxArticle.Url = article.URL
	wxArticle.SourceURL = article.SourceURL
	wxArticle.UrlMd5 = simple.MD5(article.URL)
	wxArticle.Title = article.Title
	wxArticle.Intro = article.Intro
	wxArticle.HtmlContent = article.HtmlContent
	wxArticle.MdContent = article.MdContent
	wxArticle.TextContent = article.TextContent
	wxArticle.Cover = article.Cover
	wxArticle.Author = article.Author
	wxArticle.PubAt = article.PubAt
	wxArticle.Copyright = article.Copyright
	wxArticle.WxID = article.WxID
	wxArticle.WxIntro = article.WxIntro
	wxArticle.Video = article.Video
	wxArticle.Audio = article.Audio
	wxArticle.Category = `其它`

	categories := baiduai.GetCategories(wxArticle.Title, wxArticle.TextContent)
	if categories != nil {
		// 一级分类
		var topCategories []string
		for _, t := range categories.Item.TopCategory {
			topCategories = append(topCategories, t.Tag)
		}
		if len(topCategories) > 0 {
			wxArticle.Category = strings.Join(topCategories, ",")
		}

		// 二级分类
		var secondCategories []string
		for _, t := range categories.Item.SecondCatrgory {
			secondCategories = append(secondCategories, t.Tag)
		}
		if len(secondCategories) > 0 {
			wxArticle.Categories = strings.Join(secondCategories, ",")
		}
	}

	tags := baiduai.GetTags(wxArticle.Title, wxArticle.TextContent)
	if tags != nil {
		var tagArr []string
		for _, t := range tags.Items {
			tagArr = append(tagArr, t.Tag)
		}
		if len(tagArr) > 0 {
			wxArticle.Tags = strings.Join(tagArr, ",")
		}
	}

	if len(article.Images) > 0 {
		var imgArr []string
		for _, img := range article.Images {
			if checkImage(img) {
				imgArr = append(imgArr, img)
			}
		}
		if len(imgArr) > 0 {
			wxArticle.Images = strings.Join(imgArr, ";")
		}
	}

	// // 增加音频标签
	// if wxArticle.Audio != `` {
	// 	if wxArticle.Tags != `` {
	// 		wxArticle.Tags = fmt.Sprintf(`%v,音频`, wxArticle.Tags)
	// 	} else {
	// 		wxArticle.Tags = `音频`
	// 	}
	// }
	// // 视频标签
	// if wxArticle.Video != `` {
	// 	if wxArticle.Tags != `` {
	// 		wxArticle.Tags = fmt.Sprintf(`%v,视频`, wxArticle.Tags)
	// 	} else {
	// 		wxArticle.Tags = `视频`
	// 	}
	// }

	logrus.Println("collector", wxArticle.Id, wxArticle.Title, wxArticle.Url, wxArticle.Category, wxArticle.Categories, wxArticle.Tags)
	return &wxArticle
}

// 文章保存起来
func save(article *WxArticle) *WxArticle {
	err := simple.GetDB().Create(article).Error
	if err != nil {
		return nil
	}
	return article
}

// 发布
func publish(article *WxArticle) {
	json, err := simple.FormatJson(article)
	if err != nil {
		logrus.Error(err)
		return
	}
	response, err := resty.SetTimeout(time.Second * 3).R().
		SetBody(json).
		Post(config.Conf.PublishApi + "?token=" + config.Conf.PublishToken)
	if err != nil {
		logrus.Error(err)
		return
	}
	ret := gjson.Get(string(response.Body()), "data.id")

	articleId := ret.Int()
	if articleId > 0 {
		logrus.Info("成功发布文章..." + strconv.FormatInt(articleId, 10))
		simple.GetDB().Model(&WxArticle{}).Where("id = ?", article.Id).Updates(map[string]interface{}{
			"article_id":   articleId,
			"publish_time": simple.NowTimestamp(),
		})
	} else {
		logrus.Info("文章发布失败..." + string(response.Body()))
	}
}

// checkImage 检查图片合法性，宽高大于或等于320
func checkImage(imageURL string) bool {
	resp, err := http.Get(imageURL)
	if err != nil {
		return false
	}
	c, _, err := image.DecodeConfig(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return false
	}
	width := c.Width
	height := c.Height
	if width >= 320 && height >= 320 {
		return true
	}
	return false
}

// 启动的时候将没发布的文章重新发布下
func PublishOnStart() {
	scan(func(article WxArticle) {
		if article.ArticleId <= 0 {
			publish(&article)
		}
	})
}

func scan(callback WxArticleCallback) {
	var cursor int64 = 0
	for {
		var articles []WxArticle
		simple.GetDB().Where("id > ?", cursor).Order("id asc").Limit(3000).Find(&articles)
		if len(articles) == 0 {
			break
		}
		fmt.Println("scan login account cursor...", cursor)
		for _, article := range articles {
			cursor = article.Id
			callback(article)
		}
	}
}
