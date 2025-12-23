package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"student-management-system/config"
	"student-management-system/internal/models"

	"gorm.io/gorm"
)

// Scenario 定义一组压测参数
type Scenario struct {
	Name         string
	Concurrency  int
	BatchSize    int
	TotalBatches int
}

// ScenarioResult 记录压测结果指标
type ScenarioResult struct {
	ScenarioName    string
	Concurrency     int
	BatchSize       int
	Batches         int
	Inserted        int64
	Failures        int64
	DurationSeconds float64
	Throughput      float64
	AvgLatencyMs    float64
	P95LatencyMs    float64
	MaxLatencyMs    float64
	Timestamp       string
}

const resultFile = "grade_stress_results.csv"

func main() {
	log.Println("启动成绩表压测...（确保已导入初始数据）")

	// 初始化数据库连接
	config.InitDB()
	db := config.GetDB()

	// 准备可用的 Enrollment ID
	enrollmentIDs, err := loadEnrollmentIDs(db)
	if err != nil {
		log.Fatalf("获取选课ID失败: %v", err)
	}
	if len(enrollmentIDs) == 0 {
		log.Fatal("未找到选课记录，无法进行成绩压测，请先导入数据")
	}
	log.Printf("已加载 %d 个选课ID，开始压测", len(enrollmentIDs))

	// 多组压测参数，可按需调整
	scenarios := []Scenario{
		{Name: "c8_b50_120", Concurrency: 8, BatchSize: 50, TotalBatches: 120},
		{Name: "c1_b5_100", Concurrency: 1, BatchSize: 5, TotalBatches: 100},     // 低并发基线
		{Name: "c6_b30_200", Concurrency: 6, BatchSize: 30, TotalBatches: 200},   // 中等批量
		{Name: "c12_b80_160", Concurrency: 12, BatchSize: 80, TotalBatches: 160}, // 高并发大批量
		{Name: "c16_b20_400", Concurrency: 16, BatchSize: 20, TotalBatches: 400}, // 高并发小批量
	}

	var results []ScenarioResult
	for _, sc := range scenarios {
		res := runScenario(db, sc, enrollmentIDs)
		results = append(results, res)
		log.Printf("[完成] 场景=%s 插入=%d 失败=%d 耗时=%.2fs 吞吐=%.1f ops/s P95=%.2fms",
			res.ScenarioName, res.Inserted, res.Failures, res.DurationSeconds, res.Throughput, res.P95LatencyMs)
	}

	if err := appendResultsCSV(results); err != nil {
		log.Fatalf("写入CSV失败: %v", err)
	}
	absPath, _ := filepath.Abs(resultFile)
	log.Printf("压测结束，结果已追加至: %s", absPath)
}

func loadEnrollmentIDs(db *gorm.DB) ([]uint, error) {
	var ids []uint
	if err := db.Model(&models.Enrollment{}).Select("id").Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func runScenario(db *gorm.DB, sc Scenario, enrollmentIDs []uint) ScenarioResult {
	start := time.Now()
	tasks := make(chan struct{}, sc.TotalBatches)
	for i := 0; i < sc.TotalBatches; i++ {
		tasks <- struct{}{}
	}
	close(tasks)

	var wg sync.WaitGroup
	var success atomic.Int64
	var failures atomic.Int64
	var mu sync.Mutex
	latencies := make([]time.Duration, 0, sc.TotalBatches)

	for worker := 0; worker < sc.Concurrency; worker++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID*997)))
			for range tasks {
				opStart := time.Now()
				if err := insertBatch(db, r, sc.BatchSize, enrollmentIDs, sc.Name); err != nil {
					failures.Add(1)
				} else {
					success.Add(int64(sc.BatchSize))
				}
				mu.Lock()
				latencies = append(latencies, time.Since(opStart))
				mu.Unlock()
			}
		}(worker)
	}

	wg.Wait()
	duration := time.Since(start)
	avg, p95, max := calcLatencyStats(latencies)

	return ScenarioResult{
		ScenarioName:    sc.Name,
		Concurrency:     sc.Concurrency,
		BatchSize:       sc.BatchSize,
		Batches:         sc.TotalBatches,
		Inserted:        success.Load(),
		Failures:        failures.Load(),
		DurationSeconds: duration.Seconds(),
		Throughput:      float64(success.Load()) / duration.Seconds(),
		AvgLatencyMs:    avg,
		P95LatencyMs:    p95,
		MaxLatencyMs:    max,
		Timestamp:       time.Now().Format(time.RFC3339),
	}
}

func insertBatch(db *gorm.DB, r *rand.Rand, batchSize int, enrollmentIDs []uint, scenarioName string) error {
	if len(enrollmentIDs) == 0 {
		return fmt.Errorf("无可用选课ID")
	}

	grades := make([]models.Grade, batchSize)
	for i := 0; i < batchSize; i++ {
		enrollmentID := enrollmentIDs[r.Intn(len(enrollmentIDs))]
		grades[i] = models.Grade{
			EnrollmentID: enrollmentID,
			ScoreType:    fmt.Sprintf("stress_%s", scenarioName),
			Score:        50 + r.Float64()*50, // 50~100 区间的随机分数
		}
	}

	return db.CreateInBatches(grades, batchSize).Error
}

func calcLatencyStats(latencies []time.Duration) (avgMs, p95Ms, maxMs float64) {
	if len(latencies) == 0 {
		return 0, 0, 0
	}

	var total time.Duration
	max := latencies[0]
	for _, l := range latencies {
		total += l
		if l > max {
			max = l
		}
	}

	sorted := append([]time.Duration(nil), latencies...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	// 防止索引越界，最低取最后一个元素
	idx := int(float64(len(sorted)) * 0.95)
	if idx <= 0 {
		idx = len(sorted) - 1
	} else if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	p95 := sorted[idx]

	return float64(total.Milliseconds()) / float64(len(latencies)),
		float64(p95.Microseconds()) / 1000.0,
		float64(max.Microseconds()) / 1000.0
}

func appendResultsCSV(results []ScenarioResult) error {
	fileExists := false
	if info, err := os.Stat(resultFile); err == nil && info.Size() > 0 {
		fileExists = true
	}

	f, err := os.OpenFile(resultFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	if !fileExists {
		header := []string{
			"scenario", "concurrency", "batch_size", "batches", "inserted", "failures",
			"duration_seconds", "throughput_ops_per_s", "avg_latency_ms", "p95_latency_ms",
			"max_latency_ms", "timestamp",
		}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	for _, r := range results {
		record := []string{
			r.ScenarioName,
			fmt.Sprintf("%d", r.Concurrency),
			fmt.Sprintf("%d", r.BatchSize),
			fmt.Sprintf("%d", r.Batches),
			fmt.Sprintf("%d", r.Inserted),
			fmt.Sprintf("%d", r.Failures),
			fmt.Sprintf("%.4f", r.DurationSeconds),
			fmt.Sprintf("%.2f", r.Throughput),
			fmt.Sprintf("%.2f", r.AvgLatencyMs),
			fmt.Sprintf("%.2f", r.P95LatencyMs),
			fmt.Sprintf("%.2f", r.MaxLatencyMs),
			r.Timestamp,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}
