package gomem

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-memdb"
	json "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	order2 "gomem/order"
	"gomem/protopb/order"
	"gomem/snapshot"
	"time"
)

var num = 5000000

// Person Create a sample struct
type Person struct {
	Id    uint64
	Email string
	Name  string
	Age   int
}

// Order is orders table struct
type Order struct {
	ID           uint64          `gorm:"primary_key" json:"id"`
	CustomerID   uint64          `json:"customer_id"`
	MarketUUID   string          `gorm:"column:market_uuid" json:"market_uuid"`
	Price        decimal.Decimal `sql:"type:decimal(32,16);" json:"price"`
	Amount       decimal.Decimal `sql:"type:decimal(32,16);" json:"amount"`
	FilledAmount decimal.Decimal `sql:"type:decimal(32,16);" json:"filled_amount"`
	Hidden       bool            `json:"hidden"`
	IOC          bool            `json:"ioc"`
	InsertedAt   time.Time       `json:"inserted_at"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"autoUpdateTime:milli"`
	AvgDealPrice decimal.Decimal `sql:"type:decimal(32,16);" json:"avg_deal_price"`
	MakerFeeRate decimal.Decimal `sql:"type:decimal(32,16);" json:"maker_fee_rate"`
	TakerFeeRate decimal.Decimal `sql:"type:decimal(32,16);" json:"taker_fee_rate"`
	LockedFunds  decimal.Decimal `sql:"type:decimal(32,16);" json:"locked_funds"`
	StopPrice    decimal.Decimal `sql:"type:decimal(32,16);" json:"stop_price"`
	FilledFees   decimal.Decimal `sql:"type:decimal(32,16);" json:"filled_fees"`
	PostOnly     bool            `json:"post_only"`
}

