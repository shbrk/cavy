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
	etcdClient     *etcd.Client
	publicServer   *net.TCPServer
	internalServer *net.TCPServer
	keepAliveRetry int
	timerManager   *share.TimerManager
}

func NewNodeGate(ctx context.Context, wg *sync.WaitGroup) *NodeGate {
	gNodeGate = &NodeGate{Base: node.Base{Ctx: ctx, Wg: wg}}
	return gNodeGate
}

func (n *NodeGate) Name() string {
	return "Gateway"
}

func (n *NodeGate) Init() error {
	n.timerManager = share.NewTimerManager()
	cli, err := etcd.NewClient(&etcd.Config{
		Endpoints: share.Env.EtcdAddr,
		Timeout:   10,
		Username:  share.Env.EtcdUsr,
		Password:  share.Env.EtcdPwd,
	})
	if err != nil {
		return err
	}
	n.etcdClient = cli

	aliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_ALIVE_PATH, strconv.Itoa(share.Env.BootID))
	exist, _, err := n.etcdClient.SyncGet(aliveKey)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("[ETCD] alive key already exists: " + aliveKey)
	}
	if err := n.readEtcdConfig(); err != nil {
		return err
	}
	var publicListenAddr = "0.0.0.0:" + strconv.Itoa(Config.PublicPort)
	n.publicServer = net.NewTCPServer(publicListenAddr, &net.ConnConfig{
		ReadBufferSize:  CommonConfig.ReadBufferSize,
		WriteBufferSize: CommonConfig.WriteBufferSize,
		WriteQueueSize:  CommonConfig.WriteQueueSize,
	}, NewClientSessionManager())
	if err := n.publicServer.ListenAndServe(); err != nil {
		return err
	}
	n.internalServer = net.NewTCPServer(":0", &net.ConnConfig{
		ReadBufferSize:  CommonConfig.ReadBufferSize,
		WriteBufferSize: CommonConfig.WriteBufferSize,
		WriteQueueSize:  CommonConfig.WriteQueueSize,
	}, NewServerSessionManager())
	if err := n.internalServer.ListenAndServe(); err != nil {
		return err
	}
	if err := n.etcdKeepAlive(); err != nil {
		return err
	}
	n.etcdWatch()
	return nil
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
		}
		n.keepAliveRetry++
		_ = n.etcdKeepAlive()
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

func (n *NodeGate) GetTimer() *share.TimerManager {
	return n.timerManager
}

func (n *NodeGate) GetServerSessionManager() *ServerSessionManager {
	return n.internalServer.GetSessionManager().(*ServerSessionManager)
}

func (n *NodeGate) GetClientSessionManager() *ClientSessionManager {
	return n.publicServer.GetSessionManager().(*ClientSessionManager)
}

func (n *NodeGate) GetServerSession(sessionID uint64) *ServerSession {
	sessionManager := n.GetServerSessionManager()
	if sessionManager == nil {
		return nil
	}
	return sessionManager.GetSession(sessionID)
}

func (n *NodeGate) GetClientSession(sessionID uint64) *ClientSession {
	sessionManager := n.GetClientSessionManager()
	if sessionManager == nil {
		return nil
	}
	return sessionManager.GetSession(sessionID)
}

func (n *NodeGate) stop() {
	gNodeGate = nil
	n.publicServer.Close()
	n.internalServer.Close()
	log.Info(n.Name() + " done!")
}

func (n *NodeGate) Run() {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		n.etcdClient.Run()
		n.publicServer.Run()
		n.internalServer.Run()
		select {
		case _ = <-n.Ctx.Done():
			n.stop()
			n.Wg.Done()
			return
		case now := <-timer.C:
			n.timerManager.Run(now.UnixNano(), 0)
		default:
		}
	}
}

func AddTimer(onTimer func(now int64, param interface{}), delay time.Duration, interval int, param interface{}) *share.Timer {
	return gNodeGate.GetTimer().AddTimer(onTimer, delay, interval, param)
}