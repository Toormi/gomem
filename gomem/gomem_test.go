package gomem

import (
	"github.com/hashicorp/go-memdb"
	xerrors "github.com/pkg/errors"
	"gomem/protopb/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func BenchmarkInsertWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn := db.Txn(true)
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
		txn.Commit()
	}
}

func BenchmarkInsertConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn := db.Txn(true)
			n := rand.Intn(num) + 1
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			txn.Insert("order", o)
			txn.Commit()
		}
	})

}

func BenchmarkGetWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(false)
		_, err = txn.Get("order", "id", uint64(n))
		if err != nil {
			b.Fatal(err)
		}
		txn.Abort()
	}
}

func BenchmarkGetConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(false)
			n := rand.Intn(num) + 1
			_, err := txn.Get("order", "id", uint64(n))
			if err != nil {
				b.Fatal(err)
			}
			txn.Abort()
		}
	})
}

func BenchmarkMultipleGetWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(false)
		_, err = txn.Get("order", "customer_id_market_uuid", uint64(n), "a")
		if err != nil {
			b.Fatal(err)
		}
		txn.Abort()
	}
}

func BenchmarkMultipleGetConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(false)
			n := rand.Intn(num) + 1
			_, err = txn.Get("order", "customer_id_market_uuid", uint64(n), "a")
			if err != nil {
				b.Fatal(err)
			}
			txn.Abort()
		}
	})
}

func BenchmarkConcurrentWriteAndGetWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		err = txn.Insert("order", o)
		if err != nil {
			b.Fatal(err)
		}
	}
	txn.Commit()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if n/2 == 0 {
				txn2 := db.Txn(true)
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
				err = txn2.Insert("order", o)
				if err != nil {
					b.Fatal(err)
				}
				txn2.Commit()
			} else {
				txn3 := db.Txn(false)
				_, err = txn3.Get("order", "id", uint64(n))
				if err != nil {
					b.Fatal(err)
				}
				txn3.Abort()
			}
		}
	})
}

func BenchmarkDeleteWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(true)
		_, err = txn.DeleteAll("order", "id", uint64(n))
		if err != nil {
			b.Fatal(err)
		}
		txn.Abort()
	}
}

func BenchmarkDeleteConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(true)
			n := rand.Intn(num) + 1
			_, err = txn.DeleteAll("order", "id", uint64(n))
			if err != nil {
				b.Fatal(err)
			}
			txn.Abort()
		}
	})
}

func BenchmarkUpdateWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(true)
		ord, err := txn.First("order", "id", uint64(n))
		if err != nil && !xerrors.Is(err, memdb.ErrNotFound) {
			b.Fatal(err)
		}

		if ord != nil {
			o := ord.(*order.Order)
			o.Price = "3"
			txn.Insert("order", o)
		}

		txn.Abort()
	}
}

func BenchmarkUpdateConcurrentWithIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMem()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(true)
			n := rand.Intn(num) + 1
			ord, err := txn.First("order", "id", uint64(n))
			if err != nil && !xerrors.Is(err, memdb.ErrNotFound) {
				b.Fatal(err)
			}

			if ord != nil {
				o := ord.(*order.Order)
				o.Price = "3"
				txn.Insert("order", o)
			}
			txn.Abort()
		}
	})
}

func BenchmarkInsertWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn := db.Txn(true)
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
		txn.Commit()
	}
}

func BenchmarkInsertConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn := db.Txn(true)
			n := rand.Intn(num) + 1
			strNum := strconv.Itoa(int(n))
			o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
			txn.Insert("order", o)
			txn.Commit()
		}
	})

}

func BenchmarkGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(false)
		_, err = txn.Get("order", "id", uint64(n))
		if err != nil {
			b.Fatal(err)
		}
		txn.Abort()
	}
}

func BenchmarkGetConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(false)
			n := rand.Intn(num) + 1
			_, err := txn.Get("order", "id", uint64(n))
			if err != nil {
				b.Fatal(err)
			}
			txn.Abort()
		}
	})
}

func BenchmarkConcurrentWriteAndGetWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		err = txn.Insert("order", o)
		if err != nil {
			b.Fatal(err)
		}
	}
	txn.Commit()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if n/2 == 0 {
				txn2 := db.Txn(true)
				strNum := strconv.Itoa(int(n))
				o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
				err = txn2.Insert("order", o)
				if err != nil {
					b.Fatal(err)
				}
				txn2.Commit()
			} else {
				txn3 := db.Txn(false)
				_, err = txn3.Get("order", "id", uint64(n))
				if err != nil {
					b.Fatal(err)
				}
				txn3.Abort()
			}
		}
	})
}

func BenchmarkDeleteWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(true)
		_, err = txn.DeleteAll("order", "id", uint64(n))
		if err != nil {
			b.Fatal(err)
		}
		txn.Abort()
	}
}

func BenchmarkDeleteConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(true)
			n := rand.Intn(num) + 1
			_, err = txn.DeleteAll("order", "id", uint64(n))
			if err != nil {
				b.Fatal(err)
			}
			txn.Abort()
		}
	})
}

func BenchmarkUpdateWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		txn = db.Txn(true)
		ord, err := txn.First("order", "id", uint64(n))
		if err != nil && !xerrors.Is(err, memdb.ErrNotFound) {
			b.Fatal(err)
		}

		if ord != nil {
			o := ord.(*order.Order)
			o.Price = "3"
			txn.Insert("order", o)
		}

		txn.Abort()
	}
}

func BenchmarkUpdateConcurrentWithOutIndex(b *testing.B) {
	now := time.Now()
	db, err := NewGoMemWithOutIndex()
	if err != nil {
		b.Fatal(err)
	}
	txn := db.Txn(true)
	for n := 1; n <= num; n++ {
		strNum := strconv.Itoa(int(n))
		o := &order.Order{Id: uint64(n), CustomerId: uint64(n), MarketUuid: "a", Price: strNum, Amount: strNum, Side: order.Order_Side(n % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
		txn.Insert("order", o)
	}
	txn.Commit()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			txn = db.Txn(true)
			n := rand.Intn(num) + 1
			ord, err := txn.First("order", "id", uint64(n))
			if err != nil && !xerrors.Is(err, memdb.ErrNotFound) {
				b.Fatal(err)
			}

			if ord != nil {
				o := ord.(*order.Order)
				o.Price = "3"
				txn.Insert("order", o)
			}
			txn.Abort()
		}
	})
}
