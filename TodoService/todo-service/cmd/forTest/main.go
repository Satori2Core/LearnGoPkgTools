package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Assignee    string `json:"assignee"`
}

var (
	baseURL    = "http://localhost:8080/api"
	users      = []string{"Alice", "Bob", "Charlie", "David", "Eva", "Frank", "Grace", "Henry"}
	taskTitles = []string{
		"修复登录BUG", "优化性能问题", "设计用户资料页", "重构API模块",
		"测试支付流程", "更新文档", "添加新功能", "解决安全漏洞",
	}
)

func main() {
	concurrentUsers := 20
	opsPerUser := 10000

	log.Println("开始压力测试...")
	startTime := time.Now()

	var wg sync.WaitGroup
	wg.Add(concurrentUsers)

	for i := 0; i < concurrentUsers; i++ {
		go func(userID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			taskIDs := make([]int, 0, opsPerUser)

			for j := 0; j < opsPerUser; j++ {
				op := r.Intn(100)

				switch {
				case op < 40: // 40% 创建任务
					taskID := createTask(userID, r)
					if taskID != 0 {
						taskIDs = append(taskIDs, taskID)
					}

				case op < 60: // 20% 开始任务
					if len(taskIDs) > 0 {
						taskID := taskIDs[r.Intn(len(taskIDs))]
						startTask(userID, taskID)
					}

				case op < 80: // 20% 完成任务
					if len(taskIDs) > 0 {
						taskID := taskIDs[r.Intn(len(taskIDs))]
						completeTask(userID, taskID)
					}

				default: // 20% 删除任务
					if len(taskIDs) > 0 {
						taskID := taskIDs[r.Intn(len(taskIDs))]
						deleteTask(userID, taskID)
						// 从本地列表中移除
						for i, id := range taskIDs {
							if id == taskID {
								taskIDs = append(taskIDs[:i], taskIDs[i+1:]...)
								break
							}
						}
					}
				}

				time.Sleep(time.Duration(r.Intn(500)) * time.Millisecond)
				// time.Sleep(time.Duration(time.Second))
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(startTime)
	totalRequests := concurrentUsers * opsPerUser
	qps := float64(totalRequests) / duration.Seconds()

	log.Printf("压力测试完成!")
	log.Printf("总请求数: %d", totalRequests)
	log.Printf("总耗时: %v", duration)
	log.Printf("QPS: %.2f", qps)
}

func createTask(userID int, r *rand.Rand) int {
	userName := users[userID%len(users)]

	task := CreateTaskRequest{
		Title:       taskTitles[r.Intn(len(taskTitles))],
		Description: "这是用户" + userName + "创建的任务",
		Creator:     userName,
		Assignee:    users[r.Intn(len(users))],
	}

	payload, _ := json.Marshal(task)
	resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(payload))

	if err != nil {
		log.Printf("[用户 %d] 创建任务失败: %v", userID, err)
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("[用户 %d] 创建任务失败，状态码: %d", userID, resp.StatusCode)
		return 0
	}

	// 解析返回的任务ID
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("[用户 %d] 解析响应失败: %v", userID, err)
		return 0
	}

	if id, ok := result["id"].(float64); ok {
		return int(id)
	}

	return 0
}

func startTask(userID, taskID int) {
	url := fmt.Sprintf("%s/tasks/%d/start", baseURL, taskID)
	req, _ := http.NewRequest("PATCH", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[用户 %d] 开始任务 %d 失败: %v", userID, taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("[用户 %d] 开始任务 %d 失败，状态码: %d", userID, taskID, resp.StatusCode)
	}
}

func completeTask(userID, taskID int) {
	url := fmt.Sprintf("%s/tasks/%d/complete", baseURL, taskID)
	req, _ := http.NewRequest("PATCH", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[用户 %d] 完成任务 %d 失败: %v", userID, taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("[用户 %d] 完成任务 %d 失败，状态码: %d", userID, taskID, resp.StatusCode)
	}
}

func deleteTask(userID, taskID int) {
	url := fmt.Sprintf("%s/tasks/%d", baseURL, taskID)
	req, _ := http.NewRequest("DELETE", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("[用户 %d] 删除任务 %d 失败: %v", userID, taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		log.Printf("[用户 %d] 删除任务 %d 失败，状态码: %d", userID, taskID, resp.StatusCode)
	}
}
