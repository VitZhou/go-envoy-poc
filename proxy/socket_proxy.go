package proxy

import (
	"go-envoy-poc/analyze"
	"net"
	"fmt"
	"os"
	"log"
	"io"
	"strconv"
)

func NewSocketProxy(resources *analyze.StaticResources) {
	listener, e := net.Listen("tcp", ":" + strconv.Itoa(resources.Address.Port))
	CheckError(e)
	defer listener.Close()
	log.Println("tcp server started on port 20880 waitting for clients")
	for {
		conn, i := listener.Accept()
		if i != nil {
			continue
		}
		log.Println(conn.RemoteAddr(), "tcp connection success")

		go forward(conn,":20880")
	}
}

func forward(conn net.Conn, remoteAddr string) {
	dial, e := net.Dial("tcp", remoteAddr)
	if dial == nil || e != nil {
		log.Fatalln("remote dial failed:", e)
		return
	}
	go io.Copy(conn, dial)
	go io.Copy(dial, conn)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
