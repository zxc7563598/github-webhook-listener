package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxc7563598/github-webhook-listener/internal/bootstrap"
	"github.com/zxc7563598/github-webhook-listener/internal/config"
)

func main() {
	port := flag.Int("port", 9000, "服务器端口")
	web := flag.Bool("web", false, "是否开启web")
	workers := flag.Int("workers", 5, "Shell 任务最大并发数")
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	webUser := flag.String("user", "", "Web UI Basic Auth 用户名")
	webPass := flag.String("pass", "", "Web UI Basic Auth 密码")
	flag.Parse()

	// 初始化数据库
	db := config.InitSQLiteDB()

	// 加载配置
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("未能加载配置: %v", err)
	}

	r, scheduler, healthMonitor := bootstrap.NewApp(db, *web, *workers, cfg, *webUser, *webPass)

	addr := fmt.Sprintf(":%d", *port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		log.Printf("Webhook 监听已经启动在 %s 端口", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")
	healthMonitor.Stop()
	scheduler.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭: %v", err)
	}
	log.Println("服务器已优雅关闭")
}
