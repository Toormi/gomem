// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peatio/smaug/v3/order/order.proto

package order

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Order_Side int32

const (
	Order_BID Order_Side = 0
	Order_ASK Order_Side = 1
)

var Order_Side_name = map[int32]string{
	0: "BID",
	1: "ASK",
}

var Order_Side_value = map[string]int32{
	"BID": 0,
	"ASK": 1,
}

func (x Order_Side) String() string {
	return proto.EnumName(Order_Side_name, int32(x))
}

func (Order_Side) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 0}
}

type Order_State int32

const (
	Order_PENDING   Order_State = 0
	Order_FILLED    Order_State = 1
	Order_CANCELLED Order_State = 2
	Order_FIRED     Order_State = 3
)

var Order_State_name = map[int32]string{
	0: "PENDING",
	1: "FILLED",
	2: "CANCELLED",
	3: "FIRED",
}

var Order_State_value = map[string]int32{
	"PENDING":   0,
	"FILLED":    1,
	"CANCELLED": 2,
	"FIRED":     3,
}

func (x Order_State) String() string {
	return proto.EnumName(Order_State_name, int32(x))
}

func (Order_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 1}
}

type Order_Type int32

const (
	Order_LIMIT       Order_Type = 0
	Order_MARKET      Order_Type = 1
	Order_STOP_LIMIT  Order_Type = 2
	Order_STOP_MARKET Order_Type = 3
)

var Order_Type_name = map[int32]string{
	0: "LIMIT",
	1: "MARKET",
	2: "STOP_LIMIT",
	3: "STOP_MARKET",
}

var Order_Type_value = map[string]int32{
	"LIMIT":       0,
	"MARKET":      1,
	"STOP_LIMIT":  2,
	"STOP_MARKET": 3,
}

func (x Order_Type) String() string {
	return proto.EnumName(Order_Type_name, int32(x))
}

func (Order_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 2}
}

type Order_BU int32

const (
	// busniess unit
	Order_SPOT   Order_BU = 0
	Order_MARGIN Order_BU = 1
)

var Order_BU_name = map[int32]string{
	0: "SPOT",
	1: "MARGIN",
}

var Order_BU_value = map[string]int32{
	"SPOT":   0,
	"MARGIN": 1,
}

func (x Order_BU) String() string {
	return proto.EnumName(Order_BU_name, int32(x))
}

func (Order_BU) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 3}
}

type Order_Source int32

const (
	// busniess unit
	Order_WEB    Order_Source = 0
	Order_API    Order_Source = 1
	Order_SYSTEM Order_Source = 2
)

var Order_Source_name = map[int32]string{
	0: "WEB",
	1: "API",
	2: "SYSTEM",
}

var Order_Source_value = map[string]int32{
	"WEB":    0,
	"API":    1,
	"SYSTEM": 2,
}

func (x Order_Source) String() string {
	return proto.EnumName(Order_Source_name, int32(x))
}

func (Order_Source) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 4}
}

type Order_Operator int32

const (
	Order_LTE Order_Operator = 0
	Order_GTE Order_Operator = 1
)

var Order_Operator_name = map[int32]string{
	0: "LTE",
	1: "GTE",
}

var Order_Operator_value = map[string]int32{
	"LTE": 0,
	"GTE": 1,
}

func (x Order_Operator) String() string {
	return proto.EnumName(Order_Operator_name, int32(x))
}

func (Order_Operator) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0, 5}
}

type ListOrdersRequest_Side int32

const (
	ListOrdersRequest_UNSPECIFIED ListOrdersRequest_Side = 0
	ListOrdersRequest_BID         ListOrdersRequest_Side = 1
	ListOrdersRequest_ASK         ListOrdersRequest_Side = 2
)

var ListOrdersRequest_Side_name = map[int32]string{
	0: "UNSPECIFIED",
	1: "BID",
	2: "ASK",
}

var ListOrdersRequest_Side_value = map[string]int32{
	"UNSPECIFIED": 0,
	"BID":         1,
	"ASK":         2,
}

func (x ListOrdersRequest_Side) String() string {
	return proto.EnumName(ListOrdersRequest_Side_name, int32(x))
}

func (ListOrdersRequest_Side) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1, 0}
}

type ListOrdersRequest_State int32

const (
	ListOrdersRequest_UNSPECIFIEDSTATE ListOrdersRequest_State = 0
	ListOrdersRequest_PENDING          ListOrdersRequest_State = 1
	ListOrdersRequest_FILLED           ListOrdersRequest_State = 2
	ListOrdersRequest_CANCELLED        ListOrdersRequest_State = 3
	ListOrdersRequest_FIRED            ListOrdersRequest_State = 4
	ListOrdersRequest_OPENING          ListOrdersRequest_State = 5
)

