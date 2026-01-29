package queue

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/zxc7563598/github-webhook-listener/pkg/utils"
)

// ShellTaskStatus 任务状态
type ShellTaskStatus int

const (
	StatusPending   ShellTaskStatus = iota // 0
	StatusRunning                          // 1
	StatusSuccess                          // 2
	StatusFailed                           // 3
	StatusTimeout                          // 4
	StatusCancelled                        // 5
)

func (s ShellTaskStatus) String() string {
	switch s {
	case StatusPending:
		return "待处理"
	case StatusRunning:
		return "运行中"
	case StatusSuccess:
		return "成功"
	case StatusFailed:
		return "失败"
	case StatusTimeout:
		return "超时"
	case StatusCancelled:
		return "取消"
	default:
		return "未知"
	}
}

// ShellTask 任务定义
type ShellTask struct {
	ID         string        // 任务ID
	Name       string        // 任务名称
	Cmd        string        // 执行命令
	Args       []string      // 命令参数
	Timeout    time.Duration // 超时时间
	RetryCount int           // 重试次数
	RetryDelay time.Duration // 重试延迟
	Env        []string      // 环境变量
	WorkDir    string        // 工作目录
}

// ShellTaskResult 任务执行结果
type ShellTaskResult struct {
	TaskID        string          // 任务ID
	TaskName      string          // 任务名称
	Status        ShellTaskStatus // 任务状态
	StartTime     time.Time       // 开始时间
	EndTime       time.Time       // 结束时间
	Duration      time.Duration   // 持续时间
	ExitCode      int             // 退出code
	StdoutLogPath string          // Stdout输出路径
	StderrLogPath string          // Stderr输出路径
	Error         error           // 错误信息
	RetryCount    int             // 重试次数
}

// Scheduler 调度器
type ShellScheduler struct {
	maxWorkers      int                         // 最大并发数
	tasks           map[string]*ShellTask       // 所有任务
	taskResults     map[string]*ShellTaskResult // 所有任务结果
	taskQueue       chan *ShellTask             // 任务队列
	taskResultQueue chan *ShellTaskResult       // 结果队列
	wg              sync.WaitGroup              // 等待组
	mu              sync.Mutex                  // 读写锁
	ctx             context.Context             // 上下文
	cancel          context.CancelFunc          // 取消函数
	isRunning       bool                        // 是否正在运行
	completedTasks  map[string]bool             // 已完成任务
}

// ShellNewScheduler 创建调度器
func ShellNewScheduler(maxWorkers int) *ShellScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &ShellScheduler{
		maxWorkers:      maxWorkers,
		tasks:           make(map[string]*ShellTask),
		taskResults:     make(map[string]*ShellTaskResult),
		taskQueue:       make(chan *ShellTask, 100),
		taskResultQueue: make(chan *ShellTaskResult, 100),
		ctx:             ctx,
		cancel:          cancel,
		completedTasks:  make(map[string]bool),
	}
}

// AddTask 添加任务
func (s *ShellScheduler) AddTask(task *ShellTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if task.ID == "" {
		task.ID = fmt.Sprintf("task-%d", len(s.tasks)+1)
	}
	if task.Name == "" {
		task.Name = task.ID
	}
	if task.Timeout == 0 {
		task.Timeout = 300 * time.Second
	}
	s.tasks[task.ID] = task
	s.taskQueue <- task
	return nil
}

// copyAndLogToFile 复制并同时写入文件
func (s *ShellScheduler) copyAndLogToFile(dst io.Writer, src io.Reader, taskName, filePath string) error {
	// 创建或打开文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(src)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// 1. 写入原始目标（如标准输出）
		if dst != nil {
			io.WriteString(dst, line+"\n")
		}

		// 3. 写入文件
		if _, err := file.WriteString(line + "\n"); err != nil {
			log.Printf("[queue] 写入文件失败: %v", err)
		}
	}

	// 检查Scanner错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取数据失败: %w", err)
	}

	log.Printf("[queue][%s] 已写入文件: %s，共 %d 行", taskName, filePath, lineNumber)
	return nil
}

