package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// GroupMessageRequest 定义群消息请求结构体
type GroupMessageRequest struct {
	GroupID int64         `json:"group_id"`
	Message []MessageItem `json:"message"`
}

// MessageItem 定义消息项结构体
type MessageItem struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// SendGroupMessage 发送群消息的封装函数
func SendGroupMessage(url string, groupID int64, message []MessageItem) (interface{}, error) {
	// 创建请求体
	requestBody := GroupMessageRequest{
		GroupID: groupID,
		Message: message,
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	var result map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %v", err)
	}
	return result, nil
}
