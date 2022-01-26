package gobuntdb

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
	order2 "gomem/order"
	"gomem/protopb/order"
	"gomem/snapshot"
	"strconv"
	"time"
)

var num = 5000000

//var num = 50000

func NewBuntDb() *buntdb.DB {
	db, _ := buntdb.Open(":memory:")
	db.CreateIndex("id", "*", buntdb.IndexJSON("id"))
	db.CreateIndex("customer_id", "*", buntdb.IndexJSON("customer_id"))
	db.CreateIndex("market_uuid", "*", buntdb.IndexJSON("market_uuid"))
	db.CreateIndex("customer_id_market_uuid", "*", buntdb.IndexJSON("customer_id"), buntdb.IndexJSON("market_uuid"))

	return db
}
func NewBuntDbWithOutIndex() *buntdb.DB {
	db, _ := buntdb.Open(":memory:")
	db.CreateIndex("id", "*", buntdb.IndexJSON("id"))
	return db
}

func TestBuntDb() {
	//now := time.Now().UnixNano() / 1e6
	//fmt.Println(now)
	//fmt.Println(time.Now().Add(-time.Hour*24).UnixNano() / 1e6)
	//fmt.Println(now / 1000 / 60 / 60 / 24)
	db := NewBuntDb()
	orders := []*order.Order{
		{Id: 6, CustomerId: 1100752, MarketUuid: "ETH"},
		{Id: 8, CustomerId: 1100752, MarketUuid: "EOS"},
		{Id: 1, CustomerId: 1100752, MarketUuid: "BTC"},
		{Id: 3, CustomerId: 1100752, MarketUuid: "USDT"},
		{Id: 2, CustomerId: 1100750, MarketUuid: "BTC"},
		{Id: 4, CustomerId: 1100750, MarketUuid: "USDT"},
	}
	if err := db.Update(func(tx *buntdb.Tx) error {
		for _, o := range orders {
			data, err := json.Marshal(o)
			if err != nil {
				return err
			}
			_, _, err = tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	//
	//tx,_:=db.Begin(true)
	//tx.Commit()

	db.View(func(tx *buntdb.Tx) error {
		fmt.Println("Descend by id:")
		if err := tx.Descend("id", func(key, value string) bool {
			a := &order.Order{}
			err := json.Unmarshal([]byte(value), a)
			if err != nil {
				fmt.Println("==== ", err)
				return false
			}
			fmt.Printf("%v %+v\n", key, *a)
			return true
		}); err != nil {
			return err
		}

		fmt.Println("DescendEqual by customer_id:")
		if err := tx.DescendEqual("customer_id", `{"customer_id":1100752}`, func(key, value string) bool {
			a := &order.Order{}
			err := json.Unmarshal([]byte(value), a)
			if err != nil {
				fmt.Println("==== ", err)
				return false
			}
			fmt.Printf("%v %+v\n", key, *a)
			return true
		}); err != nil {
			panic(err)
			return err
		}

		fmt.Println("AscendEqual by customer_id market_uuid:")
		if err := tx.AscendEqual("customer_id_market_uuid", `{"customer_id":1100752,"market_uuid":"EOST"}`, func(key, value string) bool {
			a := &order.Order{}
			err := json.Unmarshal([]byte(value), a)
			if err != nil {
				fmt.Println("==== ", err)
				return false
			}
			fmt.Printf("%v %+v\n", key, *a)
			return true
		}); err != nil {
			panic(err)
			return err
		}
		return nil
	})
	//go func() {
	//	num := 0
	//	db.Update(func(tx *buntdb.Tx) error {
	//		fmt.Println("delete aaaaaa")
	//		now := time.Now().UnixNano() / 1e6
	//		if err := tx.DescendRange("id", `{"id":5000000}`, `{"id":4996000}`, func(key, value string) bool {
	//			a := &order.Order{}
	//			err := json.Unmarshal([]byte(value), a)
	//			if err != nil {
	//				return false
	//			}
	//			tx.Delete(strconv.Itoa(int(a.Id)))
	//			num++
	//			return true
	//		}); err != nil {
	//			panic(err)
	//			return err
	//		}
	//		eslapse := time.Now().UnixNano()/1e6 - now
	//		fmt.Println("delete sollect:", eslapse)
	//		fmt.Println("delete bbbbbbb")
	//		return nil
	//	})
	//	fmt.Println("delete number: ", num)
	//}()

	a := make(chan struct{})

	go func() {
		//time.Sleep(time.Millisecond * 10)
		fmt.Println("insert aaaaaaaa")
		now := time.Now().UnixNano() / 1e6
		db.Update(func(tx *buntdb.Tx) error {
			for i := 1000; i <= num+50000; i++ {
				strNum := strconv.Itoa(i)
				clientID := uuid.New().String()
				ord := order2.NewOrder(uint64(i), uint64(i), "2", "2", time.Now(), clientID)
				data, err := json.Marshal(ord)
				if err != nil {
					logrus.Errorf("marshal error is %+v", err)
					return err
				}
				if _, _, err = tx.Set(strNum, string(data), nil); err != nil {
					logrus.Errorf("set error is %+v", err)
					return err
				}
				//if i%100000 == 0 {
				//	fmt.Println("insert have:", i)
				//}
				//if err := fn(i, tx); err != nil {
				//	logrus.Errorf("error is %+v", err)
				//}
			}
			return nil
		})
		now2 := time.Now().UnixNano() / 1e6
		fmt.Println("insert escape ", now2-now)
		fmt.Println("insert bbbbbbb")
		a <- struct{}{}
	}()

	<-a
	fmt.Println("start save:")
	now := time.Now().UnixNano() / 1e6
	//f, err := os.Create("temp.db")
	//if err != nil {
	//	panic(err)
	//}
	//defer func() {
	//	f.Close()
	//	//os.RemoveAll("temp.db")
	//}()
	//if err = db.Save(f); err != nil {
	//	panic(err)
	//}
	//var saveOrders []*order.Order
	var ordersStr bytes.Buffer
	total := 0
	db.View(func(tx *buntdb.Tx) error {
		if err := tx.Ascend("id", func(key, value string) bool {
			//ord := &order.Order{}
			//err := json.Unmarshal([]byte(value), ord)
			//if err != nil {
			//	return false
			//}
			//saveOrders = append(saveOrders, ord)
			ordersStr.WriteString(",")
			ordersStr.WriteString(value)
			total++
			return true
		}); err != nil {
			return err
		}
		return nil
	})
	now2 := time.Now().UnixNano() / 1e6
	eslapse := now2 - now
	fmt.Println("save sollect:", eslapse)
	fmt.Println("orders len ", total)
	time.Sleep(time.Second * 10)
	now3 := time.Now().UnixNano() / 1e6
	//data, err := json.Marshal(saveOrders)
	//if err != nil {
	//	panic(err)
	//	return
	//}
	//fmt.Println("orders ", string(data))
	//fmt.Println("orders string", ordersStr)
	snapshot.BackUp(1, ordersStr.Bytes())
	now4 := time.Now().UnixNano() / 1e6
	eslapse2 := now4 - now3
	fmt.Println("save to file:", eslapse2)
	time.Sleep(time.Minute * 10)
}
