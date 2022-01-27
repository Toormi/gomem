package _map

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	order2 "gomem/order"
	"gomem/protopb/order"
	"sync"
	"time"
)

//var num = 50000
//
var entries = 5000000

// Memory orders map[cid][market_uuid][id]order
type Memory struct {
	openingMtx    sync.RWMutex
	OpeningOrders map[uint64]*MarketOrderMemory
	closeMtx      sync.RWMutex
	CloseOrders   map[uint64]*MarketOrderMemory
	mtx           sync.RWMutex
	orderIDMtx    sync.RWMutex
	OrderIDMap    map[uint64]*order.Order
	clientIDMtx   sync.RWMutex
	ClientIDMap   map[string]*order.Order
}

type OrderMap struct {
	OrderID    uint64
	ClientID   string
	CustomerID uint64
	MarketUUID string
}

type MarketOrderMemory struct {
	mtx          sync.RWMutex
	MarketOrders map[string]*OrderMemory
}

func (m *MarketOrderMemory) GetOrStoreOrderMemory(marketUUID string) *OrderMemory {
	mo := &OrderMemory{}
	m.mtx.Lock()
	acm, ok := m.MarketOrders[marketUUID]
	if !ok {
		acm = mo
		m.MarketOrders = make(map[string]*OrderMemory)
		m.MarketOrders[marketUUID] = acm
	}
	m.mtx.Unlock()

	return acm
}

func (m *MarketOrderMemory) GetOrderMemory(marketUUID string) *OrderMemory {
	m.mtx.RLock()
	acm := m.MarketOrders[marketUUID]
	m.mtx.RUnlock()

	return acm
}

type OrderMemory struct {
	mtx    sync.RWMutex
	Orders map[uint64]*order.Order
}

func (m *OrderMemory) StoreOrder(ord *order.Order) {
	m.mtx.Lock()
	_, ok := m.Orders[ord.Id]
	if !ok {
		m.Orders = make(map[uint64]*order.Order)
	}
	m.Orders[ord.Id] = ord
	m.mtx.Unlock()
}

func (m *OrderMemory) GetOrder(id uint64) *order.Order {
	m.mtx.RLock()
	o := m.Orders[id]
	m.mtx.RUnlock()
	return o
}

func (m *OrderMemory) ListOrder() []*order.Order {
	var orders []*order.Order
	m.mtx.RLock()
	for _, o := range m.Orders {
		orders = append(orders, o)
	}
	m.mtx.RUnlock()
	return orders
}

func (m *OrderMemory) DeleteOrder(id uint64) {
	m.mtx.Lock()
	delete(m.Orders, id)
	m.mtx.Unlock()
}

func NewMemory() *Memory {
	return &Memory{
		OpeningOrders: make(map[uint64]*MarketOrderMemory),
		CloseOrders:   make(map[uint64]*MarketOrderMemory),
		OrderIDMap:    make(map[uint64]*order.Order),
		ClientIDMap:   make(map[string]*order.Order),
	}
}

func (m *Memory) GetOrStoreOpeningMarketMemory(cid uint64) *MarketOrderMemory {
	cm := &MarketOrderMemory{}
	m.openingMtx.Lock()
	ac, ok := m.OpeningOrders[cid]
	if !ok {
		ac = cm
		m.OpeningOrders[cid] = ac
	}
	m.openingMtx.Unlock()
	return ac
}

func (m *Memory) GetOpeningMarketMemory(cid uint64) *MarketOrderMemory {
	m.openingMtx.RLock()
	ac := m.OpeningOrders[cid]
	m.openingMtx.RUnlock()
	return ac
}

func (m *Memory) GetOrStoreCloseMarketMemory(cid uint64) *MarketOrderMemory {
	cm := &MarketOrderMemory{}
	m.closeMtx.Lock()
	ac, ok := m.CloseOrders[cid]
	if !ok {
		ac = cm
		m.CloseOrders[cid] = ac
	}
	m.closeMtx.Unlock()
	return ac
}

func (m *Memory) GetCloseMarketMemory(cid uint64) *MarketOrderMemory {
	m.closeMtx.Lock()
	ac := m.CloseOrders[cid]
	m.closeMtx.Unlock()
	return ac
}

func (m *Memory) SetOpeningOrder(order *order.Order) {
	ac := m.GetOrStoreOpeningMarketMemory(order.CustomerId)
	om := ac.GetOrStoreOrderMemory(order.MarketUuid)
	om.StoreOrder(order)
}

func (m *Memory) SetCloseOrder(order *order.Order) {
	m.setOrderID(order)
	ac := m.GetOrStoreCloseMarketMemory(order.CustomerId)
	om := ac.GetOrStoreOrderMemory(order.MarketUuid)
	om.StoreOrder(order)
}

func (m *Memory) setOrderID(order *order.Order) {
	m.orderIDMtx.Lock()
	m.OrderIDMap[order.Id] = order
	m.orderIDMtx.Unlock()
}

func (m *Memory) GetOpeningOrderByID(ID uint64) (*order.Order, error) {
	orderMap := m.GetOrderIDMap(ID)
	if orderMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return m.GetOpeningByOrderMap(orderMap)
}

func (m *Memory) GetOpeningOrdersByCidAndMarketUUID(cid uint64, marketUUID string) ([]*order.Order, error) {
	mo := m.GetOpeningMarketMemory(cid)
	if mo == nil {
		return []*order.Order{}, nil
	}

	cm := mo.GetOrderMemory(marketUUID)
	if cm == nil {
		return []*order.Order{}, nil
	}
	return cm.ListOrder(), nil
}

