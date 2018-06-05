package analyze

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func Parser(path string) (*StaticResources, error) {
	sr := StaticResources{}
	byte, e := ioutil.ReadFile(path)
	if e != nil {
		log.Fatalf("读取文件错误:%s,文件路径:%s", e, path)
		return nil, e
	}
	err := yaml.Unmarshal(byte, &sr)
	if err != nil {
		return nil, err
	}
	return &sr, nil
}
