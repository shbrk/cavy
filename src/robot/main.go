package main

import (
	"core/log"
	"flag"
	"fmt"
	"robot/client"
	"sync"
)

var count = 1
var port = 8848
var ip = "172.25.35.250"

func init() {
	flag.IntVar(&count, "c", 1, "客户端并发数量")
	flag.IntVar(&port, "p", 8848, "客户端连接端口")
	flag.StringVar(&ip, "h", "172.25.35.250", "客户端连接ip")
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(count)
	log.InitDefaultLogger(&log.Config{
		DevMode:    true, // 开发者模式
		CallerSkip: 3,
		Filepath:   "../robot_log/", // 日志文件夹路径
		MaxSize:    1,               // 单个日志文件最大Size
		MaxBackups: 5,               // 最多同时保留的日志数量,0是全部
		MaxAge:     2,               // 日志保留的持续天数，0是永久
		Compress:   false,           // 日志文件是否用zip压缩
	})
	log.Info("wg", log.Int("count", count))
	list := make([]*client.Client, 0, count)
	for i := 0; i < count; i++ {
		cli := &client.Client{ID: i, Wg: &wg}
		cli.Connect(fmt.Sprintf("%s:%d", ip, port))
		list = append(list, cli)
	}
	for _, cli := range list {
		go cli.Run()
	}
	wg.Wait()
}
