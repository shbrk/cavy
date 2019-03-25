package gs

import (
	"core/node"
	"encoding/json"
)

var CommonConfig *NodeCommonConfig
var Config *NodeConfig
var AliveConfig *node.GSAliveConfig
var NodeInfo = &node.GSInfoConfig{}

// Gs节点通用配置
type NodeCommonConfig struct {
	ReadBufferSize    int `json:"read_buffer_size"`
	WriteBufferSize   int `json:"write_buffer_size"`
	WriteQueueSize    int `json:"write_queue_size"`
	CompressLevelSize int `json:"compress_level_size"`
	EtcdTimeout       int `json:"etcd_timeout"`
}

// 此Gs节点私有配置
type NodeConfig struct {
	DBAddr string `json:"db_addr"`
}

func ParseCommonConfig(value string) error {
	CommonConfig = &NodeCommonConfig{}
	return json.Unmarshal([]byte(value), CommonConfig)
}

func ParseConfig(value string) error {
	Config = &NodeConfig{}
	return json.Unmarshal([]byte(value), Config)
}

func OnDynamicConfigUpdate(config *node.GSDynamicConfig) {
	//TODO 处理动态信息
}
