package monitor

import (
	"net/http"
	"strings"
	"sync/atomic"
	"sync"
	"time"
	"go-envoy-poc/log"
)

var (
	resp_1xx      int64 = 0
	resp_2xx      int64 = 0
	resp_3xx      int64 = 0
	resp_4xx      int64 = 0
	resp_5xx      int64 = 0
	resp_total    int64 = 0
	C                   = make(chan *http.Response, 20)
	T                   = make(chan int64, 20)
	resp_success  int64 = 0
	resp_fail     int64 = 0
	comsumingTime int64 = 0
	once          sync.Once
	instance      *Monitor
)

type Monitor struct {
}

func Instance() *Monitor {
	once.Do(func() {
		instance = &Monitor{}
	})
	return instance
}

//noinspection GoRedundantParens
func init() {
	go func() {
		for {
			resp := <-C
			if (strings.HasPrefix(resp.Status, "1")) {
				atomic.AddInt64(&resp_1xx, 1)
			} else if (strings.HasPrefix(resp.Status, "2")) {
				atomic.AddInt64(&resp_2xx, 1)
			} else if (strings.HasPrefix(resp.Status, "3")) {
				atomic.AddInt64(&resp_3xx, 1)
			} else if (strings.HasPrefix(resp.Status, "4")) {
				atomic.AddInt64(&resp_4xx, 1)
			} else if (strings.HasPrefix(resp.Status, "5")) {
				atomic.AddInt64(&resp_5xx, 1)
			}

			if ("200" == resp.Status) {
				atomic.AddInt64(&resp_success, 1)
			} else {
				atomic.AddInt64(&resp_fail, 1)
			}
			atomic.AddInt64(&resp_total, 1)
		}
	}()

	go func() {
		for {
			t := <-T
			atomic.AddInt64(&comsumingTime, t)
		}
	}()

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			if atomic.LoadInt64(&resp_total) > 0 {
				t :=atomic.LoadInt64(&comsumingTime) / atomic.LoadInt64(&resp_total) / 1000000
				log.Info.Printf("总请求数%d,平均耗时%d毫秒,\n", atomic.LoadInt64(&resp_total), t)
			}
		}
	}()
}

func (monitor *Monitor) NotifyResp(resp *http.Response) {
	C <- resp
}

func (monitor *Monitor) ConsumeTime(time int64) {
	T <- time
}
