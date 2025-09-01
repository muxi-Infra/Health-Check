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
	TestUrl = "http://127.0.0.1:3000/send_group_msg" // 测试用的 URL
	ProdURL = "http://127.0.0.1:3000/send_group_msg" // 生产用的 URL
	GroupID = int64(740450762)
)

// 发送请求的函数
func sendRequests(acf AppConfig) (HealthCheckResponse, error) {
	url := fmt.Sprintf("http://%s:%d%s", acf.Domain, acf.Port, acf.Path)
	req, err := http.NewRequest(acf.Type, url, nil)
	if err != nil {
		fmt.Printf("创建请求失败: %v", err)
		return HealthCheckResponse{}, fmt.Errorf("创建请求失败: %v", err)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v", err)
		return HealthCheckResponse{}, fmt.Errorf("请求失败: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("请求返回非200状态码: %d", resp.StatusCode)
		return HealthCheckResponse{}, fmt.Errorf("请求返回非200状态码: %d", resp.StatusCode)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// 读取响应
	var result HealthCheckResponse
	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return HealthCheckResponse{}, fmt.Errorf("解析响应JSON失败: %v", err)
	}
	return result, nil
}

func checkApp(name string, acf AppConfig) {
	var message []MessageItem
	resp, err := Retry(func() (HealthCheckResponse, error) {
		return sendRequests(acf)
	}, 3, 500*time.Microsecond)

	if err != nil {
		// 发送告警消息，@全体成员
		message = []MessageItem{
			{
				Type: "at",
				Data: map[string]interface{}{
					"qq": "all",
				},
			},
			{
				Type: "text",
				Data: map[string]interface{}{
					"text": fmt.Sprintf("应用 %s 出现故障，请检查！\n错误原因: %s", name, err),
				},
			},
		}
	} else {
		// 发送正常信息，无需@全体成员
		message = []MessageItem{
			{
				Type: "text",
				Data: map[string]interface{}{
					"text": fmt.Sprintf(
						"应用 %s 正常！\n"+
							"系统CPU使用率: %s\n"+
							"系统内存使用率: %s\n"+
							"系统磁盘使用量: %s\n"+
							"进程CPU使用率: %s\n"+
							"进程常驻内存: %s MB\n"+
							"当前Goroutine数: %s\n"+
							"Go堆内存分配: %s MB\n",
						name, resp.System.CPUPercent, resp.System.MemoryPercent, resp.System.DiskPercent,
						resp.Process.CPUPercent, resp.Process.MemoryRSSMB, resp.Process.Goroutines, resp.Process.GoHeapAllocMB),
				},
			},
		}
	}
	go SendGroupMessage(TestUrl, GroupID, message)
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
		go KeepBeat(func() (interface{}, error) {
			checkApp(name, acf)
			return nil, nil
		}, 45*time.Minute, 12*time.Hour)
	}
	select {}
}
