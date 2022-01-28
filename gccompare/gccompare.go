package gccompare

import (
	"fmt"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/tidwall/buntdb"
	gobuntdb "gomem/buntdb"
	"gomem/gomem"
	_map "gomem/map"
	"gomem/order"
	orderpb "gomem/protopb/order"
	"math/rand"
	"runtime"
	"runtime/debug"
	"time"
)

var (
	cids        []uint64
	marketUUIDs []string
	customerNum = 50000
	marketNum   = 300
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

var previousPause time.Duration

func gcPause() time.Duration {
	runtime.GC()
	var stats debug.GCStats
	debug.ReadGCStats(&stats)
	pause := stats.PauseTotal - previousPause
	previousPause = stats.PauseTotal
	return pause
}

const (
	entries = 20000000
	repeat  = 50
)

func GcCompare() {
	debug.SetGCPercent(10)
	fmt.Println("Number of entries: ", entries)
	fmt.Println("Number of repeats: ", repeat)

	fmt.Println("GC pause for startup: ", gcPause())

	stdmapCache()
	//buntdbCache()
	//gomemdbCache()

	fmt.Println("GC pause for warmup: ", gcPause())
	//
	//for i := 0; i < repeat; i++ {
	//	gomemdbCache()
	//}
	//fmt.Println("GC pause for memdbcache: ", gcPause())
	//for i := 0; i < repeat; i++ {
	//	buntdbCache()
	//}
	//fmt.Println("GC pause for buntdbcache: ", gcPause())
	//for i := 0; i < repeat; i++ {
	stdmapCache()
	//}
	fmt.Println("GC pause for map: ", gcPause())
}

func gomemdbCache() {
	db, _ := gomem.NewGoMemWithOutIndex()
	tx := db.Txn(true)
	for i := 0; i < entries; i++ {
		m := rand.Int63n(int64(customerNum))
		n := rand.Int63n(int64(marketNum))
		_, ord := generateKeyAndValue(uint64(i), cids[m], "2", "2", time.Now(), marketUUIDs[n])
		tx.Insert("order", ord)
	}
	tx.Commit()
}

func buntdbCache() {
	db := gobuntdb.NewBuntDbWithOutIndex()
	db.Update(func(tx *buntdb.Tx) error {
		for i := 0; i < entries; i++ {

			m := rand.Int63n(int64(customerNum))
			n := rand.Int63n(int64(marketNum))
			key, ord := generateKeyAndValue(uint64(i), cids[m], "2", "2", time.Now(), marketUUIDs[n])
			data, err := json.Marshal(ord)
			if err != nil {
				return err
			}
			if _, _, err = tx.Set(key, string(data), nil); err != nil {
				return err
			}
		}
		return nil
	})
}

func stdmapCache() {
	db := _map.NewMemory()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < entries; i++ {
		m := rand.Int63n(int64(customerNum))
		n := rand.Int63n(int64(marketNum))
		clientID := uuid.New().String()
		_, ord := generateKeyAndValue(uint64(i), cids[m], "2", "2", time.Now(), marketUUIDs[n])
		db.SetOpeningOrder(ord)
		db.SetOrderMap(ord, clientID)
	}
}

func generateKeyAndValue(id uint64, customerID uint64, price, amount string, now time.Time, marketUUID string) (string, *orderpb.Order) {
	return key(id), order.NewPbOrder(id, customerID, price, amount, now, marketUUID)
}

func key(id uint64) string {
	return fmt.Sprintf("key-%010d", id)
}
