package entity

import (
	"proto/entity"
)

//====================
// 实体系统，统一管理服务器的所有落地数据
// 自动完成客户端和数据库同步
//===================

type UUID = uint64
type Entities = map[UUID]Entity
type Container = map[entity.TYPE]Entities

//通知接口，实体改动的时候通知客户端和DB
type Notifier interface {
	SendToClientList(owner uint64, c Container)
	SendToClientAdd(owner uint64, entityType entity.TYPE, list []Entity)
	SendToClientUpdate(owner uint64, entityType entity.TYPE, list []Entity)
	SendToClientDelete(owner uint64, entityType entity.TYPE, list []UUID)
	SendToDBAdd(owner uint64, entityType entity.TYPE, list []Entity)
	SendToDBUpdate(owner uint64, entityType entity.TYPE, list []Entity)
	SendToDBDelete(owner uint64, entityType entity.TYPE, list []UUID)
}

// 实体接口
type Entity interface {
	GetUUID() uint64 // 每个实体有唯一UUID
}

func NewManager(owner uint64, notifier Notifier) *Manager {
	return &Manager{
		owner:     owner,
		container: make(Container),
		notifier:  notifier,
	}
}

type Manager struct {
	owner     uint64
	container Container
	notifier  Notifier
}

//数据初始化（从DB加载）
func (m *Manager) DataInit(c Container) {
	m.container = c
}

//数据后续追加（可能玩家登陆的时候不会一次性把所有数据加载出来）
func (m *Manager) DataAppend(c Container) {
	for t, list := range c {
		l, ok := m.container[t]
		if !ok {
			m.container[t] = list
			continue
		}
		for uuid, ett := range list {
			l[uuid] = ett
		}
	}
}

//获取一个实体对象
func (m *Manager) Get(entityType entity.TYPE, uuid UUID) Entity {
	l, ok := m.container[entityType]
	if !ok {
		return nil
	}
	return l[uuid]
}

//获取该类型的所有实体
func (m *Manager) GetAll(entityType entity.TYPE) map[UUID]Entity {
	l, ok := m.container[entityType]
	if !ok {
		return nil
	}
	return l
}

//获取该类型的任意一个实体
func (m *Manager) GetAny(entityType entity.TYPE) Entity {
	l, ok := m.container[entityType]
	if !ok {
		return nil
	}
	for _, ett := range l {
		return ett
	}
	return nil
}

//添加
func (m *Manager) Add(entityType entity.TYPE, entity Entity, notifyClient bool) {
	var list = []Entity{entity}
	m.cacheAdd(entityType, list)
	m.sendToDBAdd(entityType, list)
	if notifyClient {
		m.sendToClientAdd(entityType, list)
	}
}
func (m *Manager) AddBatch(entityType entity.TYPE, list []Entity, notifyClient bool) {
	m.cacheAdd(entityType, list)
	m.sendToDBAdd(entityType, list)
	if notifyClient {
		m.sendToClientAdd(entityType, list)
	}
}

func (m *Manager) Update(entityType entity.TYPE, entity Entity, notifyClient bool) {
	var list = []Entity{entity}
	m.cacheUpdate(entityType, list)
	m.sendToDBUpdate(entityType, list)
	if notifyClient {
		m.sendToClientUpdate(entityType, list)
	}
}

func (m *Manager) UpdateBatch(entityType entity.TYPE, list []Entity, notifyClient bool) {
	m.cacheUpdate(entityType, list)
	m.sendToDBUpdate(entityType, list)
	if notifyClient {
		m.sendToClientUpdate(entityType, list)
	}
}

func (m *Manager) Delete(entityType entity.TYPE, entity Entity, notifyClient bool) {
	var uuidList = []UUID{entity.GetUUID()}
	m.cacheDelete(entityType, uuidList)
	m.sendToDBDelete(entityType, uuidList)
	if notifyClient {
		m.sendToClientDelete(entityType, uuidList)
	}
}

func (m *Manager) DeleteBatch(entityType entity.TYPE, list []Entity, notifyClient bool) {
	var uuidList = make([]UUID, 0, len(list))
	for _, ett := range list {
		uuidList = append(uuidList, ett.GetUUID())
	}
	m.cacheDelete(entityType, uuidList)
	m.sendToDBDelete(entityType, uuidList)
	if notifyClient {
		m.sendToClientDelete(entityType, uuidList)
	}
}

func (m *Manager) DeleteByKey(entityType entity.TYPE, uuid UUID, notifyClient bool) {
	var uuidList = []UUID{uuid}
	m.cacheDelete(entityType, uuidList)
	m.sendToDBDelete(entityType, uuidList)
	if notifyClient {
		m.sendToClientDelete(entityType, uuidList)
	}
}

func (m *Manager) DeleteBatchByKey(entityType entity.TYPE, list []UUID, notifyClient bool) {
	m.cacheDelete(entityType, list)
	m.sendToDBDelete(entityType, list)
	if notifyClient {
		m.sendToClientDelete(entityType, list)
	}
}

/// ====================================发送给客户端实体列表===============================================
func (m *Manager) SendToClientAll() {
	m.notifier.SendToClientList(m.owner, m.container)
}
func (m *Manager) SendToClient(entityTypes []entity.TYPE) {
	var c = make(map[entity.TYPE]Entities)
	for _, entityType := range entityTypes {
		l, ok := m.container[entityType]
		if ok {
			c[entityType] = l
		}
	}
	m.notifier.SendToClientList(m.owner, c)
}

/// ================================== 同步给客户端============================================
func (m *Manager) sendToClientAdd(entityType entity.TYPE, list []Entity) {
	m.notifier.SendToClientAdd(m.owner, entityType, list)
}

func (m *Manager) sendToClientUpdate(entityType entity.TYPE, list []Entity) {
	m.notifier.SendToClientUpdate(m.owner, entityType, list)
}

func (m *Manager) sendToClientDelete(entityType entity.TYPE, list []UUID) {
	m.notifier.SendToClientDelete(m.owner, entityType, list)
}

/// ====================================同步到数据库===============================================
func (m *Manager) sendToDBAdd(entityType entity.TYPE, list []Entity) {
	m.notifier.SendToDBAdd(m.owner, entityType, list)
}

func (m *Manager) sendToDBUpdate(entityType entity.TYPE, list []Entity) {
	m.notifier.SendToDBUpdate(m.owner, entityType, list)
}

func (m *Manager) sendToDBDelete(entityType entity.TYPE, list []UUID) {
	m.notifier.SendToDBDelete(m.owner, entityType, list)
}

///==================================修改本地cache=============================================

func (m *Manager) cacheAdd(entityType entity.TYPE, list []Entity) {
	l, ok := m.container[entityType]
	if !ok {
		l = make(Entities)
		m.container[entityType] = l
	}
	for _, ett := range list {
		l[ett.GetUUID()] = ett
	}
}

func (m *Manager) cacheUpdate(entityType entity.TYPE, list []Entity) {
	// 因为操作的实体是指针，所以cache一直是新的 所以不用更新
}

func (m *Manager) cacheDelete(entityType entity.TYPE, list []UUID) {
	l, ok := m.container[entityType]
	if !ok {
		return
	}
	for _, uuid := range list {
		delete(l, uuid)
	}
}
