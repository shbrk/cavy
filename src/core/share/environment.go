package share

import (
	"flag"
	"os"
	"strings"
)

//命令行参数
type environment struct {
	BootID   int      // 进程ID
	DevMode  bool     // 开发者模式
	EtcdRoot string   // etcd根节点
	EtcdAddr []string // etcd地址节点 多个用逗号分开
	EtcdUsr  string   // etcd用户名
	EtcdPwd  string   // etcd密码
	Version  string   // 版本号
	LogPath  string   // 日志路径
	AreaID   int      // 区服ID
	Gate     bool     // 是否启动gate
	Game     bool     // 是否启动game
}

var Env environment
var help bool

func init() {
	flag.BoolVar(&help, "h", false, "帮助")
	flag.IntVar(&Env.BootID, "id", 0, "进程id,不能重复，不能省略")
	flag.BoolVar(&Env.DevMode, "d", false, "开发者模式")
	flag.StringVar(&Env.EtcdRoot, "r", "/DEV/", "etcd路径根节点")
	flag.Var(newSliceValue([]string{"127.0.0.1:2379"}, &Env.EtcdAddr), "a", "etcd地址列表，用逗号隔开")
	flag.StringVar(&Env.EtcdUsr, "u", "", "etcd用户名")
	flag.StringVar(&Env.EtcdPwd, "p", "", "etcd用户名")
	flag.StringVar(&Env.Version, "v", "0.0.1", "版本号")
	flag.StringVar(&Env.LogPath, "l", "", "日志路径")
	flag.IntVar(&Env.AreaID, "aid", 10000, "区服ID")
	flag.BoolVar(&Env.Gate, "gate", false, "是否启动gate")
	flag.BoolVar(&Env.Game, "game", false, "是否启动game")
}

func ParseEnvironment() {
	flag.Parse()
	if help || Env.BootID == 0 {
		PrintHelpAndExit()
	}
}

func PrintHelpAndExit() {
	flag.Usage()
	os.Exit(0)
}

//命令行解析数组
type sliceValue []string

func newSliceValue(values []string, p *[]string) *sliceValue {
	*p = values
	return (*sliceValue)(p)
}

func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}

func (s *sliceValue) Get(val string) interface{} {
	return []string(*s)
}

func (s *sliceValue) String() string {
	var str string
	for _, value := range *s {
		str = str + value + ";"
	}
	return str
}
