package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/zxc7563598/github-webhook-listener/internal/actions"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
)

const (
	// MaxRequestBodySize 最大请求体大小 (10MB)
	MaxRequestBodySize = 10 * 1024 * 1024
)

// HealthHandler 健康检查端点
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func MakeWebhookHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// 只允许POST请求
		if r.Method != http.MethodPost {
			http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
			return
		}

		// 限制请求体大小
		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

		// 读取 body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("[webhook] 读取请求体失败: %v", err)
			http.Error(w, "无法读取正文", http.StatusBadRequest)
			return
		}
		// 提取仓库名
		repoName := extractRepoName(body)
		if repoName == "" {
			log.Printf("[webhook] 无法从请求中提取仓库名")
			http.Error(w, "无法检测存储库", http.StatusBadRequest)
			return
		}

		// 查找配置
		repoCfg, ok := cfg.Repos[repoName]
		if !ok {
			log.Printf("[webhook] 在配置中找不到存储库: %s", repoName)
			http.Error(w, "未配置存储库", http.StatusNotFound)
			return
		}
		// 签名校验
		sig := r.Header.Get("X-Hub-Signature-256")
		if !ValidateGitHubSignature(repoCfg.Secret, body, sig) {
			log.Printf("[webhook] 仓库 %s 的GitHub签名验证失败", repoName)
			http.Error(w, "无效签名", http.StatusUnauthorized)
			return
		}
		// 事件类型
		event := r.Header.Get("X-GitHub-Event")
		if event == "" {
			log.Printf("[webhook] 缺少 X-GitHub-Event 头")
			http.Error(w, "不存在 X-GitHub-Event", http.StatusBadRequest)
			return
		}
		// 分支
		branch := extractBranch(body)
		log.Printf("[webhook] 仓库: %s, 事件: %s, 分支: %s", repoName, event, branch)
		// 遍历规则匹配
		for _, rule := range repoCfg.Rules {
			if rule.Event != event {
				continue
			}
			// 分支匹配
			if len(rule.Branches) > 0 {
				match := false
				for _, b := range rule.Branches {
					if b == branch {
						match = true
						break
					}
				}
				if !match {
					continue
				}
			}
			log.Printf("[webhook] 仓库 %s 的规则匹配: event=%s, branch=%s", repoName, rule.Event, branch)
			actions.ExecuteActions(rule.Actions, body)
		}

		// 记录处理时间
		duration := time.Since(startTime)
		log.Printf("[webhook] 请求处理完成，耗时: %v", duration)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

// extractRepoName 从 payload 获取 repository.full_name
func extractRepoName(body []byte) string {
	var obj struct {
		Repository struct {
			FullName string `json:"full_name"`
		} `json:"repository"`
	}

	if err := json.Unmarshal(body, &obj); err != nil {
		return ""
	}

	return obj.Repository.FullName
}

// extractBranch: 从 payload 里面获取 ref，如 "refs/heads/main"
func extractBranch(body []byte) string {
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return ""
	}
	ref, ok := obj["ref"].(string)
	if !ok {
		return ""
	}
	// refs/heads/main → main
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}
