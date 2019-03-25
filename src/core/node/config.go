package node

type GateDynamicConfig struct {
	DebugLevelLog bool `json:"debug_level_log"`
}

type GateAliveConfig struct {
	BootID       int    `json:"boot_id"`
	PublicIP     string `json:"public_ip"`
	PublicPort   int    `json:"public_port"`
	InternalIP   string `json:"internal_ip"`
	InternalPort int    `json:"internal_port"`
	RunVersion   string `json:"run_version"`
}

//节点的动态信息
type GateInfoConfig struct {
	ConnectionNum int `json:"connection_num"`
	OnlineNum     int `json:"online_num"`
}

// gs watch的节点信息
type GSDynamicConfig struct {
	DebugLevelLog bool `json:"debug_level_log"`
}

// gs keepalive的节点信息
type GSAliveConfig struct {
	RunVersion string `json:"run_version"`
}

//节点的动态信息
type GSInfoConfig struct {
	ConnectionNum int `json:"connection_num"`
	OnlineNum     int `json:"online_num"`
}
