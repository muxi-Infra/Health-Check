package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

var (
	BaseURL = "http://127.0.0.1:5700/send_group_msg"
	GroupID = int64(123456789)
)

// 发送请求的函数
func sendRequests(acf AppConfig) (interface{}, error) {
	url := fmt.Sprintf("http://%s:%d%s", acf.Domain, acf.Port, acf.Path)
	req, err := http.NewRequest(acf.Type, url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v", err)
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v", err)
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求返回非200状态码: %d", resp.StatusCode)
		return nil, fmt.Errorf("请求返回非200状态码: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// 读取响应
	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %v", err)
	}
	return result, nil
}

func main() {
	v := viper.New()
	v.SetConfigFile("config.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("读取配置失败: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("解析配置失败: %v", err)
	}

	for name, app := range cfg.App {
		fmt.Printf("name=%s, domain=%s, type=%s, port=%d, path=%s\n",
			name, app.Domain, app.Type, app.Port, app.Path)
	}

	for name, acf := range cfg.App {
		go func(name string, acf AppConfig) {
			fmt.Printf("启动应用: %s", name)
			data, err := sendRequests(acf)
			if err != nil {
				fmt.Printf("应用 %s 请求失败: %v\n", name, err)
			} else {
				fmt.Printf("应用 %s 请求成功: %v\n", name, data)
				go func() {
					_, err := SendGroupMessage(BaseURL, GroupID, []MessageItem{})
					if err != nil {
						fmt.Printf("发送群消息失败: %v\n", err)
					} else {
						fmt.Printf("发送群消息成功\n")
					}
				}()
			}
		}(name, acf)
	}
	time.Sleep(3 * time.Second)

	// 每 30 分钟执行一次
	//ticker := time.NewTicker(30 * time.Minute)
	//defer ticker.Stop()
	//
	//for {
	//	select {
	//	case <-ticker.C:
	//		sendRequests(cfg)
	//	}
	//}
}
