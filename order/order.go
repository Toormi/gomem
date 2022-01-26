package order

import (
	"gomem/protopb/order"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func NewOrder(id, customerID uint64, price, amount string, now time.Time, marketUUID string) *order.Order {
	return &order.Order{Id: id, CustomerId: customerID, MarketUuid: marketUUID, Price: price, Amount: amount, Side: order.Order_Side(id % 2), State: order.Order_PENDING, Type: order.Order_LIMIT, Ioc: false, FilledAmount: "0", AvgDealPrice: "0", InsertedAt: timestamppb.New(now), UpdatedAt: timestamppb.New(now)}
}
