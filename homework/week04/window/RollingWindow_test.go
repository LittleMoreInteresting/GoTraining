package window

import (
	"fmt"
	"testing"
	"time"
)

const duration = time.Millisecond * 50

func TestNewRollingWindow(t *testing.T) {
	const size = 3
	r := NewRollingWindow(size, duration)
	r.Add(1)
	fmt.Println(r.Reduce())
	r.Add(1)
	time.Sleep(1 * duration)
	r.Add(1)
	time.Sleep(1 * duration)
	//time.Sleep(1*duration)
	//time.Sleep(1*duration)
	fmt.Println(r.Reduce())
	r.Add(1)
	r.Add(1)
	time.Sleep(2 * duration)
	fmt.Println(r.Reduce())
}
