package gate

import (
	"context"
	"core/etcd"
	"core/log"
	"core/net"
	"core/node"
	"core/share"
	"encoding/json"
	"errors"
	"path"
	"strconv"
	"sync"
	"time"
)

var gNodeGate *NodeGate

type NodeGate struct {
	node.Base
	etcdClient           *etcd.Client
	publicServer         *net.TCPServer        // 面向客户端的Server
	internalServer       *net.TCPServer        // 面向内部连接的Server
	clientSessionManager *ClientSessionManager // 客户端session管理器
	serverSessionManager *ServerSessionManager // 内部session管理器
	keepAliveRetry       int
}

func NewNodeGate(ctx context.Context, wg *sync.WaitGroup) *NodeGate {
	if gNodeGate == nil {
		gNodeGate = &NodeGate{Base: node.Base{Ctx: ctx, Wg: wg}}
	}
	return gNodeGate
}

func (n *NodeGate) Name() string {
	return "Gateway"
}

func (n *NodeGate) Init() {
	cli, err := etcd.NewClient(&etcd.Config{
		Endpoints: share.Env.EtcdAddr,
		Timeout:   10,
		Username:  share.Env.EtcdUsr,
		Password:  share.Env.EtcdPwd,
	})
	share.CheckFatalErr("[GATE]:etcd client create error", err)
	n.etcdClient = cli

	aliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_ALIVE_PATH, strconv.Itoa(share.Env.BootID))
	exist, _, err := n.etcdClient.SyncGet(aliveKey)
	share.CheckFatalErr("[GATE]:etcd get key error："+aliveKey, err)
	if exist {
		log.Fatal("[GATE]:alive key already exists", log.String("key", aliveKey))
	}
	share.CheckFatalErr("[GATE]:read etcd config error", n.readEtcdConfig())

	var publicListenAddr = "0.0.0.0:" + strconv.Itoa(Config.PublicPort)
	n.clientSessionManager = NewClientSessionManager()
	n.publicServer = net.NewTCPServer(publicListenAddr, &net.ConnConfig{
		ReadBufferSize:  CommonConfig.ReadBufferSize,
		WriteBufferSize: CommonConfig.WriteBufferSize,
		WriteQueueSize:  CommonConfig.WriteQueueSize,
	}, n.clientSessionManager)
	share.CheckFatalErr("[GATE]:public server listen error", n.publicServer.ListenAndServe())
	log.Info("[GATE]:server listen for client at " + n.publicServer.Addr())
	n.serverSessionManager = NewServerSessionManager()
	n.internalServer = net.NewTCPServer(":0", &net.ConnConfig{
		ReadBufferSize:  CommonConfig.ReadBufferSize,
		WriteBufferSize: CommonConfig.WriteBufferSize,
		WriteQueueSize:  CommonConfig.WriteQueueSize,
	}, n.serverSessionManager)
	share.CheckFatalErr("[GATE]:internal server listen error", n.internalServer.ListenAndServe())
	log.Info("[GATE]:server listen for internal at " + n.internalServer.Addr())
	share.CheckFatalErr("[GATE]:etcd keep alive error", n.etcdKeepAlive())
	n.etcdWatch()
}

func (n *NodeGate) etcdKeepAlive() error {
	aliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_ALIVE_PATH, strconv.Itoa(share.Env.BootID))
	internalIP, err := share.GetIPByInterface(Config.NetInterfaceName)
	var aliveConfig = &node.GateAliveConfig{
		BootID:       share.Env.BootID,
		PublicIP:     Config.PublicIP,
		PublicPort:   n.publicServer.Port,
		InternalIP:   internalIP,
		InternalPort: n.internalServer.Port,
		RunVersion:   share.Env.Version,
	}
	jsonStr, err := json.Marshal(aliveConfig)
	if err != nil {
		return err
	}
	n.etcdClient.KeepAlive(aliveKey, string(jsonStr), int64(CommonConfig.EtcdTimeout), func(err error, key string) {
		if n.keepAliveRetry < 10 {
			log.Error("[ETCD] keep alive failed,retrying", log.String("key", key),
				log.Int("retry_count", n.keepAliveRetry))
			n.keepAliveRetry++
			_ = n.etcdKeepAlive()
		} else {
			log.Error("[ETCD] keep alive failed, stop retry", log.String("key", key),
				log.Int("retry_count", n.keepAliveRetry))
		}
	})

	return nil
}

func (n *NodeGate) etcdWatch() {
	var dynamicKey = path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_DYNAMIC_PATH, strconv.Itoa(share.Env.BootID))
	n.etcdClient.Watch(dynamicKey, false, func(err error, eventType etcd.EventType, key string, value string) {
		if err != nil {
			log.Error("[ETCD] keep alive failed,retrying", log.String("key", key),
				log.NamedError("err", err))
			return
		}
		if eventType == etcd.EventType_PUT {
			var config = &node.GateDynamicConfig{}
			err := json.Unmarshal([]byte(value), config)
			if err != nil {
				log.Error("[ETCD] keep alive failed,retrying", log.String("key", key),
					log.NamedError("err", err), log.String("value", value))
				return
			}
			OnDynamicConfigUpdate(config)
		}
	})
}

func (n *NodeGate) readEtcdConfig() error {
	commonConfigKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH)
	exist, value, err := n.etcdClient.SyncGet(commonConfigKey)
	if err != nil {
		return err
	}
	if exist == false {
		return errors.New("[ETCD] common config key not exist " + commonConfigKey)
	}
	if err := ParseCommonConfig(value); err != nil {
		return err
	}
	configKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_CONFIG_PATH, strconv.Itoa(share.Env.BootID))
	exist, value, err = n.etcdClient.SyncGet(configKey)
	if err != nil {
		return nil
	}
	if exist == false {
		return errors.New("[ETCD] config key not exist " + configKey)
	}
	if err := ParseConfig(value); err != nil {
		return err
	}
	return nil
}

func (n *NodeGate) GetServerSessionManager() *ServerSessionManager {
	return n.serverSessionManager
}

func (n *NodeGate) GetClientSessionManager() *ClientSessionManager {
	return n.clientSessionManager
}

func (n *NodeGate) GetServerSession(sessionID uint64) *ServerSession {
	return n.serverSessionManager.GetSession(sessionID)
}

func (n *NodeGate) GetClientSession(sessionID uint64) *ClientSession {
	return n.clientSessionManager.GetSession(sessionID)
}

func (n *NodeGate) stop() {
	n.publicServer.Close()
	n.internalServer.Close()
	log.Info(n.Name() + " done!")
}

func (n *NodeGate) Run() {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case _ = <-n.Ctx.Done():
			n.stop()
			n.Wg.Done()
			return
		case event := <-n.etcdClient.ChanOut:
			event.HandleEvent()
		case event := <-n.clientSessionManager.EventChan:
			n.clientSessionManager.HandleEvent(event)
		case event := <-n.serverSessionManager.EventChan:
			n.serverSessionManager.HandleEvent(event)
		case now := <-timer.C:
			n.Tick(now.UnixNano() / int64(time.Millisecond))
		}
	}
}

func (n *NodeGate) Tick(now int64) {

}