var ListOrdersRequest_State_name = map[int32]string{
	0: "UNSPECIFIEDSTATE",
	1: "PENDING",
	2: "FILLED",
	3: "CANCELLED",
	4: "FIRED",
	5: "OPENING",
}

var ListOrdersRequest_State_value = map[string]int32{
	"UNSPECIFIEDSTATE": 0,
	"PENDING":          1,
	"FILLED":           2,
	"CANCELLED":        3,
	"FIRED":            4,
	"OPENING":          5,
}

func (x ListOrdersRequest_State) String() string {
	return proto.EnumName(ListOrdersRequest_State_name, int32(x))
}

func (ListOrdersRequest_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1, 1}
}

type ListOrdersRequest_Type int32

const (
	ListOrdersRequest_UNSPECIFIEDTYPE ListOrdersRequest_Type = 0
	ListOrdersRequest_LIMIT           ListOrdersRequest_Type = 1
	ListOrdersRequest_MARKET          ListOrdersRequest_Type = 2
	ListOrdersRequest_STOP_LIMIT      ListOrdersRequest_Type = 3
	ListOrdersRequest_STOP_MARKET     ListOrdersRequest_Type = 4
)

var ListOrdersRequest_Type_name = map[int32]string{
	0: "UNSPECIFIEDTYPE",
	1: "LIMIT",
	2: "MARKET",
	3: "STOP_LIMIT",
	4: "STOP_MARKET",
}

var ListOrdersRequest_Type_value = map[string]int32{
	"UNSPECIFIEDTYPE": 0,
	"LIMIT":           1,
	"MARKET":          2,
	"STOP_LIMIT":      3,
	"STOP_MARKET":     4,
}

func (x ListOrdersRequest_Type) String() string {
	return proto.EnumName(ListOrdersRequest_Type_name, int32(x))
}

func (ListOrdersRequest_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1, 2}
}

type ListOrdersRequest_BU int32

const (
	ListOrdersRequest_UNSPECIFIEDBU ListOrdersRequest_BU = 0
	ListOrdersRequest_SPOT          ListOrdersRequest_BU = 1
	ListOrdersRequest_MARGIN        ListOrdersRequest_BU = 2
)

var ListOrdersRequest_BU_name = map[int32]string{
	0: "UNSPECIFIEDBU",
	1: "SPOT",
	2: "MARGIN",
}

var ListOrdersRequest_BU_value = map[string]int32{
	"UNSPECIFIEDBU": 0,
	"SPOT":          1,
	"MARGIN":        2,
}

func (x ListOrdersRequest_BU) String() string {
	return proto.EnumName(ListOrdersRequest_BU_name, int32(x))
}

func (ListOrdersRequest_BU) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1, 3}
}

type ListOrdersRequest_Source int32

const (
	ListOrdersRequest_UNSPECIFIEDSOURCE ListOrdersRequest_Source = 0
	ListOrdersRequest_WEB               ListOrdersRequest_Source = 1
	ListOrdersRequest_API               ListOrdersRequest_Source = 2
	ListOrdersRequest_SYSTEM            ListOrdersRequest_Source = 3
)

var ListOrdersRequest_Source_name = map[int32]string{
	0: "UNSPECIFIEDSOURCE",
	1: "WEB",
	2: "API",
	3: "SYSTEM",
}

var ListOrdersRequest_Source_value = map[string]int32{
	"UNSPECIFIEDSOURCE": 0,
	"WEB":               1,
	"API":               2,
	"SYSTEM":            3,
}

func (x ListOrdersRequest_Source) String() string {
	return proto.EnumName(ListOrdersRequest_Source_name, int32(x))
}

func (ListOrdersRequest_Source) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1, 4}
}

