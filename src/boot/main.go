package main

import (
	"context"
	"core/log"
	"core/node"
	"core/share"
	"fmt"
	"game"
	"gate"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	//格式化命令行参数
	share.ParseEnvironment()
	//初始化log
	log.InitDefaultLogger(&log.Config{
		DevMode:    share.Env.DevMode,
		Filepath:   share.Env.LogPath,
		MaxSize:    20,
		MaxBackups: 20,
		MaxAge:     5,
		CallerSkip: 3,
	})
	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	var nodes = LoadNodes(ctx, &wg)
	if len(nodes) == 0 {
		share.PrintHelpAndExit()
	}
	wg.Add(len(nodes))
	//加载所有节点
	for i := range nodes {
		nodes[i].Init()
		log.Info("[BOOT]:"+ nodes[i].Name() +" start")
	}
	// 每个节点运行在独立的线程上
	for i := range nodes {
		go nodes[i].Run()
	}
	//收到信号 主动退出
	signalChan := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGPIPE)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGFPE)
	sig := <-signalChan
	fmt.Println("[BOOT]:received a signal:" + sig.String())
	cancel()  // 关闭所有node
	wg.Wait() //等待所有node关闭成功后，主线程做清理工作
	fmt.Println("[BOOT]:boot closed")
}

func LoadNodes(ctx context.Context, wg *sync.WaitGroup) []node.Node {
	var nodes = make([]node.Node, 0)
	if share.Env.Gate {
		nodes = append(nodes, gate.NewNodeGate(ctx, wg))
	}
	if share.Env.Game {
		nodes = append(nodes, game.NewNodeGame(ctx, wg))
	}
	return nodes
}
