package memory

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	json "github.com/json-iterator/go"
	order2 "gomem/order"
	"gomem/protopb/order"
	"gomem/snapshot"
	"sync"
	"time"
)

//var num = 50000
//
var num = 5000000

// Memory orders map[cid][market_uuid][id]order
type Memory struct {
	OpeningOrders sync.Map
	CloseOrders   sync.Map
	mtx           sync.RWMutex
	OrderIDMap    sync.Map
	ClientIDMap   sync.Map
}

type OrderMap struct {
	OrderID    uint64
	ClientID   string
	CustomerID uint64
	MarketUUID string
}

type MarketOrderMemory struct {
	MarketOrders sync.Map
}

type OrderMemory struct {
	Orders sync.Map
}

func NewMemory() *Memory {
	return &Memory{
		OpeningOrders: sync.Map{},
		CloseOrders:   sync.Map{},
		OrderIDMap:    sync.Map{},
		ClientIDMap:   sync.Map{},
	}
}

func (m *Memory) SetOpeningOrder(order *order.Order) {
	cm := &MarketOrderMemory{}
	ac, _ := m.OpeningOrders.LoadOrStore(order.CustomerId, cm)
	//if !ok {
	//	ac=
	//}
	mo := &OrderMemory{}
	acm, _ := ac.(*MarketOrderMemory).MarketOrders.LoadOrStore(order.MarketUuid, mo)
	acm.(*OrderMemory).Orders.Store(order.Id, order)
}

func (m *Memory) SetCloseOrder(order *order.Order) {
	cm := &MarketOrderMemory{}
	ac, _ := m.CloseOrders.LoadOrStore(order.CustomerId, cm)
	//if !ok {
	//	ac=
	//}
	mo := &OrderMemory{}
	acm, _ := ac.(*MarketOrderMemory).MarketOrders.LoadOrStore(order.MarketUuid, mo)
	acm.(*OrderMemory).Orders.Store(order.Id, order)
}

