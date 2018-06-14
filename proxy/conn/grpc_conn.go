package conn

import (
	"log"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"sync"
)

var cache = sync.Map{}

func GetConn(context context.Context) (conn *grpc.ClientConn, err error) {
	value, ok := cache.Load("a")
	if ok {
		return value.(*grpc.ClientConn), nil
	}else {
		clientConn, e := grpc.DialContext(context, "localhost:50051", grpc.WithInsecure())
		if e != nil {
			log.Println("创建客户端连接失败", e)
			return nil, e
		}
		cache.Store("a", clientConn)
		return clientConn, nil
	}
}
