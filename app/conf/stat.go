package conf

import (
	"fmt"
	"sync"
)

const maxRecords = 1000

type stat struct {
	mu      sync.RWMutex
	records []bool // true表示成功，false表示失败
	index   int    // 当前写入位置
	count   int    // 已记录的总数（最多maxRecords）
}

var (
	data sync.Map // map[string]*stat
)

func getStat(net string) *stat {
	val, _ := data.LoadOrStore(net, &stat{
		records: make([]bool, maxRecords),
	})
	return val.(*stat)
}

func RecordSuccess(net string) {
	s := getStat(net)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[s.index] = true
	s.index = (s.index + 1) % maxRecords
	if s.count < maxRecords {
		s.count++
	}
}

func RecordFailure(net string) {
	s := getStat(net)
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[s.index] = false
	s.index = (s.index + 1) % maxRecords
	if s.count < maxRecords {
		s.count++
	}
}

func GetSuccessRate(net string) string {
	s := getStat(net)
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.count == 0 {
		return "100.00%"
	}

	successCount := 0
	for i := 0; i < s.count; i++ {
		if s.records[i] {
			successCount++
		}
	}

	return fmt.Sprintf("%.2f%%", float64(successCount)/float64(s.count)*100)
}
