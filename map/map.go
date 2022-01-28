package _map

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	order2 "gomem/order"
	"gomem/protopb/order"
	snapshotpb "gomem/protopb/snapshot"
	"gomem/snapshot"
	"math/rand"
	"sync"
	"time"
)

var (
	cids        []uint64
	marketUUIDs []string
	customerNum = 5000
	marketNum   = 300
	px          = decimal.NewFromFloat(1)
)

func init() {
	for i := 1; i <= customerNum; i++ {
		cids = append(cids, uint64(i))
	}
	fmt.Println("cid", len(cids))

	for i := 1; i <= marketNum; i++ {
		marketUUIDs = append(marketUUIDs, uuid.New().String())
	}
	fmt.Println("market uuid", len(marketUUIDs))
}

//var num = 50000
//
var num = 20000000

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
	if len(m.Orders) <= 0 {
		// 将旧map释放，确保gc能回收
		m.Orders = make(map[uint64]*order.Order)
	}
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
	ac := m.GetOrStoreCloseMarketMemory(order.CustomerId)
	om := ac.GetOrStoreOrderMemory(order.MarketUuid)
	om.StoreOrder(order)
}

func (m *Memory) SetOrderMap(order *order.Order, clientID string) {
	m.orderIDMtx.Lock()
	m.OrderIDMap[order.Id] = order
	m.orderIDMtx.Unlock()

	if clientID != "" {
		m.clientIDMtx.Lock()
		m.ClientIDMap[clientID] = order
		m.clientIDMtx.Unlock()
	}
}

func (m *Memory) GetOpeningOrderByID(ID uint64) (*order.Order, error) {
	ord := m.GetOrderIDMap(ID)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return ord, nil
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
	ord := m.GetClientIDMap(ClientID)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return ord, nil
}

func (m *Memory) GetCloseOrderByID(ID uint64) (*order.Order, error) {
	ord := m.GetOrderIDMap(ID)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return ord, nil
}

func (m *Memory) GetCloseOrderByClientID(ClientID string) (*order.Order, error) {
	ord := m.GetClientIDMap(ClientID)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	return ord, nil
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
	ordMap := m.GetOrderIDMap(id)
	if ordMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	mo := m.GetOpeningMarketMemory(ordMap.CustomerId)
	if mo == nil {
		return nil, gorm.ErrRecordNotFound
	}

	cu := mo.GetOrderMemory(ordMap.MarketUuid)
	if cu == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord := cu.GetOrder(ordMap.Id)
	if ord == nil {
		return nil, gorm.ErrRecordNotFound
	}

	ord.Price = price
	cu.StoreOrder(ord)
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
	CloseOrders   []*order.Order
	OpeningOrders []*order.Order
}

func (m *Memory) Snapshot() *snapshotpb.Snapshot {
	snap := &snapshotpb.Snapshot{}

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
	m.mtx.Lock()
	for _, o := range m.OrderIDMap {
		snap.Orders = append(snap.Orders, o)
	}
	m.mtx.Unlock()
	now2 := time.Now().UnixNano() / 1e6
	fmt.Println("close order copy ", now2-now)
	//}()

	//go func() {
	//m.orderIDMtx.RLock()
	//for _, ordMap := range m.OrderIDMap {
	//	snap.OrderMaps[ordMap.OrderID] = ordMap
	//}
	//m.orderIDMtx.RUnlock()

	//go func() {
	//m.clientIDMtx.RLock()
	//for _, ordMap := range m.ClientIDMap {
	//	snap.OrderMaps[ordMap.OrderID] = ordMap
	//}
	//m.clientIDMtx.RUnlock()

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
	fmt.Println("close order len", len(snap.Orders))
	//fmt.Println("opening order len", len(snap.OpeningOrders))

	for i := 0; i < 5; i++ {
		SaveToFile(snap)
	}
}

func Insert(db *Memory) {
	fmt.Println("start insert")
	rand.Seed(time.Now().UnixNano())
	now := time.Now().UnixNano() / 1e6
	for i := 1; i <= num; i++ {
		m := rand.Int63n(int64(customerNum))
		n := rand.Int63n(int64(marketNum))
		clientID := uuid.New().String()
		//clientIDs = append(clientIDs, clientID)
		ord := order2.NewPbOrder(uint64(i), cids[m], "2", "2", time.Now(), marketUUIDs[n])
		db.SetCloseOrder(ord)
		db.SetOrderMap(ord, clientID)
	}
	a := num + 50000
	for i := num + 1; i <= a; i++ {
		m := rand.Int63n(int64(customerNum))
		n := rand.Int63n(int64(marketNum))
		clientID := uuid.New().String()
		//clientIDs = append(clientIDs, clientID)
		ord := order2.NewPbOrder(uint64(i), cids[m], "2", "2", time.Now(), marketUUIDs[n])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, clientID)
	}

	now2 := time.Now().UnixNano() / 1e6
	fmt.Println("insert into eslap ", now2-now)
}

func SaveToFile(snap *snapshotpb.Snapshot) {
	now5 := time.Now().UnixNano() / 1e6
	data, err := proto.Marshal(snap)
	if err != nil {
		panic(err)
		return
	}
	snapshot.BackUp(1, data)
	now6 := time.Now().UnixNano() / 1e6
	fmt.Println("save to file:", now6-now5)
}
