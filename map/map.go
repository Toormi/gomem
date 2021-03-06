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
var num = 5000000

// Memory orders map[cid][market_uuid][id]order
type Memory struct {
	openingMtx    sync.RWMutex
	OpeningOrders map[uint64]*MarketOrderMemory
	closeMtx      sync.RWMutex
	CloseOrders   map[uint64]*MarketOrderMemory
	mtx           sync.RWMutex
	orderIDMtx    sync.RWMutex
	OrderIDMap    map[uint64]*OrderMap
	clientIDMtx   sync.RWMutex
	ClientIDMap   map[string]*OrderMap
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
		OrderIDMap:    make(map[uint64]*OrderMap),
		ClientIDMap:   make(map[string]*OrderMap),
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
	ac := m.GetOrStoreCloseMarketMemory(order.CustomerId)
	om := ac.GetOrStoreOrderMemory(order.MarketUuid)
	om.StoreOrder(order)
}

func (m *Memory) SetOrderMap(order *order.Order, clientID string) {
	orderMap := &OrderMap{
		OrderID:    order.Id,
		ClientID:   clientID,
		CustomerID: order.CustomerId,
		MarketUUID: order.MarketUuid,
	}
	m.orderIDMtx.Lock()
	m.OrderIDMap[order.Id] = orderMap
	m.orderIDMtx.Unlock()

	if clientID != "" {
		m.clientIDMtx.Lock()
		m.ClientIDMap[clientID] = orderMap
		m.clientIDMtx.Unlock()
	}
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

func (m *Memory) GetOpeningByOrderMap(orderMap *OrderMap) (*order.Order, error) {
	mo := m.GetOpeningMarketMemory(orderMap.CustomerID)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cm := mo.GetOrderMemory(orderMap.MarketUUID)
	if cm == nil {
		return nil, gorm.ErrRecordNotFound
	}

	o := cm.GetOrder(orderMap.OrderID)
	if o == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return o, nil
}

func (m *Memory) GetCloseByOrderMap(orderMap *OrderMap) (*order.Order, error) {
	mo := m.GetCloseMarketMemory(orderMap.CustomerID)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cm := mo.GetOrderMemory(orderMap.MarketUUID)
	if cm == nil {
		return nil, gorm.ErrRecordNotFound
	}

	o := cm.GetOrder(orderMap.OrderID)
	if o == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return o, nil
}

func (m *Memory) DeleteOpeningOrderByID(id uint64) error {
	ordMap := m.GetOrderIDMap(id)
	if ordMap == nil {
		return gorm.ErrRecordNotFound
	}

	mo := m.GetOpeningMarketMemory(ordMap.CustomerID)
	if mo == nil {
		return gorm.ErrRecordNotFound
	}

	cu := mo.GetOrderMemory(ordMap.MarketUUID)
	if cu == nil {
		return gorm.ErrRecordNotFound
	}
	cu.DeleteOrder(ordMap.OrderID)
	return nil
}

func (m *Memory) UpdateOpeningOrderByID(id uint64, price string) (*order.Order, error) {
	ordMap := m.GetOrderIDMap(id)
	if ordMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	mo := m.GetOpeningMarketMemory(ordMap.CustomerID)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cu := mo.GetOrderMemory(ordMap.MarketUUID)
	if cu == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord := cu.GetOrder(ordMap.OrderID)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord.Price = price
	cu.StoreOrder(ord)
	return ord, nil
}

func (m *Memory) GetOrderIDMap(ID uint64) *OrderMap {
	m.orderIDMtx.RLock()
	v := m.OrderIDMap[ID]
	m.orderIDMtx.RUnlock()
	return v
}

func (m *Memory) GetClientIDMap(clientID string) *OrderMap {
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
	CloseOrders   []*order.Order
	OpeningOrders []*order.Order
	orderMapMtx   sync.Mutex
	OrderMaps     map[uint64]*OrderMap
}

func (m *Memory) Snapshot() *Snapshot {
	snap := &Snapshot{
		OrderMaps:   make(map[uint64]*OrderMap),
		CloseOrders: make([]*order.Order, 0, 5000000),
	}

	//var wg sync.WaitGroup
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//	now := time.Now().UnixNano() / 1e6
	//	for _, om := range m.OpeningOrders {
	//		for _, co := range om.MarketOrders {
	//			for _, ord := range co.Orders {
	//				snap.OpeningOrders = append(snap.OpeningOrders, &order.Order{
	//					Id:           ord.Id,
	//					Price:        ord.Price,
	//					StopPrice:    ord.StopPrice,
	//					Amount:       ord.Amount,
	//					MarketUuid:   ord.MarketUuid,
	//					Side:         ord.Side,
	//					State:        ord.State,
	//					FilledAmount: ord.FilledAmount,
	//					FilledFees:   ord.FilledFees,
	//					AvgDealPrice: ord.AvgDealPrice,
	//					InsertedAt:   ord.InsertedAt,
	//					UpdatedAt:    ord.UpdatedAt,
	//					CustomerId:   ord.CustomerId,
	//					Type:         ord.Type,
	//					Bu:           ord.Bu,
	//					Source:       ord.Source,
	//					Operator:     ord.Operator,
	//					Ioc:          ord.Ioc,
	//				})
	//			}
	//		}
	//	}
	//	now2 := time.Now().UnixNano() / 1e6
	//	fmt.Println("opening order deep copy ", now2-now)
	//}()

	//go func() {
	//	defer wg.Done()
	now := time.Now().UnixNano() / 1e6
	for _, om := range m.CloseOrders {
		for _, co := range om.MarketOrders {
			for _, o := range co.Orders {
				//snap.CloseOrders = append(snap.CloseOrders, o)
				_ = o
			}
		}
	}
	now2 := time.Now().UnixNano() / 1e6
	fmt.Println("close order copy ", now2-now)
	//}()

	//go func() {
	m.orderIDMtx.RLock()
	for _, ordMap := range m.OrderIDMap {
		snap.OrderMaps[ordMap.OrderID] = ordMap
	}
	m.orderIDMtx.RUnlock()

	//go func() {
	m.clientIDMtx.RLock()
	for _, ordMap := range m.ClientIDMap {
		snap.OrderMaps[ordMap.OrderID] = ordMap
	}
	m.clientIDMtx.RUnlock()

	//wg.Wait()

	return snap
}

func TestMap() {
	db := NewMemory()
	//var clientIDs []string
	//clientID := uuid.New().String()
	Insert(db)
	//time.Sleep(time.Second * 10)
	fmt.Println("start to deep copy ")
	now3 := time.Now().UnixNano() / 1e6
	snap := db.Snapshot()
	now4 := time.Now().UnixNano() / 1e6
	fmt.Println("deep copy snapshot ", now4-now3)
	fmt.Println("close order len", len(snap.CloseOrders))
	fmt.Println("opening order len", len(snap.OpeningOrders))
	fmt.Println("map len", len(snap.OrderMaps))

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
	now := time.Now().UnixNano() / 1e6
	for i := 1; i <= num; i++ {
		clientID := uuid.New().String()
		//clientIDs = append(clientIDs, clientID)
		ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
		db.SetCloseOrder(ord)
		db.SetOrderMap(ord, clientID)
	}
	a := num + 50000
	for i := num + 1; i <= a; i++ {
		clientID := uuid.New().String()
		//clientIDs = append(clientIDs, clientID)
		ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, clientID)
	}

	now2 := time.Now().UnixNano() / 1e6
	fmt.Println("insert into eslap ", now2-now)
}
