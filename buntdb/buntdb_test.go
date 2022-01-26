package gobuntdb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	json "github.com/json-iterator/go"
	xerrors "github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/tidwall/buntdb"
	"gomem/protopb/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkJsonInsertWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		data, _ := json.Marshal(o)
		if err := db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(strconv.Itoa(n), string(data), nil)
			if err != nil {
				panic(err)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonInsertConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := db.Update(func(tx *buntdb.Tx) error {
				n := rand.Intn(num) + 1
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}

				data, _ := json.Marshal(o)
				_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
				if err != nil {
					panic(err)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
	//for n := 0; n < b.N; n++ {
	//	o := &Order{ID: n, CustomerID: n, Price: decimal.NewFromFloat(float64(n))}
	//	txn.Insert("order", o)
	//}
}

func BenchmarkJsonGetWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strNum, string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.View(func(tx *buntdb.Tx) error {
			data, err := tx.Get(strconv.Itoa(n))
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = json.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonGetConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.View(func(tx *buntdb.Tx) error {
				data, err := tx.Get(strconv.Itoa(n))
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = json.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkJsonMultipleGetWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.View(func(tx *buntdb.Tx) error {
			strN := strconv.Itoa(n)
			if err := tx.AscendEqual("customer_id_market_uuid", `{"customer_id":`+strN+`,"market_uuid":"a"}`, func(key, value string) bool {
				a := &order.Order{}
				err := json.Unmarshal([]byte(value), a)
				if err != nil {
					fmt.Println("==== ", err)
					return false
				}
				return true
			}); err != nil {
				panic(err)
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonMultipleGetConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.View(func(tx *buntdb.Tx) error {
				strN := strconv.Itoa(n)
				if err := tx.AscendEqual("customer_id_market_uuid", `{"customer_id":`+strN+`,"market_uuid":"a"}`, func(key, value string) bool {
					a := &order.Order{}
					err := json.Unmarshal([]byte(value), a)
					if err != nil {
						fmt.Println("==== ", err)
						return false
					}
					return true
				}); err != nil {
					panic(err)
					return err
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkJsonConcurrentWriteAndGetWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if n/2 == 0 {
				if err := db.Update(func(tx *buntdb.Tx) error {
					strNum := strconv.Itoa(int(n))
					o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
					data, _ := json.Marshal(o)
					_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
					if err != nil {
						panic(err)
					}
					return nil
				}); err != nil {
					panic(err)
				}
			} else {
				if err := db.View(func(tx *buntdb.Tx) error {
					data, err := tx.Get(strconv.Itoa(int(n)))
					if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
						b.Fatal(err)
					}
					a := &order.Order{}
					if err == nil {
						if err = json.Unmarshal([]byte(data), a); err != nil {
							panic(err)
						}
					}
					return nil
				}); err != nil {
					panic(err)
				}
			}
		}
	})
}

func BenchmarkJsonDeleteWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonDeleteConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					return err
				}
				return nil
			}); err != nil {
				return
			}
		}
	})
}

func BenchmarkJsonUpdateWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strN := strconv.Itoa(n)
			data, err := tx.Get(strN)
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = json.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
				a.Price = "3"
				d, err := json.Marshal(a)
				if err != nil {
					return err
				}
				if _, _, err = tx.Set(strN, string(d), nil); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonUpdateConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				strN := strconv.Itoa(n)
				data, err := tx.Get(strN)
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = json.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
					a.Price = "3"
					d, err := json.Marshal(a)
					if err != nil {
						return err
					}
					if _, _, err = tx.Set(strN, string(d), nil); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkProtoInsertWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoInsertConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := db.Update(func(tx *buntdb.Tx) error {
				n := rand.Intn(num) + 1
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
				data, _ := proto.Marshal(o)
				_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
				if err != nil {
					panic(err)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
	//for n := 0; n < b.N; n++ {
	//	o := &Order{ID: n, CustomerID: n, Price: decimal.NewFromFloat(float64(n))}
	//	txn.Insert("order", o)
	//}
}

func BenchmarkProtoGetWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.View(func(tx *buntdb.Tx) error {
			data, err := tx.Get(strconv.Itoa(n))
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = proto.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoGetConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.View(func(tx *buntdb.Tx) error {
				data, err := tx.Get(strconv.Itoa(n))
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = proto.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkProtoConcurrentWriteAndGetWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := rand.Intn(100) + 1
			n := uint64(rand.Intn(m))
			if n/2 == 0 {
				if err := db.Update(func(tx *buntdb.Tx) error {
					o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: decimal.NewFromFloat(float64(n)).String()}
					data, _ := proto.Marshal(o)
					_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
					if err != nil {
						panic(err)
					}
					return nil
				}); err != nil {
					panic(err)
				}
			} else {
				if err := db.View(func(tx *buntdb.Tx) error {
					data, err := tx.Get(strconv.Itoa(int(n % 10)))
					if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
						b.Fatal(err)
					}
					a := &order.Order{}
					if err == nil {
						if err = proto.Unmarshal([]byte(data), a); err != nil {
							panic(err)
						}
					}
					return nil
				}); err != nil {
					panic(err)
				}
			}
		}
	})
}

func BenchmarkProtoDeleteWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				return err
			}
			return nil
		}); err != nil {
			return
		}
	}
}

func BenchmarkProtoDeleteConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					return err
				}
				return nil
			}); err != nil {
				return
			}
		}
	})
}

func BenchmarkProtoUpdateWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strN := strconv.Itoa(n)
			data, err := tx.Get(strN)
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = proto.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
				a.Price = "3"
				d, err := proto.Marshal(a)
				if err != nil {
					return err
				}
				if _, _, err = tx.Set(strN, string(d), nil); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoUpdateConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDb()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				strN := strconv.Itoa(n)
				data, err := tx.Get(strN)
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = proto.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
					a.Price = "3"
					d, err := proto.Marshal(a)
					if err != nil {
						return err
					}
					if _, _, err = tx.Set(strN, string(d), nil); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkJsonInsertWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	b.ResetTimer()
	//fmt.Println(string(data))
	for n := 0; n < b.N; n++ {
		//s := strconv.Itoa(n)
		//data := `{"id":` + s + `,"price":` + s + `,"market_uuid":"a","customer_id":` + s + `}`
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		data, _ := json.Marshal(o)
		if err := db.Update(func(tx *buntdb.Tx) error {
			_, _, err := tx.Set(strconv.Itoa(n), string(data), nil)
			if err != nil {
				panic(err)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonInsertConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := db.Update(func(tx *buntdb.Tx) error {
				n := rand.Intn(num) + 1
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}

				data, _ := json.Marshal(o)
				_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
				if err != nil {
					panic(err)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
	//for n := 0; n < b.N; n++ {
	//	o := &Order{ID: n, CustomerID: n, Price: decimal.NewFromFloat(float64(n))}
	//	txn.Insert("order", o)
	//}
}

func BenchmarkJsonGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.View(func(tx *buntdb.Tx) error {
			data, err := tx.Get(strconv.Itoa(n))
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = json.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonGetConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.View(func(tx *buntdb.Tx) error {
				data, err := tx.Get(strconv.Itoa(n))
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = json.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkJsonConcurrentWriteAndGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if n/2 == 0 {
				if err := db.Update(func(tx *buntdb.Tx) error {
					strNum := strconv.Itoa(int(n))
					o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
					data, _ := json.Marshal(o)
					_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
					if err != nil {
						panic(err)
					}
					return nil
				}); err != nil {
					panic(err)
				}
			} else {
				if err := db.View(func(tx *buntdb.Tx) error {
					data, err := tx.Get(strconv.Itoa(int(n)))
					if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
						b.Fatal(err)
					}
					a := &order.Order{}
					if err == nil {
						if err = json.Unmarshal([]byte(data), a); err != nil {
							panic(err)
						}
					}
					return nil
				}); err != nil {
					panic(err)
				}
			}
		}
	})
}

func BenchmarkJsonDeleteWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				return err
			}
			return nil
		}); err != nil {
			return
		}
	}
}

func BenchmarkJsonDeleteConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					return err
				}
				return nil
			}); err != nil {
				return
			}
		}
	})
}

func BenchmarkJsonUpdateWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strN := strconv.Itoa(n)
			data, err := tx.Get(strN)
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = json.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
				a.Price = "3"
				d, err := json.Marshal(a)
				if err != nil {
					return err
				}
				if _, _, err = tx.Set(strN, string(d), nil); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkJsonUpdateConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := json.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				strN := strconv.Itoa(n)
				data, err := tx.Get(strN)
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = json.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
					a.Price = "3"
					d, err := json.Marshal(a)
					if err != nil {
						return err
					}
					if _, _, err = tx.Set(strN, string(d), nil); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkProtoInsertWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoInsertConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := db.Update(func(tx *buntdb.Tx) error {
				n := rand.Intn(num) + 1
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
				data, _ := proto.Marshal(o)
				_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
				if err != nil {
					panic(err)
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
	//for n := 0; n < b.N; n++ {
	//	o := &Order{ID: n, CustomerID: n, Price: decimal.NewFromFloat(float64(n))}
	//	txn.Insert("order", o)
	//}
}

func BenchmarkProtoGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.View(func(tx *buntdb.Tx) error {
			data, err := tx.Get(strconv.Itoa(n))
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = proto.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoGetConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.View(func(tx *buntdb.Tx) error {
				data, err := tx.Get(strconv.Itoa(n))
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = proto.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkProtoConcurrentWriteAndGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := rand.Intn(100) + 1
			n := uint64(rand.Intn(m))
			if n/2 == 0 {
				if err := db.Update(func(tx *buntdb.Tx) error {
					o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: decimal.NewFromFloat(float64(n)).String()}
					data, _ := proto.Marshal(o)
					_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
					if err != nil {
						panic(err)
					}
					return nil
				}); err != nil {
					panic(err)
				}
			} else {
				if err := db.View(func(tx *buntdb.Tx) error {
					data, err := tx.Get(strconv.Itoa(int(n % 10)))
					if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
						b.Fatal(err)
					}
					a := &order.Order{}
					if err == nil {
						if err = proto.Unmarshal([]byte(data), a); err != nil {
							panic(err)
						}
					}
					return nil
				}); err != nil {
					panic(err)
				}
			}
		}
	})
}

func BenchmarkProtoDeleteWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoDeleteConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				if _, err := tx.Delete(strconv.Itoa(n)); err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					return err
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func BenchmarkProtoUpdateWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		if err := db.Update(func(tx *buntdb.Tx) error {
			strN := strconv.Itoa(n)
			data, err := tx.Get(strN)
			if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
				b.Fatal(err)
			}
			a := &order.Order{}
			if err == nil {
				if err = proto.Unmarshal([]byte(data), a); err != nil {
					panic(err)
				}
				a.Price = "3"
				d, err := proto.Marshal(a)
				if err != nil {
					return err
				}
				if _, _, err = tx.Set(strN, string(d), nil); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			panic(err)
		}
	}
}

func BenchmarkProtoUpdateConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db := NewBuntDbWithOutIndex()
	if err := db.Update(func(tx *buntdb.Tx) error {
		for n := 1; n <= num; n++ {
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			data, _ := proto.Marshal(o)
			_, _, err := tx.Set(strconv.Itoa(int(o.Id)), string(data), nil)
			if err != nil {
				panic(err)
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if err := db.Update(func(tx *buntdb.Tx) error {
				strN := strconv.Itoa(n)
				data, err := tx.Get(strN)
				if err != nil && !xerrors.Is(err, buntdb.ErrNotFound) {
					b.Fatal(err)
				}
				a := &order.Order{}
				if err == nil {
					if err = proto.Unmarshal([]byte(data), a); err != nil {
						panic(err)
					}
					a.Price = "3"
					d, err := proto.Marshal(a)
					if err != nil {
						return err
					}
					if _, _, err = tx.Set(strN, string(d), nil); err != nil {
						return err
					}
				}
				return nil
			}); err != nil {
				panic(err)
			}
		}
	})
}

func NewString(num string) string {
	return `{"id":` + num + `,"price":"` + num + `","amount":"` + num + `","market_uuid":"ETH","filled_amount":"0","avg_deal_price":"0","inserted_at":{"seconds":1643017494,"nanos":695600000},"updated_at":{"seconds":1643017494,"nanos":695600000},"customer_id":` + num + `}`
}