func (m *Memory) GetOpeningOrdersByCidAndMarketUUIDs(cid uint64, marketUUIDs []string) ([]*order.Order, error) {
	mo := m.GetOpeningMarketMemory(cid)
	if mo == nil {
		return []*order.Order{}, nil
	}

	return GetOrders(mo, marketUUIDs)
}

func (m *Memory) GetCloseOrdersByCidAndMarketUUID(cid uint64, marketUUID string) ([]*order.Order, error) {
	mo := m.GetCloseMarketMemory(cid)
	if mo == nil {
		return []*order.Order{}, nil
	}

	cm := mo.GetOrderMemory(marketUUID)
	if cm == nil {
		return []*order.Order{}, nil
	}
	return cm.ListOrder(), nil
}

func (m *Memory) GetCloseOrdersByCidAndMarketUUIDs(cid uint64, marketUUIDs []string) ([]*order.Order, error) {
	mo := m.GetCloseMarketMemory(cid)
	if mo == nil {
		return []*order.Order{}, nil
	}

	return GetOrders(mo, marketUUIDs)
}

func GetOrders(mo *MarketOrderMemory, marketUUIDs []string) ([]*order.Order, error) {
	var orders []*order.Order
	mMap := toMap(marketUUIDs)
	mLen := len(marketUUIDs)
	mo.mtx.RLock()
	for marketUUID, om := range mo.MarketOrders {
		if _, ok := mMap[marketUUID]; ok || mLen <= 0 {
			orders = append(orders, om.ListOrder()...)
		}
	}
	mo.mtx.RUnlock()
	return orders, nil
}

func (m *Memory) GetOpeningOrderByClientID(ClientID string) (*order.Order, error) {
	orderMap := m.GetClientIDMap(ClientID)
	if orderMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return m.GetOpeningByOrderMap(orderMap)
}

func (m *Memory) GetCloseOrderByID(ID uint64) (*order.Order, error) {
	orderMap := m.GetOrderIDMap(ID)
	if orderMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return m.GetCloseByOrderMap(orderMap)
}

func (m *Memory) GetCloseOrderByClientID(ClientID string) (*order.Order, error) {
	orderMap := m.GetClientIDMap(ClientID)
	if orderMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return m.GetCloseByOrderMap(orderMap)
}

func (m *Memory) GetOpeningByOrderMap(o *order.Order) (*order.Order, error) {
	mo := m.GetOpeningMarketMemory(o.CustomerId)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cm := mo.GetOrderMemory(o.MarketUuid)
	if cm == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord := cm.GetOrder(o.Id)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return ord, nil
}

func (m *Memory) GetCloseByOrderMap(o *order.Order) (*order.Order, error) {
	mo := m.GetCloseMarketMemory(o.CustomerId)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cm := mo.GetOrderMemory(o.MarketUuid)
	if cm == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord := cm.GetOrder(o.Id)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return ord, nil
}

func (m *Memory) DeleteOpeningOrderByID(id uint64) error {
	ordMap := m.GetOrderIDMap(id)

	mo := m.GetOpeningMarketMemory(ordMap.CustomerId)
	if mo == nil {
		return gorm.ErrRecordNotFound
	}

	cu := mo.GetOrderMemory(ordMap.MarketUuid)
	if cu == nil {
		return gorm.ErrRecordNotFound
	}
	cu.DeleteOrder(ordMap.Id)
	return nil
}

func (m *Memory) UpdateOpeningOrderByID(id uint64, price string) (*order.Order, error) {
	ord := m.GetOrderIDMap(id)
	ord.Price = price
	return ord, nil
}

func (m *Memory) GetOrderIDMap(ID uint64) *order.Order {
	m.orderIDMtx.RLock()
	v := m.OrderIDMap[ID]
	m.orderIDMtx.RUnlock()
	return v
}

func (m *Memory) GetClientIDMap(clientID string) *order.Order {
	m.clientIDMtx.RLock()
	v := m.ClientIDMap[clientID]
	m.clientIDMtx.RUnlock()
	return v
}

func toMap(uuids []string) map[string]struct{} {
	mMap := make(map[string]struct{})
	for _, uuid := range uuids {
		mMap[uuid] = struct{}{}
	}

	return mMap
}

type Snapshot struct {
	Orders   []*order.Order
}

func (m *Memory) Snapshot() *Snapshot {
	snap := &Snapshot{
		Orders: make([]*order.Order, 0, entries),
	}

	now := time.Now()
	for _, o := range m.OrderIDMap {
		snap.Orders = append(snap.Orders, o)
	}
	now2 := time.Now()
	fmt.Println("close order copy ", now2.Sub(now))

	return snap
}

func TestMap() {
	db := NewMemory()
	//var clientIDs []string
	//clientID := uuid.New().String()
	Insert(db)
	//time.Sleep(time.Second * 10)
	fmt.Println("start to deep copy ")
	now3 := time.Now()
	snap := db.Snapshot()
	now4 := time.Now()
	fmt.Println("deep copy snapshot ", now4.Sub(now3))
	fmt.Println("snapshot order len", len(snap.Orders))

	//now5 := time.Now().UnixNano() / 1e6
	//data, err := json.Marshal(snap)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//snapshot.BackUp(1, data)
	//now6 := time.Now().UnixNano() / 1e6
	//fmt.Println("save to file:", now6-now5)
}

func Insert(db *Memory) {
	fmt.Println("start insert")
	now := time.Now()
	for i := 1; i <= entries; i++ {
		if i % 100000 == 0 {
			fmt.Println("insert processing: ", i)
		}
		clientID := uuid.New().String()
		//clientIDs = append(clientIDs, clientID)
		ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
		db.SetCloseOrder(ord)
	}
	now2 := time.Now()
	fmt.Println("insert into eslap ", now2.Sub(now))
}
