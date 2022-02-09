package _map

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	xerrors "github.com/pkg/errors"
	"gomem/protopb/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"testing"
	"time"
)

var (
	//cids  = []uint64{1100752, 1100750, 1100749, 1100758}
	uuids = []string{"ffd79462-48b6-43f7-8cb0-12bdeaf8915e", "ff447f95-b7cc-4523-ac5b-10da6cf2253c", "feefa85f-99d4-4eda-97c7-8de4c0c45931", "fdf18ca8-1e30-4191-8e57-d3f7c0b92c82", "fd7d29b6-10dd-4ade-925c-3e1f61ef9a3a"}
)

func BenchmarkInsert(b *testing.B) {
	db := NewMemory()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i := n % 4
		j := n % 5

		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
}

func BenchmarkInsertConcurrent(b *testing.B) {
	db := NewMemory()
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			i := n % 4
			j := n % 5
			id := uuid.New().String()
			ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
			db.SetOpeningOrder(ord)
			db.SetOrderMap(ord, id)
		}
	})

}

func BenchmarkGet(b *testing.B) {
	db := NewMemory()
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := db.GetOpeningOrderByID(uint64(n % 10))
		if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetConcurrent(b *testing.B) {
	db := NewMemory()
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			_, err := db.GetOpeningOrderByID(uint64(n % 10))
			if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkMultipleGet(b *testing.B) {
	db := NewMemory()
	var clientIDs []string
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
		clientIDs = append(clientIDs, id)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i := n % num
		if i >= num {
			i = 3
		}
		clientID := clientIDs[i]
		_, err := db.GetOpeningOrdersByCidAndMarketUUID(uint64(n), clientID)
		if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
			b.Fatal(err)
		}
	}
}

func BenchmarkMultipleGetConcurrent(b *testing.B) {
	db := NewMemory()
	var clientIDs []string
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
		clientIDs = append(clientIDs, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			i := n % num
			if i <= 0 || i >= num {
				i = 3
			}
			clientID := clientIDs[i]
			_, err := db.GetOpeningOrdersByCidAndMarketUUID(uint64(n), clientID)
			if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkConcurrentWriteAndGet(b *testing.B) {
	db := NewMemory()
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			if n/2 == 0 {
				i := n % 4
				j := n % 5
				id := uuid.New().String()
				ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
				db.SetOpeningOrder(ord)
				db.SetOrderMap(ord, id)
			} else {
				_, err := db.GetOpeningOrderByID(uint64(n % 10))
				if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
					b.Fatal(err)
				}
			}
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	db := NewMemory()
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i := rand.Intn(num) + 1
		err := db.DeleteOpeningOrderByID(uint64(i))
		if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeleteConcurrent(b *testing.B) {
	db := NewMemory()
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			err := db.DeleteOpeningOrderByID(uint64(n))
			if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUpdate(b *testing.B) {
	db := NewMemory()
	var clientIDs []string
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
		clientIDs = append(clientIDs, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		i := rand.Intn(num) + 1
		_, err := db.UpdateOpeningOrderByID(uint64(i), "12")
		if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
			b.Fatal(err)
		}
	}
}

func BenchmarkUpdateConcurrent(b *testing.B) {
	db := NewMemory()
	var clientIDs []string
	for n := 1; n <= num; n++ {
		i := n % 4
		j := n % 5
		id := uuid.New().String()
		ord := NewOrder(uint64(n), cids[i], "2", "2", time.Now(), uuids[j])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, id)
		clientIDs = append(clientIDs, id)
	}
	rand.Seed(time.Now().UnixNano())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := rand.Intn(num) + 1
			_, err := db.UpdateOpeningOrderByID(uint64(n), "12")
			if err != nil && !xerrors.Is(err, gorm.ErrRecordNotFound) {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkTestMap(b *testing.B) {
	db := NewMemory()
	Insert(db)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		db.Snapshot()
	}
}

func NewOrder(id, customerID uint64, price, amount string, now time.Time, marketUUID string) *order.Order {
	return &order.Order{Id: id, CustomerId: customerID, MarketUuid: marketUUID, Price: price, Amount: amount, Side: order.Order_Side(id % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
}
