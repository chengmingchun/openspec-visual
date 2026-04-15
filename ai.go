package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type LLMConfig struct {
	APIKey  string `json:"apiKey"`
	BaseURL string `json:"baseUrl"`
	Model   string `json:"model"`
}

var configCache *LLMConfig

func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".openspec-visualizer.json")
}

func LoadLLMConfig() LLMConfig {
	if configCache != nil {
		return *configCache
	}
	path := getConfigPath()
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return LLMConfig{BaseURL: "https://api.openai.com/v1", Model: "gpt-3.5-turbo"}
	}
	var c LLMConfig
	json.Unmarshal(data, &c)
	configCache = &c
	return c
}

func SaveLLMConfig(c LLMConfig) error {
	path := getConfigPath()
	data, _ := json.MarshalIndent(c, "", "  ")
	configCache = &c
	return ioutil.WriteFile(path, data, 0644)
}

func SendPrompt(prompt string, systemPrompt string) (string, error) {
	c := LoadLLMConfig()
	if c.APIKey == "" {
		return "", fmt.Errorf("未配置 API Key，请先进入设置填写。")
	}

	url := c.BaseURL
	if url == "" {
		url = "https://api.openai.com/v1"
	}
	// Trim trailing slash
	if len(url) > 0 && url[len(url)-1] == '/' {
		url = url[:len(url)-1]
	}
	endpoint := url + "/chat/completions"

	reqBody := map[string]interface{}{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API 错误 %d: %s", resp.StatusCode, string(body))
	}

	var res struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return "", fmt.Errorf("解析结果失败: %v", err)
	}

	if len(res.Choices) > 0 {
		return res.Choices[0].Message.Content, nil
	}
	return "", fmt.Errorf("无结果返回")
}
