package load_balance

import (
	"testing"
	"go-envoy-poc/analyze/addr"
)

func TestBalancing(t *testing.T) {
	t.Run("正常轮询", func(t *testing.T) {
		round := RoundRobin{}
		address := []addr.SocketAddress{{Host: "localhost", Port: 1}, {Host: "localhost", Port: 2}}
		target := round.Balancing(address)
		if target.Port != 1 {
			t.Error("fail")
		}

		target2 := round.Balancing(address)
		if target2.Port != 2 {
			t.Error("fail")
		}

		target3 := round.Balancing(address)
		if target3.Port != 1 {
			t.Error("fail")
		}
	})
}
