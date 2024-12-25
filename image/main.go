package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	Url   = "https://ark.cn-beijing.volces.com/api/v3/chat/completions"
	Model = "ep-*"
	Key   = "8429f8ab-*"
)

type Response struct {
	Choices []Choice `json:"choices"`
	Created int64    `json:"created"`
	Id      string   `json:"id"`
	Model   string   `json:"model"`
	Object  string   `json:"object"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	FinishReason string  `json:"finish_reason"`
	Index        int     `json:"index"`
	LogProbs     bool    `json:"log_probs"`
	Message      Message `json:"message"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Usage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func main() {
	var res Response

	buf := fmt.Sprintf(`{
		"model": "%s",
		"messages": [
			{
				"role": "user",
				"content": [
					{
						"type": "text",
						"text": "图片主要讲了什么?"
					},
					{
						"type": "image_url",
						"image_url": {
							"url": "https://ark-project.tos-cn-beijing.volces.com/images/view.jpeg"
						}
					},
					{
						"type": "image_url",
						"image_url": {
							"url": "https://portal.volccdn.com/obj/volcfe/cloud-universal-doc/upload_a81e3cdd3e30617ecd524a132fdb2736.png"
						}
					}
				]
			}
		]
	}`, Model)

	req, err := http.NewRequest("POST", Url, bytes.NewBuffer([]byte(buf)))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	fmt.Printf("Response status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := json.Unmarshal((body), &res); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Response body: %v\n", res)
}
