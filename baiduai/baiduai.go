package baiduai

import (
	"encoding/json"
	"github.com/mlogclub/simple"
	"gopkg.in/resty.v1"
)

func GetTags(title, content string) *AiTags {
	if title == "" || content == "" {
		return nil
	}
	data := make(map[string]interface{})
	data["title"] = title
	data["content"] = simple.Substr(content, 0, 10000)

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/keyword?charset=UTF-8&access_token=" + GetToken()
	response, err := resty.R().SetBody(string(bytesData)).Post(url)
	if err != nil {
		return nil
	}

	tags := &AiTags{}
	err = json.Unmarshal(response.Body(), tags)
	if err != nil {
		return nil
	}
	return tags
}

func GetCategories(title, content string) *AiCategories {
	if title == "" || content == "" {
		return nil
	}

	data := make(map[string]interface{})
	data["title"] = title
	data["content"] = simple.Substr(content, 0, 10000)

	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	url := "https://aip.baidubce.com/rpc/2.0/nlp/v1/topic?charset=UTF-8&access_token=" + GetToken()
	response, err := resty.R().SetBody(string(bytesData)).Post(url)
	if err != nil {
		return nil
	}

	categories := &AiCategories{}
	err = json.Unmarshal(response.Body(), categories)
	if err != nil {
		return nil
	}
	return categories
}
