package gs

import (
	"context"
	"core/etcd"
	"core/log"
	"core/net"
	"core/node"
	"core/share"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strconv"
	"sync"
	"time"
)

type NodeGS struct {
	node.Base
	etcdClient     *etcd.Client
	internalClient *net.TCPClient
	keepAliveRetry int
	timerManager   *share.TimerManager
}

func NewNodeGs(ctx context.Context, wg *sync.WaitGroup) *NodeGS {
	return &NodeGS{Base: node.Base{Ctx: ctx, Wg: wg}}
}

func (n *NodeGS) Name() string {
	return "GameServer"
}

func (n *NodeGS) Init() error {
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
	aliveKey := path.Join(share.Env.EtcdRoot, strconv.Itoa(share.Env.AreaID), share.ETCD_GS_PATH, share.ETCD_ALIVE_PATH,
		strconv.Itoa(share.Env.BootID))
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
	n.internalClient = net.NewTCPClient(5*time.Second, &net.ConnConfig{
		ReadBufferSize:  CommonConfig.ReadBufferSize,
		WriteBufferSize: CommonConfig.WriteBufferSize,
		WriteQueueSize:  CommonConfig.WriteQueueSize,
	}, net.NewGateSessionManager())

	gateAliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_ALIVE_PATH)
	_, values, err := n.etcdClient.SyncGetWithPrefix(gateAliveKey + "/")
	if err != nil {
		return err
	}
	if len(values) == 0 {
		return errors.New("there has no gate for gs to connect")
	}
	for index := range values {
		gateConfig := &node.GateAliveConfig{}
		err := json.Unmarshal([]byte(values[index]), gateConfig)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("%s:%d", gateConfig.InternalIP, gateConfig.InternalPort)
		_, err = n.internalClient.SyncConnect(addr, 5*time.Second, gateConfig.BootID)
		if err != nil {
			return err
		}
	}
	if err := n.etcdKeepAlive(); err != nil {
		return err
	}
	n.etcdWatch()
	return nil
}

func (n *NodeGS) readEtcdConfig() error {
	commonConfigKey := path.Join(share.Env.EtcdRoot, share.ETCD_GS_PATH)
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
	configKey := path.Join(share.Env.EtcdRoot, share.ETCD_GS_PATH, strconv.Itoa(share.Env.AreaID),
		share.ETCD_CONFIG_PATH, strconv.Itoa(share.Env.BootID))
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

func (n *NodeGS) etcdKeepAlive() error {
	aliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GS_PATH, strconv.Itoa(share.Env.AreaID),
		share.ETCD_ALIVE_PATH, strconv.Itoa(share.Env.BootID))
	var aliveConfig = &node.GSAliveConfig{
		RunVersion: share.Env.Version,
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
func (n *NodeGS) etcdWatch() {
	var dynamicKey = path.Join(share.Env.EtcdRoot, share.ETCD_GS_PATH, strconv.Itoa(share.Env.AreaID),
		share.ETCD_DYNAMIC_PATH, strconv.Itoa(share.Env.BootID))
	n.etcdClient.Watch(dynamicKey, false, func(err error, eventType etcd.EventType, key string, value string) {
		if err != nil {
			log.Error("[ETCD] watch failed", log.String("key", key),
				log.NamedError("err", err))
			return
		}
		if eventType == etcd.EventType_PUT {
			var config = &node.GSDynamicConfig{}
			err := json.Unmarshal([]byte(value), config)
			if err != nil {
				log.Error("[ETCD] parse json error", log.String("key", key),
					log.NamedError("err", err), log.String("value", value))
				return
			}
			OnDynamicConfigUpdate(config)
		}
	})
	gateAliveKey := path.Join(share.Env.EtcdRoot, share.ETCD_GATE_PATH, share.ETCD_ALIVE_PATH)
	n.etcdClient.Watch(gateAliveKey+"/", true, func(err error, eventType etcd.EventType, key string, value string) {
		if err != nil {
			log.Error("[ETCD] watch gate error", log.String("key", key),
				log.NamedError("err", err))
			return
		}
		sessionManager := n.internalClient.GetSessionManager().(*net.GateSessionManager)
		var bootID, _ = strconv.Atoi(key)
		if eventType == etcd.EventType_PUT {
			var config = &node.GateAliveConfig{}
			err := json.Unmarshal([]byte(value), config)
			if err != nil {
				log.Error("[ETCD] parse json error", log.String("key", key),
					log.NamedError("err", err), log.String("value", value))
				return
			}
			sessionManager.HandleEtcdEventPut(config.BootID, fmt.Sprintf("%s:%d", config.InternalIP, config.InternalPort))
		} else if eventType == etcd.EventType_DELETE {
			sessionManager.HandleEtcdEventDelete(bootID)
		}
	})

}

func (n *NodeGS) stop() {
	n.internalClient.Close()
	log.Info(n.Name() + " done!")
}

func (n *NodeGS) Run() {
	timer := time.NewTicker(100 * time.Millisecond)
	for {
		n.etcdClient.Run()
		n.internalClient.Run()
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