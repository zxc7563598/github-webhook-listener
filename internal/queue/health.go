package queue

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/zxc7563598/github-webhook-listener/internal/config"
	"github.com/zxc7563598/github-webhook-listener/internal/repository"
)

type HealthMonitor struct {
	cfg    config.Config
	client *http.Client
	ctx    context.Context    // 上下文
	cancel context.CancelFunc // 取消函数
}

func NewHealthMonitor(cfg config.Config) *HealthMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthMonitor{
		cfg: cfg,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (m *HealthMonitor) Start() error {
	for reposName, repos := range m.cfg.Repos {
		if repos.HealthCheck != nil {
			healthCheck := repos.HealthCheck
			go m.runSiteChecker(m.ctx, reposName, healthCheck)
		}
	}
	return nil
}

func (m *HealthMonitor) runSiteChecker(
	ctx context.Context,
	reposName string,
	hc *config.HealthCheckConfig,
) {
	if hc.Interval <= 0 {
		log.Printf("[health] %s 未配置监听事件", reposName)
		return
	}

	if hc.Url == "" {
		log.Printf("[health] %s 未配置监听链接", reposName)
		return
	}

	interval := time.Duration(hc.Interval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("[health] %s 健康监控启动, 间隔: %s 秒", reposName, interval)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[health] %s 停止健康监控", reposName)
			return

		case <-ticker.C:
			m.checkOnce(hc, reposName)
		}
	}
}

func (m *HealthMonitor) checkOnce(hc *config.HealthCheckConfig, reposName string) {
	start := time.Now()

	resp, err := m.client.Get(hc.Url)
	if err != nil {
		log.Printf("[health] %s 请求失败: %v", reposName, err)
		if err := repository.HealthMonitoringCreate(reposName, err.Error(), -1, 0); err != nil {
			log.Printf("[database] 健康监控数据库记录失败: %v", err)
		}
		return
	}
	defer resp.Body.Close()

	cost := time.Since(start)

	log.Printf(
		"[health] %s 请求完成，状态:%d 时间:%s",
		reposName,
		resp.StatusCode,
		cost,
	)
	if err := repository.HealthMonitoringCreate(reposName, "", resp.StatusCode, cost.Milliseconds()); err != nil {
		log.Printf("[database] 健康监控数据库记录失败: %v", err)
	}
}

func (m *HealthMonitor) Stop() {
	log.Println("[health] 停止调度器...")
	m.cancel()
	log.Println("[health] 调度器已停止")
}