func NewGoMem() (*memdb.MemDB, error) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
					"email": &memdb.IndexSchema{
						Name:    "email",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Email"},
					},
					"age": &memdb.IndexSchema{
						Name:    "age",
						Unique:  false,
						Indexer: &memdb.IntFieldIndex{Field: "Age"},
					},
				},
			},
			"order": {
				Name: "order",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
					"customer_id": {
						Name:    "customer_id",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "CustomerId"},
					},
					"market_uuid": {
						Name:    "market_uuid",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "MarketUuid"},
					},
					"customer_id_market_uuid": {
						Name:   "customer_id_market_uuid",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							Indexes: []memdb.Indexer{
								&memdb.UintFieldIndex{Field: "CustomerId"},
								&memdb.StringFieldIndex{Field: "MarketUuid"},
							},
						},
					},
				},
			},
			"close_order": {
				Name: "close_order",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
					"customer_id": {
						Name:    "customer_id",
						Unique:  false,
						Indexer: &memdb.UintFieldIndex{Field: "CustomerId"},
					},
					"market_uuid": {
						Name:    "market_uuid",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "MarketUuid"},
					},
					"customer_id_market_uuid": {
						Name:   "customer_id_market_uuid",
						Unique: false,
						Indexer: &memdb.CompoundMultiIndex{
							Indexes: []memdb.Indexer{
								&memdb.UintFieldIndex{Field: "CustomerId"},
								&memdb.StringFieldIndex{Field: "MarketUuid"},
							},
						},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewGoMemWithOutIndex() (*memdb.MemDB, error) {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"person": &memdb.TableSchema{
				Name: "person",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
				},
			},
			"order": {
				Name: "order",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
				},
			},
			"close_order": {
				Name: "close_order",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UintFieldIndex{Field: "Id"},
					},
				},
			},
		},
	}

	// Create a new data base
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestMemDb() {
	db, err := NewGoMem()
	if err != nil {
		fmt.Println(err)
		return
	}
	txn := db.Txn(true)
	//// Insert some people
	//people := []*gomem.Person{
	//	{1, "joe@aol.com", "Joe", 30},
	//	{3, "lucy@aol.com", "Lucy", 35},
	//	{2, "tariq@aol.com", "Tariq", 21},
	//	{4, "dorothy@aol.com", "Dorothy", 53},
	//}
	//for _, p := range people {
	//	if err := txn.Insert("person", p); err != nil {
	//		panic(err)
	//	}
	//}
	timeNow := time.Now()
	orders := []*order.Order{
		order2.NewOrder(6, 1100752, "10", "10", timeNow, "ETH"),
		order2.NewOrder(10, 1100752, "10", "10", timeNow, "EOS"),
		order2.NewOrder(8, 1100752, "10", "10", timeNow, "EOS"),
		order2.NewOrder(1, 1100752, "10", "10", timeNow, "BTC"),
		order2.NewOrder(3, 1100752, "10", "10", timeNow, "USDT"),
		order2.NewOrder(2, 1100750, "10", "10", timeNow, "BTC"),
		order2.NewOrder(4, 1100750, "10", "10", timeNow, "USDT"),
	}
	for _, p := range orders {
		if err = txn.Insert("order", p); err != nil {
			panic(err)
		}
	}
	d, err := txn.First("order", "id", uint64(1))
	if err != nil {
		panic(err)
		return
	}
	d.(*order.Order).Price = "10"
	err = txn.Insert("order", d)
	if err != nil {
		panic(err)
		return
	}
	txn.Commit()
	db2 := db.Snapshot()

	txn = db.Txn(false)
	it, err := txn.Get("order", "customer_id", uint64(1100752))
	if err != nil {
		fmt.Println("err is ====", err)
		return
	}
	filter := memdb.NewFilterIterator(it, func(i interface{}) bool {
		a := i.(*order.Order)
		return a.Id < 1
	})
	for obj := filter.Next(); obj != nil; obj = filter.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %v   %v  %v \n", p.Id, p.MarketUuid, p.Price)
	}
	txn.Abort()
	txn = db.Txn(true)
	d, err = txn.First("order", "id", uint64(1))
	if err != nil {
		panic(err)
		return
	}
	d.(*order.Order).Price = "20"
	err = txn.Insert("order", d)
	if err != nil {
		panic(err)
		return
	}
	err = txn.Insert("order", order2.NewOrder(5, 1100752, "10", "10", timeNow, "ETH888"))
	if err != nil {
		panic(err)
		return
	}
	err = txn.Insert("order", order2.NewOrder(5, 1100752, "10", "10", timeNow, "ETH6"))
	if err != nil {
		panic(err)
		return
	}
	//effect, err := txn.DeletePrefix("order", "market_uuid_prefix", "DDD")
	//if err != nil {
	//	fmt.Println("delete prefix err ", err)
	//	return
	//}
	//fmt.Println("effect ", effect)
	txn.Commit()

	txn2 := db2.Txn(false)
	defer txn2.Abort()
	it, err = txn2.Get("order", "customer_id", uint64(1100752))
	if err != nil {
		fmt.Println("err is ====", err)
		return
	}
	filter = memdb.NewFilterIterator(it, func(i interface{}) bool {
		a := i.(*order.Order)
		return a.Id < 1
	})
	fmt.Println("snapshot object:")
	for obj := filter.Next(); obj != nil; obj = filter.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %v   %v  %v \n", p.Id, p.MarketUuid, p.Price)
	}

	txn = db.Txn(false)
	defer txn.Abort()
	it, err = txn.Get("order", "customer_id", uint64(1100752))
	if err != nil {
		fmt.Println("err is ====", err)
		return
	}
	filter = memdb.NewFilterIterator(it, func(i interface{}) bool {
		a := i.(*order.Order)
		return a.Id < 1
	})
	fmt.Println("row object:")
	for obj := filter.Next(); obj != nil; obj = filter.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %v   %v  %v \n", p.Id, p.MarketUuid, p.Price)
	}
	//Commit the transaction
	rows, err := db.Txn(false).Get("order", "id")
	if err != nil {
		fmt.Println("get rows failed ", err)
		return
	}
	for obj := rows.Next(); obj != nil; obj = rows.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %v   %v   %v \n", p.Id, p.MarketUuid, p.Price)
	}
	data, _ := json.Marshal(rows.Next())
	fmt.Println(string(data))
	//txn = db.Txn(false)
	//defer txn.Abort()
	//// Lookup by email
	//raw, err := txn.First("person", "email", "joe@aol.com")
	//if err != nil {
	//	panic(err)
	//}
	//
	//// Say hi!
	//fmt.Printf("Hello %s!\n", raw.(*gomem.Person).Name)
	//
	//// List all the people
	//it, err := txn.Get("person", "id")
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("All the people:")
	//for obj := it.Next(); obj != nil; obj = it.Next() {
	//	p := obj.(*gomem.Person)
	//	fmt.Printf("  %s   %d   \n", p.Name, p.Id)
	//}
	//
	//// Range scan over people with ages between 25 and 35 inclusive
	//it, err = txn.LowerBound("person", "age", 25)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("People aged 25 - 35:")
	//for obj := it.Next(); obj != nil; obj = it.Next() {
	//	p := obj.(*gomem.Person)
	//	if p.Age > 35 {
	//		break
	//	}
	//	fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
	//}

	fmt.Println("get by customer id market uuid:")

	it, err = txn.GetReverse("order", "customer_id_market_uuid", uint64(1100752), "ETH")
	if err != nil {
		panic(err)
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %+v \n", p)
	}
	txn = db.Txn(true)
	txn.DeleteAll("order", "customer_id_market_uuid", uint64(1100752), "ETH")
	txn.Commit()

	txn = db.Txn(false)
	defer txn.Abort()
	fmt.Println("delete get by customer_id market_uuid:")
	it, err = txn.GetReverse("order", "customer_id_market_uuid", uint64(1100752), "ETH")
	if err != nil {
		panic(err)
	}
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*order.Order)
		fmt.Printf("  %+v \n", p)
	}

	a := make(chan struct{})

	//go func() {
	//	time.Sleep(time.Millisecond * 20)
	//	num := 0
	//	txn := db.Txn(true)
	//	fmt.Println("delete remove lock")
	//	now := time.Now().UnixNano() / 1e6
	//	// Get all the objects
	//	iter, err := txn.LowerBound("close_order", "id", uint64(4996000))
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// Put them into a slice so there are no safety concerns while actually
	//	// performing the deletes
	//	for {
	//		obj := iter.Next()
	//		if obj == nil {
	//			break
	//		}
	//		if err = txn.Delete("close_order", obj); err != nil {
	//			panic(err)
	//			return
	//		}
	//		num++
	//	}
	//	eslapse := time.Now().UnixNano()/1e6 - now
	//	fmt.Println("delete sollect:", eslapse)
	//	txn.Commit()
	//	fmt.Println("delete number: ", num)
	//}()

	var clientIDs []string
	go func() {
		txn := db.Txn(true)
		fmt.Println("insert aaaaaaaa")
		now := time.Now().UnixNano() / 1e6
		num2 := num + 50000
		for i := num + 1; i <= num2; i++ {
			clientID := uuid.New().String()
			clientIDs = append(clientIDs, clientID)
			ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
			txn.Insert("order", ord)
		}
		for i := 1; i <= num; i++ {
			clientID := uuid.New().String()
			clientIDs = append(clientIDs, clientID)
			ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
			txn.Insert("close_order", ord)
			//if i%100000 == 0 {
			//	fmt.Println("insert have:", i)
			//}
		}
		now2 := time.Now().UnixNano() / 1e6
		fmt.Println("insert escape ", now2-now)
		fmt.Println("insert bbbbbbbb")
		txn.Commit()
		a <- struct{}{}
	}()

	<-a

	txn = db.Txn(true)
	fmt.Println("start save:")
	now := time.Now().UnixNano() / 1e6
	defer txn.Commit()
	it, err = txn.Get("order", "id")
	if err != nil {
		return
	}
	var saveOrders []*order.Order
	for obj := it.Next(); obj != nil; obj = it.Next() {
		ord := obj.(*order.Order)
		saveOrders = append(saveOrders, &order.Order{
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
	}
	db2 = db.Snapshot()

	eslapse := time.Now().UnixNano()/1e6 - now
	fmt.Println("save sollect:", eslapse)
	time.Sleep(time.Second * 10)

	total := 0
	now2 := time.Now().UnixNano() / 1e6
	tx := db2.Txn(false)
	defer tx.Abort()
	it, err = tx.Get("close_order", "id")
	if err != nil {
		panic(err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*order.Order)
		saveOrders = append(saveOrders, p)
		total++
	}

	fmt.Println("orders len ", len(saveOrders))
	data, err = json.Marshal(saveOrders)
	if err != nil {
		panic(err)
		return
	}
	snapshot.BackUp(1, data)

	now3 := time.Now().UnixNano() / 1e6
	fmt.Println("save to file:", now3-now2)

	time.Sleep(time.Minute)
}