// runCommand 执行shell命令
func (s *ShellScheduler) runCommand(task *ShellTask, output io.Writer, stdoutLogPath, stderrLogPath string) (int, error) {
	ctx, cancel := context.WithTimeout(s.ctx, task.Timeout)
	defer cancel()
	// 创建命令
	var cmd *exec.Cmd
	if len(task.Args) > 0 {
		cmd = exec.CommandContext(ctx, task.Cmd, task.Args...)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", task.Cmd)
	}

	// 设置工作目录
	if task.WorkDir != "" {
		cmd.Dir = task.WorkDir
	}

	// 设置环境变量
	if len(task.Env) > 0 {
		cmd.Env = append(os.Environ(), task.Env...)
	}

	// 设置输出
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return -1, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return -1, err
	}

	// 并发读取 stdout 和 stderr
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.copyAndLogToFile(output, stdoutPipe, task.Name, stdoutLogPath)
	}()
	go func() {
		defer wg.Done()
		s.copyAndLogToFile(output, stderrPipe, task.Name, stderrLogPath)
	}()
	wg.Wait()

	// 等待命令完成
	err = cmd.Wait()
	exitCode := cmd.ProcessState.ExitCode()
	if ctx.Err() == context.DeadlineExceeded {
		return exitCode, fmt.Errorf("任务执行超时(限时: %v)", task.Timeout)
	}
	return exitCode, err
}

// executeTask 执行单个任务
func (s *ShellScheduler) executeTask(workerID int, task *ShellTask) *ShellTaskResult {
	stdoutLog, _ := utils.GenerateRequestLogPath()
	stderrLog, _ := utils.GenerateRequestLogPath()

	result := &ShellTaskResult{
		TaskID:        task.ID,
		TaskName:      task.Name,
		Status:        StatusRunning,
		StartTime:     time.Now(),
		StdoutLogPath: stdoutLog,
		StderrLogPath: stderrLog,
		RetryCount:    0,
	}

	log.Printf("[queue] Worker-%d 开始执行%s: %s", workerID, task.Name, task.Cmd)

	// 执行命令
	var output bytes.Buffer
	var err error
	var exitCode int

	for attempt := 0; attempt <= task.RetryCount; attempt++ {
		if attempt > 0 {
			log.Printf("[queue] 任务 %s 第 %d 次重试...", task.Name, attempt)
			time.Sleep(task.RetryDelay)
		}

		result.RetryCount = attempt
		output.Reset()
		exitCode, err = s.runCommand(task, &output, result.StdoutLogPath, result.StderrLogPath)

		if err == nil {
			result.Status = StatusSuccess
			break
		}

		if attempt == task.RetryCount {
			result.Status = StatusFailed
		}
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.ExitCode = exitCode
	result.Error = err

	return result
}

// worker 工作协程
func (s *ShellScheduler) worker(id int) {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.taskQueue:
			result := s.executeTask(id, task)
			s.taskResultQueue <- result
		}
	}
}

// resultProcessor 处理任务结果
func (s *ShellScheduler) resultProcessor() {
	for result := range s.taskResultQueue {
		s.mu.Lock()
		s.taskResults[result.TaskID] = result
		s.completedTasks[result.TaskID] = true
		s.mu.Unlock()
		// 打印结果
		fmt.Printf("[queue] 任务完成: %s (%s)\n", result.TaskName, result.TaskID)
		fmt.Printf("  状态: %s", result.Status)
		fmt.Printf("  耗时: %v", result.Duration)
		fmt.Printf("  开始: %s", result.StartTime.Format(time.DateTime))
		fmt.Printf("  结束: %s", result.EndTime.Format(time.DateTime))
		fmt.Printf("  退出码: %d", result.ExitCode)
		fmt.Printf("  错误信息: %v", result.Error)
		fmt.Printf("  重试次数: %d", result.RetryCount)
	}
}

// AddTasks 批量添加任务
func (s *ShellScheduler) AddTasks(tasks ...*ShellTask) {
	for _, task := range tasks {
		s.AddTask(task)
	}
}

// Start 启动调度器
func (s *ShellScheduler) Start() error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return fmt.Errorf("程序已经在运行")
	}
	s.isRunning = true
	s.mu.Unlock()

	// 启动work
	for i := 0; i < s.maxWorkers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	// 启动结果处理器
	go s.resultProcessor()

	log.Printf("[queue] 调度器启动，最大并发数: %d", s.maxWorkers)
	return nil
}

// Stop 停止调度器
func (s *ShellScheduler) Stop() {
	log.Println("[queue] 停止调度器...")
	s.cancel()
	s.wg.Wait()
	close(s.taskQueue)
	close(s.taskResultQueue)
	s.isRunning = false
	log.Println("[queue] 调度器已停止")
}
