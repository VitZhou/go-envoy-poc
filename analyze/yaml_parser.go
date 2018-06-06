package analyze

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"go-envoy-poc/log"
)

func Parser(path string) (*StaticResources, error) {
	sr := StaticResources{}
	bytes, e := ioutil.ReadFile(path)
	if e != nil {
		log.Error.Fatalf("读取文件错误:%s,文件路径:%s", e, path)
		return nil, e
	}
	err := yaml.Unmarshal(bytes, &sr)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}