type Order struct {
	Id                   uint64               `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Price                string               `protobuf:"bytes,2,opt,name=price,proto3" json:"price,omitempty"`
	StopPrice            string               `protobuf:"bytes,3,opt,name=stop_price,json=stopPrice,proto3" json:"stop_price,omitempty"`
	Amount               string               `protobuf:"bytes,4,opt,name=amount,proto3" json:"amount,omitempty"`
	MarketUuid           string               `protobuf:"bytes,5,opt,name=market_uuid,json=marketUuid,proto3" json:"market_uuid,omitempty"`
	Side                 Order_Side           `protobuf:"varint,6,opt,name=side,proto3,enum=peatio.smaug.v3.Order_Side" json:"side,omitempty"`
	State                Order_State          `protobuf:"varint,7,opt,name=state,proto3,enum=peatio.smaug.v3.Order_State" json:"state,omitempty"`
	FilledAmount         string               `protobuf:"bytes,8,opt,name=filled_amount,json=filledAmount,proto3" json:"filled_amount,omitempty"`
	FilledFees           string               `protobuf:"bytes,9,opt,name=filled_fees,json=filledFees,proto3" json:"filled_fees,omitempty"`
	AvgDealPrice         string               `protobuf:"bytes,10,opt,name=avg_deal_price,json=avgDealPrice,proto3" json:"avg_deal_price,omitempty"`
	InsertedAt           *timestamp.Timestamp `protobuf:"bytes,11,opt,name=inserted_at,json=insertedAt,proto3" json:"inserted_at,omitempty"`
	UpdatedAt            *timestamp.Timestamp `protobuf:"bytes,12,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CustomerId           uint64               `protobuf:"varint,13,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Type                 Order_Type           `protobuf:"varint,14,opt,name=type,proto3,enum=peatio.smaug.v3.Order_Type" json:"type,omitempty"`
	Bu                   Order_BU             `protobuf:"varint,15,opt,name=bu,proto3,enum=peatio.smaug.v3.Order_BU" json:"bu,omitempty"`
	Source               Order_Source         `protobuf:"varint,16,opt,name=source,proto3,enum=peatio.smaug.v3.Order_Source" json:"source,omitempty"`
	Operator             Order_Operator       `protobuf:"varint,17,opt,name=operator,proto3,enum=peatio.smaug.v3.Order_Operator" json:"operator,omitempty"`
	Ioc                  bool                 `protobuf:"varint,18,opt,name=ioc,proto3" json:"ioc,omitempty"`
	ClientId             string               `protobuf:"bytes,19,opt,name=client_id,json=clientId,proto3" json:"client_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Order) Reset()         { *m = Order{} }
func (m *Order) String() string { return proto.CompactTextString(m) }
func (*Order) ProtoMessage()    {}
func (*Order) Descriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{0}
}

func (m *Order) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Order.Unmarshal(m, b)
}
func (m *Order) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Order.Marshal(b, m, deterministic)
}
func (m *Order) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Order.Merge(m, src)
}
func (m *Order) XXX_Size() int {
	return xxx_messageInfo_Order.Size(m)
}
func (m *Order) XXX_DiscardUnknown() {
	xxx_messageInfo_Order.DiscardUnknown(m)
}

var xxx_messageInfo_Order proto.InternalMessageInfo

func (m *Order) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Order) GetPrice() string {
	if m != nil {
		return m.Price
	}
	return ""
}

func (m *Order) GetStopPrice() string {
	if m != nil {
		return m.StopPrice
	}
	return ""
}

func (m *Order) GetAmount() string {
	if m != nil {
		return m.Amount
	}
	return ""
}

func (m *Order) GetMarketUuid() string {
	if m != nil {
		return m.MarketUuid
	}
	return ""
}

func (m *Order) GetSide() Order_Side {
	if m != nil {
		return m.Side
	}
	return Order_BID
}

func (m *Order) GetState() Order_State {
	if m != nil {
		return m.State
	}
	return Order_PENDING
}

func (m *Order) GetFilledAmount() string {
	if m != nil {
		return m.FilledAmount
	}
	return ""
}

func (m *Order) GetFilledFees() string {
	if m != nil {
		return m.FilledFees
	}
	return ""
}

func (m *Order) GetAvgDealPrice() string {
	if m != nil {
		return m.AvgDealPrice
	}
	return ""
}

func (m *Order) GetInsertedAt() *timestamp.Timestamp {
	if m != nil {
		return m.InsertedAt
	}
	return nil
}

func (m *Order) GetUpdatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.UpdatedAt
	}
	return nil
}

func (m *Order) GetCustomerId() uint64 {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

func (m *Order) GetType() Order_Type {
	if m != nil {
		return m.Type
	}
	return Order_LIMIT
}

func (m *Order) GetBu() Order_BU {
	if m != nil {
		return m.Bu
	}
	return Order_SPOT
}

func (m *Order) GetSource() Order_Source {
	if m != nil {
		return m.Source
	}
	return Order_WEB
}

func (m *Order) GetOperator() Order_Operator {
	if m != nil {
		return m.Operator
	}
	return Order_LTE
}

func (m *Order) GetIoc() bool {
	if m != nil {
		return m.Ioc
	}
	return false
}

func (m *Order) GetClientId() string {
	if m != nil {
		return m.ClientId
	}
	return ""
}

type ListOrdersRequest struct {
	CustomerId           uint64                   `protobuf:"varint,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	MarketUuid           string                   `protobuf:"bytes,2,opt,name=market_uuid,json=marketUuid,proto3" json:"market_uuid,omitempty"`
	State                ListOrdersRequest_State  `protobuf:"varint,3,opt,name=state,proto3,enum=peatio.smaug.v3.ListOrdersRequest_State" json:"state,omitempty"`
	Side                 ListOrdersRequest_Side   `protobuf:"varint,4,opt,name=side,proto3,enum=peatio.smaug.v3.ListOrdersRequest_Side" json:"side,omitempty"`
	InsertedAtStart      *timestamp.Timestamp     `protobuf:"bytes,5,opt,name=inserted_at_start,json=insertedAtStart,proto3" json:"inserted_at_start,omitempty"`
	InsertedAtEnd        *timestamp.Timestamp     `protobuf:"bytes,6,opt,name=inserted_at_end,json=insertedAtEnd,proto3" json:"inserted_at_end,omitempty"`
	PageToken            string                   `protobuf:"bytes,7,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	Limit                uint32                   `protobuf:"varint,8,opt,name=limit,proto3" json:"limit,omitempty"`
	Bu                   ListOrdersRequest_BU     `protobuf:"varint,9,opt,name=bu,proto3,enum=peatio.smaug.v3.ListOrdersRequest_BU" json:"bu,omitempty"`
	Source               ListOrdersRequest_Source `protobuf:"varint,10,opt,name=source,proto3,enum=peatio.smaug.v3.ListOrdersRequest_Source" json:"source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ListOrdersRequest) Reset()         { *m = ListOrdersRequest{} }
func (m *ListOrdersRequest) String() string { return proto.CompactTextString(m) }
func (*ListOrdersRequest) ProtoMessage()    {}
func (*ListOrdersRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{1}
}

func (m *ListOrdersRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListOrdersRequest.Unmarshal(m, b)
}
func (m *ListOrdersRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListOrdersRequest.Marshal(b, m, deterministic)
}
func (m *ListOrdersRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListOrdersRequest.Merge(m, src)
}
func (m *ListOrdersRequest) XXX_Size() int {
	return xxx_messageInfo_ListOrdersRequest.Size(m)
}
func (m *ListOrdersRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ListOrdersRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ListOrdersRequest proto.InternalMessageInfo

func (m *ListOrdersRequest) GetCustomerId() uint64 {
	if m != nil {
		return m.CustomerId
	}
	return 0
}

func (m *ListOrdersRequest) GetMarketUuid() string {
	if m != nil {
		return m.MarketUuid
	}
	return ""
}

func (m *ListOrdersRequest) GetState() ListOrdersRequest_State {
	if m != nil {
		return m.State
	}
	return ListOrdersRequest_UNSPECIFIEDSTATE
}

func (m *ListOrdersRequest) GetSide() ListOrdersRequest_Side {
	if m != nil {
		return m.Side
	}
	return ListOrdersRequest_UNSPECIFIED
}

func (m *ListOrdersRequest) GetInsertedAtStart() *timestamp.Timestamp {
	if m != nil {
		return m.InsertedAtStart
	}
	return nil
}

func (m *ListOrdersRequest) GetInsertedAtEnd() *timestamp.Timestamp {
	if m != nil {
		return m.InsertedAtEnd
	}
	return nil
}

func (m *ListOrdersRequest) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

func (m *ListOrdersRequest) GetLimit() uint32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ListOrdersRequest) GetBu() ListOrdersRequest_BU {
	if m != nil {
		return m.Bu
	}
	return ListOrdersRequest_UNSPECIFIEDBU
}

func (m *ListOrdersRequest) GetSource() ListOrdersRequest_Source {
	if m != nil {
		return m.Source
	}
	return ListOrdersRequest_UNSPECIFIEDSOURCE
}

type ListOrdersResponse struct {
	Orders               []*Order `protobuf:"bytes,1,rep,name=orders,proto3" json:"orders,omitempty"`
	PageToken            string   `protobuf:"bytes,2,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListOrdersResponse) Reset()         { *m = ListOrdersResponse{} }
func (m *ListOrdersResponse) String() string { return proto.CompactTextString(m) }
func (*ListOrdersResponse) ProtoMessage()    {}
func (*ListOrdersResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_25d95bbd5825932a, []int{2}
}

func (m *ListOrdersResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListOrdersResponse.Unmarshal(m, b)
}
func (m *ListOrdersResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListOrdersResponse.Marshal(b, m, deterministic)
}
func (m *ListOrdersResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListOrdersResponse.Merge(m, src)
}
func (m *ListOrdersResponse) XXX_Size() int {
	return xxx_messageInfo_ListOrdersResponse.Size(m)
}
func (m *ListOrdersResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListOrdersResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListOrdersResponse proto.InternalMessageInfo

func (m *ListOrdersResponse) GetOrders() []*Order {
	if m != nil {
		return m.Orders
	}
	return nil
}

func (m *ListOrdersResponse) GetPageToken() string {
	if m != nil {
		return m.PageToken
	}
	return ""
}

func init() {
	proto.RegisterEnum("peatio.smaug.v3.Order_Side", Order_Side_name, Order_Side_value)
	proto.RegisterEnum("peatio.smaug.v3.Order_State", Order_State_name, Order_State_value)
	proto.RegisterEnum("peatio.smaug.v3.Order_Type", Order_Type_name, Order_Type_value)
	proto.RegisterEnum("peatio.smaug.v3.Order_BU", Order_BU_name, Order_BU_value)
	proto.RegisterEnum("peatio.smaug.v3.Order_Source", Order_Source_name, Order_Source_value)
	proto.RegisterEnum("peatio.smaug.v3.Order_Operator", Order_Operator_name, Order_Operator_value)
	proto.RegisterEnum("peatio.smaug.v3.ListOrdersRequest_Side", ListOrdersRequest_Side_name, ListOrdersRequest_Side_value)
	proto.RegisterEnum("peatio.smaug.v3.ListOrdersRequest_State", ListOrdersRequest_State_name, ListOrdersRequest_State_value)
	proto.RegisterEnum("peatio.smaug.v3.ListOrdersRequest_Type", ListOrdersRequest_Type_name, ListOrdersRequest_Type_value)
	proto.RegisterEnum("peatio.smaug.v3.ListOrdersRequest_BU", ListOrdersRequest_BU_name, ListOrdersRequest_BU_value)
	proto.RegisterEnum("peatio.smaug.v3.ListOrdersRequest_Source", ListOrdersRequest_Source_name, ListOrdersRequest_Source_value)
	proto.RegisterType((*Order)(nil), "peatio.smaug.v3.Order")
	proto.RegisterType((*ListOrdersRequest)(nil), "peatio.smaug.v3.ListOrdersRequest")
	proto.RegisterType((*ListOrdersResponse)(nil), "peatio.smaug.v3.ListOrdersResponse")
}

func init() {
	proto.RegisterFile("peatio/smaug/v3/order/order.proto", fileDescriptor_25d95bbd5825932a)
}

var fileDescriptor_25d95bbd5825932a = []byte{
	// 977 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x96, 0x6d, 0x6f, 0xda, 0x48,
	0x10, 0xc7, 0xe3, 0x07, 0x08, 0x1e, 0x0a, 0x98, 0x6d, 0xaf, 0xf2, 0xa5, 0xad, 0xc2, 0xb9, 0xf7,
	0x40, 0x74, 0xaa, 0x91, 0x88, 0xfa, 0xe2, 0x2e, 0xba, 0x4a, 0x10, 0x4c, 0x64, 0x95, 0x00, 0xb2,
	0x8d, 0x4e, 0x39, 0xe9, 0x84, 0x0c, 0xde, 0x50, 0xab, 0x80, 0x7d, 0xf6, 0x3a, 0x52, 0xbf, 0xea,
	0x7d, 0x80, 0xfb, 0x1c, 0xa7, 0xdd, 0x35, 0xe0, 0x40, 0x93, 0xe6, 0x4d, 0xc4, 0xce, 0xfe, 0xff,
	0xeb, 0xf1, 0xcc, 0xfe, 0x26, 0x86, 0x1f, 0x22, 0xec, 0x91, 0x20, 0x6c, 0x25, 0x2b, 0x2f, 0x5d,
	0xb4, 0xee, 0xce, 0x5b, 0x61, 0xec, 0xe3, 0x98, 0xff, 0x35, 0xa2, 0x38, 0x24, 0x21, 0xaa, 0x71,
	0x89, 0xc1, 0x24, 0xc6, 0xdd, 0xf9, 0xc9, 0xe9, 0x22, 0x0c, 0x17, 0x4b, 0xdc, 0x62, 0xdb, 0xb3,
	0xf4, 0xb6, 0x45, 0x82, 0x15, 0x4e, 0x88, 0xb7, 0x8a, 0xb8, 0x43, 0xff, 0xef, 0x18, 0x0a, 0x23,
	0x7a, 0x02, 0xaa, 0x82, 0x18, 0xf8, 0x9a, 0xd0, 0x10, 0x9a, 0xb2, 0x2d, 0x06, 0x3e, 0x7a, 0x01,
	0x85, 0x28, 0x0e, 0xe6, 0x58, 0x13, 0x1b, 0x42, 0x53, 0xb1, 0xf9, 0x02, 0xbd, 0x01, 0x48, 0x48,
	0x18, 0x4d, 0xf9, 0x96, 0xc4, 0xb6, 0x14, 0x1a, 0x19, 0xb3, 0xed, 0x97, 0x50, 0xf4, 0x56, 0x61,
	0xba, 0x26, 0x9a, 0xcc, 0xb6, 0xb2, 0x15, 0x3a, 0x85, 0xf2, 0xca, 0x8b, 0x3f, 0x63, 0x32, 0x4d,
	0xd3, 0xc0, 0xd7, 0x0a, 0x6c, 0x13, 0x78, 0x68, 0x92, 0x06, 0x3e, 0x6a, 0x81, 0x9c, 0x04, 0x3e,
	0xd6, 0x8a, 0x0d, 0xa1, 0x59, 0x6d, 0xbf, 0x32, 0xf6, 0x5e, 0xc4, 0x60, 0x39, 0x1a, 0x4e, 0xe0,
	0x63, 0x9b, 0x09, 0x51, 0x1b, 0x0a, 0x09, 0xf1, 0x08, 0xd6, 0x8e, 0x99, 0xe3, 0xf5, 0x43, 0x0e,
	0xaa, 0xb1, 0xb9, 0x14, 0xbd, 0x85, 0xca, 0x6d, 0xb0, 0x5c, 0x62, 0x7f, 0x9a, 0x25, 0x59, 0x62,
	0x79, 0x3c, 0xe3, 0xc1, 0xce, 0x36, 0xd5, 0x4c, 0x74, 0x8b, 0x71, 0xa2, 0x29, 0x3c, 0x55, 0x1e,
	0xea, 0x63, 0x9c, 0xa0, 0x1f, 0xa1, 0xea, 0xdd, 0x2d, 0xa6, 0x3e, 0xf6, 0x96, 0x59, 0x19, 0x80,
	0x1f, 0xe3, 0xdd, 0x2d, 0x7a, 0xd8, 0x5b, 0xf2, 0x4a, 0x5c, 0x40, 0x39, 0x58, 0x27, 0x38, 0x26,
	0xf4, 0x69, 0x44, 0x2b, 0x37, 0x84, 0x66, 0xb9, 0x7d, 0x62, 0xf0, 0x7e, 0x18, 0x9b, 0x7e, 0x18,
	0xee, 0xa6, 0x1f, 0x36, 0x6c, 0xe4, 0x1d, 0x82, 0x7e, 0x03, 0x48, 0x23, 0xdf, 0xcb, 0xbc, 0xcf,
	0xbe, 0xe9, 0x55, 0x32, 0x75, 0x87, 0xa5, 0x3f, 0x4f, 0x13, 0x12, 0xae, 0x70, 0x3c, 0x0d, 0x7c,
	0xad, 0xc2, 0xfa, 0x09, 0x9b, 0x90, 0xc5, 0x2a, 0x4d, 0xbe, 0x44, 0x58, 0xab, 0x3e, 0x5a, 0x69,
	0xf7, 0x4b, 0x84, 0x6d, 0x26, 0x44, 0x67, 0x20, 0xce, 0x52, 0xad, 0xc6, 0xe4, 0xdf, 0x3f, 0x20,
	0xef, 0x4e, 0x6c, 0x71, 0x96, 0xa2, 0xf7, 0x50, 0x4c, 0xc2, 0x34, 0x9e, 0x63, 0x4d, 0x65, 0xf2,
	0x37, 0x0f, 0x75, 0x85, 0x89, 0xec, 0x4c, 0x8c, 0x2e, 0xa0, 0x14, 0x46, 0x38, 0xf6, 0x48, 0x18,
	0x6b, 0x75, 0x66, 0x3c, 0x7d, 0xc0, 0x38, 0xca, 0x64, 0xf6, 0xd6, 0x80, 0x54, 0x90, 0x82, 0x70,
	0xae, 0xa1, 0x86, 0xd0, 0x2c, 0xd9, 0xf4, 0x27, 0x7a, 0x05, 0xca, 0x7c, 0x19, 0xe0, 0x35, 0xa1,
	0x05, 0x78, 0xce, 0x7a, 0x53, 0xe2, 0x01, 0xcb, 0xd7, 0x35, 0x90, 0xe9, 0x2d, 0x42, 0xc7, 0x20,
	0x75, 0xad, 0x9e, 0x7a, 0x44, 0x7f, 0x74, 0x9c, 0x8f, 0xaa, 0xa0, 0xff, 0x0e, 0x05, 0x76, 0x5b,
	0x50, 0x19, 0x8e, 0xc7, 0xe6, 0xb0, 0x67, 0x0d, 0xaf, 0xd4, 0x23, 0x04, 0x50, 0xec, 0x5b, 0x83,
	0x81, 0xd9, 0x53, 0x05, 0x54, 0x01, 0xe5, 0xb2, 0x33, 0xbc, 0x34, 0xd9, 0x52, 0x44, 0x0a, 0x14,
	0xfa, 0x96, 0x6d, 0xf6, 0x54, 0x49, 0xff, 0x00, 0x32, 0xad, 0x18, 0x0d, 0x0d, 0xac, 0x6b, 0xcb,
	0xe5, 0xc6, 0xeb, 0x8e, 0xfd, 0xd1, 0x74, 0x55, 0x01, 0x55, 0x01, 0x1c, 0x77, 0x34, 0x9e, 0xf2,
	0x3d, 0x11, 0xd5, 0xa0, 0xcc, 0xd6, 0x99, 0x40, 0xd2, 0x4f, 0x40, 0xec, 0x4e, 0x50, 0x09, 0x64,
	0x67, 0x3c, 0xda, 0x99, 0xaf, 0xac, 0xa1, 0x2a, 0xe8, 0x3f, 0x43, 0x91, 0xd7, 0x8b, 0xa6, 0xfa,
	0xa7, 0xd9, 0xcd, 0x72, 0x1e, 0x5b, 0xaa, 0x40, 0x75, 0xce, 0x8d, 0xe3, 0x9a, 0xd7, 0xaa, 0xa8,
	0xbf, 0x86, 0xd2, 0xa6, 0x3c, 0x54, 0x30, 0x70, 0x4d, 0xae, 0xbc, 0x72, 0x4d, 0x55, 0xd0, 0xff,
	0x2d, 0x42, 0x7d, 0x10, 0x24, 0x84, 0xd5, 0x31, 0xb1, 0xf1, 0x3f, 0x29, 0x4e, 0x0e, 0x6e, 0x8b,
	0x70, 0x70, 0x5b, 0xf6, 0xc0, 0x15, 0x0f, 0xc0, 0xfd, 0xb0, 0xe1, 0x50, 0x62, 0x8d, 0x6b, 0x1e,
	0x34, 0xee, 0xe0, 0xa1, 0xf7, 0x99, 0xbc, 0xc8, 0xc0, 0x97, 0x99, 0xfd, 0x97, 0xa7, 0xd8, 0x77,
	0x43, 0xa0, 0x0f, 0xf5, 0x1c, 0x64, 0xd3, 0x84, 0x78, 0x31, 0x61, 0xc3, 0xe5, 0x71, 0x5c, 0x6a,
	0x3b, 0xd4, 0x1c, 0x6a, 0x41, 0x5d, 0xa8, 0xe5, 0xcf, 0xc1, 0x6b, 0x9f, 0x0d, 0xa2, 0xc7, 0x4f,
	0xa9, 0xec, 0x4e, 0x31, 0xd7, 0x3e, 0x9d, 0x8c, 0x91, 0xb7, 0xc0, 0x53, 0x12, 0x7e, 0xc6, 0x6b,
	0x36, 0x95, 0x14, 0x5b, 0xa1, 0x11, 0x97, 0x06, 0xe8, 0x38, 0x5d, 0x06, 0xab, 0x80, 0xcf, 0x9c,
	0x8a, 0xcd, 0x17, 0xe8, 0x3d, 0x63, 0x4b, 0x61, 0xef, 0xfe, 0xd3, 0x13, 0xde, 0x3d, 0xe3, 0xac,
	0xb3, 0xe5, 0x0c, 0x98, 0xf5, 0xec, 0x29, 0x65, 0xbb, 0xc7, 0x9c, 0x7e, 0x96, 0x71, 0x50, 0x83,
	0xf2, 0x64, 0xe8, 0x8c, 0xcd, 0x4b, 0xab, 0x6f, 0x99, 0x19, 0x0f, 0x14, 0x0c, 0x61, 0x03, 0x86,
	0xa8, 0xff, 0xbd, 0x01, 0xe3, 0x05, 0xa8, 0x39, 0xad, 0xe3, 0x76, 0xd8, 0x15, 0xcb, 0xe1, 0x22,
	0xe4, 0x70, 0x11, 0xef, 0xe3, 0x22, 0xed, 0x70, 0x91, 0xa9, 0x65, 0x34, 0x36, 0x87, 0xd4, 0x52,
	0xd0, 0x9d, 0x8c, 0x9d, 0xe7, 0x50, 0xcb, 0x9d, 0xee, 0xde, 0x8c, 0xe9, 0xe1, 0x5b, 0xa0, 0x84,
	0x1c, 0x50, 0xe2, 0x1e, 0x50, 0xd2, 0x3e, 0x50, 0xb2, 0xfe, 0x8e, 0x01, 0x55, 0x87, 0x4a, 0xee,
	0xc8, 0xee, 0x44, 0x3d, 0xda, 0x32, 0x26, 0xe4, 0x18, 0x13, 0xf5, 0x3f, 0xb6, 0x8c, 0x7d, 0x07,
	0xf5, 0xfc, 0x3b, 0x8e, 0x26, 0xf6, 0x65, 0xc6, 0x11, 0x45, 0x4f, 0xd8, 0xa0, 0x27, 0xe6, 0xd0,
	0x93, 0xf4, 0x39, 0xa0, 0x7c, 0xc1, 0x93, 0x28, 0x5c, 0x27, 0x18, 0x19, 0x50, 0x64, 0xff, 0x9c,
	0x13, 0x4d, 0x68, 0x48, 0xcd, 0x72, 0xfb, 0xe5, 0xd7, 0x87, 0x9a, 0x9d, 0xa9, 0xf6, 0x6e, 0x90,
	0xb8, 0x77, 0x83, 0xda, 0x37, 0x50, 0xe4, 0x0f, 0x40, 0x23, 0x90, 0xe9, 0xe3, 0x90, 0xfe, 0xed,
	0xb6, 0x9f, 0xbc, 0x7d, 0x54, 0xc3, 0x33, 0xed, 0xbe, 0xfb, 0xeb, 0xd7, 0x45, 0x40, 0x3e, 0xa5,
	0x33, 0x63, 0x1e, 0xae, 0x5a, 0xd9, 0x77, 0xc6, 0xa7, 0x74, 0xd6, 0xfa, 0xea, 0x27, 0xc7, 0xac,
	0xc8, 0x68, 0x38, 0xff, 0x3f, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x4c, 0x25, 0xd3, 0x92, 0x08, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// OrdersClient is the client API for Orders service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type OrdersClient interface {
	List(ctx context.Context, in *ListOrdersRequest, opts ...grpc.CallOption) (*ListOrdersResponse, error)
}

type ordersClient struct {
	cc grpc.ClientConnInterface
}

func NewOrdersClient(cc grpc.ClientConnInterface) OrdersClient {
	return &ordersClient{cc}
}

func (c *ordersClient) List(ctx context.Context, in *ListOrdersRequest, opts ...grpc.CallOption) (*ListOrdersResponse, error) {
	out := new(ListOrdersResponse)
	err := c.cc.Invoke(ctx, "/peatio.smaug.v3.Orders/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OrdersServer is the server API for Orders service.
type OrdersServer interface {
	List(context.Context, *ListOrdersRequest) (*ListOrdersResponse, error)
}

// UnimplementedOrdersServer can be embedded to have forward compatible implementations.
type UnimplementedOrdersServer struct {
}

func (*UnimplementedOrdersServer) List(ctx context.Context, req *ListOrdersRequest) (*ListOrdersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

func RegisterOrdersServer(s *grpc.Server, srv OrdersServer) {
	s.RegisterService(&_Orders_serviceDesc, srv)
}

func _Orders_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListOrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OrdersServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/peatio.smaug.v3.Orders/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OrdersServer).List(ctx, req.(*ListOrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Orders_serviceDesc = grpc.ServiceDesc{
	ServiceName: "peatio.smaug.v3.Orders",
	HandlerType: (*OrdersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Orders_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "peatio/smaug/v3/order/order.proto",
}
