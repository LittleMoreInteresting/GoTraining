package window

import (
	"sync"
	"time"
)

//go-zero/core/collection/rollingwindow.go

var initTime = time.Now().AddDate(-1, -1, -1)

type RollingWindow struct {
	lock     sync.RWMutex
	size     int
	interval time.Duration
	offset   int
	lastTime time.Duration
	buckets  []int64
}

func NewRollingWindow(size int, interval time.Duration) *RollingWindow {
	if size < 1 {
		panic("size must be greater than 0")
	}
	w := &RollingWindow{
		size:     size,
		buckets:  make([]int64, size),
		interval: interval,
		lastTime: time.Since(initTime),
	}
	return w
}

func (rw *RollingWindow) Add(v int64) {
	rw.lock.Lock()
	defer rw.lock.Unlock()
	rw.updateOffset()
	rw.buckets[rw.offset%rw.size] += v
}

//包含当前桶
func (rw *RollingWindow) Reduce() int64 {
	rw.lock.RLock()
	defer rw.lock.RUnlock()

	var diff int
	span := rw.span()
	diff = rw.size - span
	var count int64
	if diff > 0 {
		offset := (rw.offset + span + 1) % rw.size
		for i := 0; i < diff; i++ {
			count += rw.buckets[(offset+i)%rw.size]
		}
	}
	return count
}

// 计算当前距离最后写入数据经过多少个桶
func (rw *RollingWindow) span() int {
	offset := int((time.Since(initTime) - rw.lastTime) / rw.interval)
	if 0 <= offset && offset < rw.size {
		return offset
	}
	return rw.size
}

func (w *RollingWindow) resetBucket(offset int) {
	w.buckets[offset%w.size] = 0
}

func (rw *RollingWindow) updateOffset() {
	span := rw.span()
	if span <= 0 {
		return
	}
	offset := rw.offset

	for i := 0; i < span; i++ {
		rw.resetBucket((offset + i + 1) % rw.size)
	}
	// 更新offset
	rw.offset = (offset + span) % rw.size
	f := time.Since(initTime)
	// 更新操作时间
	rw.lastTime = f - (f-rw.lastTime)%rw.interval
}
