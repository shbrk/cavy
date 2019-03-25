package gate

import (
	"core/node"
	"encoding/json"
)

var CommonConfig *NodeCommonConfig
var Config *NodeConfig
var AliveConfig *node.GateAliveConfig
var NodeInfo = &node.GateInfoConfig{}

// Gate节点通用配置
type NodeCommonConfig struct {
	ReadBufferSize    int `json:"read_buffer_size"`
	WriteBufferSize   int `json:"write_buffer_size"`
	WriteQueueSize    int `json:"write_queue_size"`
	CompressLevelSize int `json:"compress_level_size"`
	EtcdTimeout       int `json:"etcd_timeout"`
}

// 此Gate节点私有配置
type NodeConfig struct {
	PublicIP         string `json:"public_ip"`
	PublicPort       int    `json:"public_port"`
	NetInterfaceName string `json:"net_interface_name"` // 根据网卡获取内部地址，内部端口使用随机端口
}

func ParseCommonConfig(value string) error {
	CommonConfig = &NodeCommonConfig{}
	return json.Unmarshal([]byte(value), CommonConfig)
}

func ParseConfig(value string) error {
	Config = &NodeConfig{}
	return json.Unmarshal([]byte(value), Config)
}

func OnDynamicConfigUpdate(config *node.GateDynamicConfig) {
	//TODO 处理动态信息
}
