package log

import (
	"log"
	"os"
	"io"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	Init("")
}

func Init(path string){
	if path != ""{
		errFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("打开日志文件失败：", err)
		}
		Error = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
		Info = log.New(io.MultiWriter(os.Stderr, errFile), "Info:", log.Ldate|log.Ltime|log.Lshortfile)
		Warning = log.New(io.MultiWriter(os.Stderr, errFile), "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
	}else {
		Info = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
		Warning = log.New(os.Stdout, "Warning:", log.Ldate|log.Ltime|log.Lshortfile)
		Error = log.New(os.Stdout, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	}
}