func (m *Memory) SetOrderMap(order *order.Order, clientID string) {
	orderMap := &OrderMap{
		OrderID:    order.Id,
		ClientID:   clientID,
		CustomerID: order.CustomerId,
		MarketUUID: order.MarketUuid,
	}
	m.OrderIDMap.Store(order.Id, orderMap)

	if clientID != "" {
		m.ClientIDMap.Store(clientID, orderMap)
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
	var orders []*order.Order
	mo, has := m.OpeningOrders.Load(cid)
	if !has {
		return orders, nil
	}

	cm, ok := mo.(*MarketOrderMemory).MarketOrders.Load(marketUUID)
	if !ok {
		return orders, nil
	}

	cm.(*OrderMemory).Orders.Range(func(key, value interface{}) bool {
		orders = append(orders, value.(*order.Order))
		return true
	})

	return orders, nil
}

func (m *Memory) GetOpeningOrdersByCidAndMarketUUIDs(cid uint64, marketUUIDs []string) ([]*order.Order, error) {
	var orders []*order.Order
	mo, has := m.OpeningOrders.Load(cid)
	if !has {
		return orders, nil
	}

	return GetOrders(mo.(*MarketOrderMemory), marketUUIDs)
}

func (m *Memory) GetCloseOrdersByCidAndMarketUUID(cid uint64, marketUUID string) ([]*order.Order, error) {
	var orders []*order.Order
	mo, has := m.CloseOrders.Load(cid)
	if !has {
		return orders, nil
	}

	cm, ok := mo.(*MarketOrderMemory).MarketOrders.Load(marketUUID)
	if !ok {
		return orders, nil
	}

	cm.(*OrderMemory).Orders.Range(func(key, value interface{}) bool {
		orders = append(orders, value.(*order.Order))
		return true
	})

	return orders, nil
}

func (m *Memory) GetCloseOrdersByCidAndMarketUUIDs(cid uint64, marketUUIDs []string) ([]*order.Order, error) {
	var orders []*order.Order
	mo, has := m.CloseOrders.Load(cid)
	if !has {
		return orders, nil
	}

	return GetOrders(mo.(*MarketOrderMemory), marketUUIDs)
}

func GetOrders(mo *MarketOrderMemory, marketUUIDs []string) ([]*order.Order, error) {
	var orders []*order.Order
	mMap := toMap(marketUUIDs)
	mLen := len(marketUUIDs)
	mo.MarketOrders.Range(func(marketUUID, val interface{}) bool {
		if _, ok := mMap[marketUUID.(string)]; ok || mLen <= 0 {
			orderMemory := val.(*OrderMemory)
			orderMemory.Orders.Range(func(key, value interface{}) bool {
				orders = append(orders, value.(*order.Order))
				return true
			})
		}
		return true
	})
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
	mo, ok := m.OpeningOrders.Load(orderMap.CustomerID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	cm, ok := mo.(*MarketOrderMemory).MarketOrders.Load(orderMap.MarketUUID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	o, ok := cm.(*OrderMemory).Orders.Load(orderMap.OrderID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return o.(*order.Order), nil
}

func (m *Memory) GetCloseByOrderMap(orderMap *OrderMap) (*order.Order, error) {
	mo, ok := m.CloseOrders.Load(orderMap.CustomerID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	cm, ok := mo.(*MarketOrderMemory).MarketOrders.Load(orderMap.MarketUUID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	o, ok := cm.(*OrderMemory).Orders.Load(orderMap.OrderID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return o.(*order.Order), nil
}

func (m *Memory) DeleteOpeningOrderByID(id uint64) error {
	ordMap := m.GetOrderIDMap(id)
	if ordMap == nil {
		return gorm.ErrRecordNotFound
	}

	mo, ok := m.OpeningOrders.Load(ordMap.CustomerID)
	if !ok {
		return gorm.ErrRecordNotFound
	}

	cu, ok := mo.(*MarketOrderMemory).MarketOrders.Load(ordMap.MarketUUID)
	if !ok {
		return gorm.ErrRecordNotFound
	}
	cu.(*OrderMemory).Orders.Delete(ordMap.OrderID)
	return nil
}

func (m *Memory) UpdateOpeningOrderByID(id uint64, price string) (*order.Order, error) {
	ordMap := m.GetOrderIDMap(id)
	if ordMap == nil {
		return nil, gorm.ErrRecordNotFound
	}

	mo, ok := m.OpeningOrders.Load(ordMap.CustomerID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}

	cu, ok := mo.(*MarketOrderMemory).MarketOrders.Load(ordMap.MarketUUID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	value, ok := cu.(*OrderMemory).Orders.Load(ordMap.OrderID)
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	ord := value.(*order.Order)
	ord.Price = price
	cu.(*OrderMemory).Orders.Store(ordMap.OrderID, ord)
	return ord, nil
}

func (m *Memory) GetOrderIDMap(ID uint64) *OrderMap {
	v, ok := m.OrderIDMap.Load(ID)
	if !ok {
		return nil
	}
	return v.(*OrderMap)
}

func (m *Memory) GetClientIDMap(clientID string) *OrderMap {
	v, ok := m.ClientIDMap.Load(clientID)
	if !ok {
		return nil
	}
	return v.(*OrderMap)
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

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		now := time.Now().UnixNano() / 1e6
		m.OpeningOrders.Range(func(key, value interface{}) bool {
			mo := value.(*MarketOrderMemory)
			mo.MarketOrders.Range(func(key, value interface{}) bool {
				co := value.(*OrderMemory)
				co.Orders.Range(func(key, value interface{}) bool {
					ord := value.(*order.Order)
					snap.OpeningOrders = append(snap.OpeningOrders, &order.Order{
						Id:           ord.Id,
						Price:        ord.Price,
						StopPrice:    ord.StopPrice,
						Amount:       ord.Amount,
						MarketUuid:   ord.MarketUuid,
						Side:         ord.Side,
						State:        ord.State,
						FilledAmount: ord.FilledAmount,
						FilledFees:   ord.FilledFees,
						AvgDealPrice: ord.AvgDealPrice,
						InsertedAt:   ord.InsertedAt,
						UpdatedAt:    ord.UpdatedAt,
						CustomerId:   ord.CustomerId,
						Type:         ord.Type,
						Bu:           ord.Bu,
						Source:       ord.Source,
						Operator:     ord.Operator,
						Ioc:          ord.Ioc,
					})
					return true
				})
				return true
			})
			return true
		})
		now2 := time.Now().UnixNano() / 1e6
		fmt.Println("opening order deep copy ", now2-now)
	}()

	go func() {
		defer wg.Done()
		now := time.Now().UnixNano() / 1e6
		m.CloseOrders.Range(func(key, value interface{}) bool {
			mo := value.(*MarketOrderMemory)
			mo.MarketOrders.Range(func(key, value interface{}) bool {
				co := value.(*OrderMemory)
				co.Orders.Range(func(key, value interface{}) bool {
					ord := value.(*order.Order)
					snap.CloseOrders = append(snap.CloseOrders, ord)
					return true
				})
				return true
			})
			return true
		})
		now2 := time.Now().UnixNano() / 1e6
		fmt.Println("close order deep copy ", now2-now)
	}()

	//go func() {
	m.OrderIDMap.Range(func(key, value interface{}) bool {
		ordMap := value.(*OrderMap)
		snap.OrderMaps[ordMap.OrderID] = ordMap
		return true
	})

	//go func() {
	m.ClientIDMap.Range(func(key, value interface{}) bool {
		ordMap := value.(*OrderMap)
		snap.OrderMaps[ordMap.OrderID] = ordMap
		return true
	})

	wg.Wait()

	return snap
}

func TestMemory() {
	db := NewMemory()
	//var clientIDs []string
	fmt.Println("start insert")
	//clientID := uuid.New().String()
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
	time.Sleep(time.Second * 10)
	fmt.Println("start to deep copy ")
	now3 := time.Now().UnixNano() / 1e6
	snap := db.Snapshot()
	now4 := time.Now().UnixNano() / 1e6
	fmt.Println("deep copy snapshot ", now4-now3)
	fmt.Println("close order len", len(snap.CloseOrders))
	fmt.Println("opening order len", len(snap.OpeningOrders))
	fmt.Println("map len", len(snap.OrderMaps))

	now5 := time.Now().UnixNano() / 1e6
	data, err := json.Marshal(snap)
	if err != nil {
		panic(err)
		return
	}
	snapshot.BackUp(1, data)
	now6 := time.Now().UnixNano() / 1e6
	fmt.Println("save to file:", now6-now5)
	time.Sleep(time.Minute * 10)
}
