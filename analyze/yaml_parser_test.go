package analyze

import "testing"

func TestParser(t *testing.T) {
	t.Run("文件路径正确", func(t *testing.T) {
		_, e := Parser("./envoy_test.yaml")
		if e != nil{
			t.Error("fail")
		}
	})

	t.Run("文件路径不正确", func(t *testing.T) {
		_, e := Parser("./not_exist.yaml")
		if e == nil{
			t.Error("fail")
		}
	})
}