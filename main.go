package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

// 发送请求的函数
func sendRequests(acf AppConfig) interface{} {
	url := fmt.Sprintf("http://%s:%d%s", acf.Domain, acf.Port, acf.Path)
	req, err := http.NewRequest(acf.Type, url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v", err)
		return nil
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v", err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("请求成功 [%d]: %s", resp.StatusCode, string(body))
	return body
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
			sendRequests(acf)
		}(name, acf)
	}
	time.Sleep(3 * time.Second)

	//// 每 30 分钟执行一次
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
