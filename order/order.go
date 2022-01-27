package order

import (
	"fmt"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"time"
)

// Type represents type of orders
type Type uint8

const (
	// TypeLimit represents limit order
	TypeLimit Type = iota
	// TypeMarket represents market order
	TypeMarket
	// TypeStopLimit represents stop market order
	TypeStopLimit
	// TypeStopMarket represents stop limit order
	TypeStopMarket
	// TypeAll represents all order
	TypeAll Type = 66
)

// UnmarshalJSON parse bytes to order type
func (t *Type) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	default:
		return errors.New("invalid order type")
	case "LIMIT":
		*t = TypeLimit
	case "MARKET":
		*t = TypeMarket
	case "STOP_MARKET":
		*t = TypeStopMarket
	case "STOP_LIMIT":
		*t = TypeStopLimit
	}

	return nil
}

// MarshalJSON parse order type to string
func (t Type) MarshalJSON() ([]byte, error) {
	var s string
	switch t {
	default:
		return nil, errors.New("invalid order type")
	case TypeLimit:
		s = "LIMIT"
	case TypeMarket:
		s = "MARKET"
	case TypeStopMarket:
		s = "STOP_MARKET"
	case TypeStopLimit:
		s = "STOP_LIMIT"
	}

	return json.Marshal(s)
}

// UnmarshalParam parse query/form param to order type in echo framework
func (t *Type) UnmarshalParam(src string) error {
	switch src {
	default:
		return errors.New("invalid order type")
	case "LIMIT":
		*t = TypeLimit
	case "MARKET":
		*t = TypeMarket
	case "STOP_MARKET":
		*t = TypeStopMarket
	case "STOP_LIMIT":
		*t = TypeStopLimit
	}

	return nil
}

// Side is order's side, one of ASK(sell), BID(buy)
type Side uint32

const (
	// SideBID is BID(buy) side
	SideBID Side = iota
	// SideASK is ASK(sell) side
	SideASK
	// SideALL is all side
	SideALL
)

// UnmarshalJSON parse bytes to order side
func (side *Side) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	default:
		return fmt.Errorf("invalid order side: %s", s)
	case "ASK":
		*side = SideASK
	case "BID":
		*side = SideBID
	}

	return nil
}

// MarshalJSON parse order side to string
func (side Side) MarshalJSON() ([]byte, error) {
	var s string
	switch side {
	default:
		return nil, errors.New("invalid order side")
	case SideASK:
		s = "ASK"
	case SideBID:
		s = "BID"
	}

	return json.Marshal(s)
}

// UnmarshalParam parse query/form param to order side in echo framework
func (side *Side) UnmarshalParam(src string) error {
	switch src {
	default:
		return fmt.Errorf("invalid order side: %s", src)
	case "ASK":
		*side = SideASK
	case "BID":
		*side = SideBID
	}

	return nil
}

// State is order state
type State uint32

const (
	// StatePending represents order state of pending
	StatePending State = iota
	// StateFilled represents order state of filled
	StateFilled
	// StateCancelled represents order state of cancelled
	StateCancelled
	// StateFired represents stop order state of created
	StateFired
	// StateRejected represents order state of rejected
	StateRejected
	// StateClosed represents order state of filled and cancelled
	StateClosed
	// StateOpening is opening state of pending and fired
	StateOpening
	// StateNofill is nofill state of closed that filled amount is zero
	StateNofill
	// StateALL is all state
	StateALL
)

// UnmarshalJSON parse bytes to order state
func (state *State) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	default:
		return errors.New("invalid order state")
	case "PENDING":
		*state = StatePending
	case "FILLED":
		*state = StateFilled
	case "CANCELLED":
		*state = StateCancelled
	case "CLOSED":
		*state = StateClosed
	case "FIRED":
		*state = StateFired
	case "REJECTED":
		*state = StateRejected
	}

	return nil
}

// MarshalJSON parse order state to string
func (state State) MarshalJSON() ([]byte, error) {
	var s string
	switch state {
	default:
		return nil, errors.New("invalid order state")
	case StatePending:
		s = "PENDING"
	case StateFilled:
		s = "FILLED"
	case StateCancelled:
		s = "CANCELLED"
	case StateClosed:
		s = "CLOSED"
	case StateFired:
		s = "FIRED"
	case StateRejected:
		s = "REJECTED"
	}

	return json.Marshal(s)
}

// UnmarshalParam parse query/form param to order state in echo framework
func (state *State) UnmarshalParam(src string) error {
	switch src {
	default:
		return errors.New("invalid order state")
	case "PENDING":
		*state = StatePending
	case "FILLED":
		*state = StateFilled
	case "CANCELLED":
		*state = StateCancelled
	case "CLOSED":
		*state = StateClosed
	case "FIRED":
		*state = StateFired
	case "REJECTED":
		*state = StateRejected
	case "OPENING":
		*state = StateOpening
	case "NONE_FILLED":
		*state = StateNofill
	case "ALL":
		*state = StateALL
	}

	return nil
}

// BusinessUnit represents `bu` field in table `orders` and `opening_orders`
type BusinessUnit uint8

const (
	// Spot represents spot account
	Spot BusinessUnit = iota
	// Margin represents margin account
	Margin
	// BUAll represents both of spot and margin account
	BUAll
)

// MarshalJSON parse orders `bu` field to string
func (b BusinessUnit) MarshalJSON() ([]byte, error) {
	var s string
	switch b {
	case Spot:
		s = "SPOT"
	case Margin:
		s = "MARGIN"
	default:
		return nil, errors.New("invalid business unit type")
	}
	return json.Marshal(s)
}

// UnmarshalJSON parse bytes to orders `bu` field
func (b *BusinessUnit) UnmarshalJSON(s []byte) error {
	var str string
	if err := json.Unmarshal(s, &str); err != nil {
		return err
	}
	switch str {
	case "SPOT":
		*b = Spot
	case "MARGIN":
		*b = Margin
	default:
		return errors.New("invalid business unit type")
	}
	return nil
}

// UnmarshalParam parse query/form parameters to order `bu` field
func (b *BusinessUnit) UnmarshalParam(str string) error {
	switch str {
	case "SPOT":
		*b = Spot
	case "MARGIN":
		*b = Margin
	default:
		return errors.New("invalid business unit type")
	}
	return nil
}

// Source represents `source` field in table `orders` and `opening_orders`
type Source uint8

const (
	// WEB represents order from web
	WEB Source = iota
	// API represents order from api
	API
	// SYSTEM represents order from system
	SYSTEM
	// SOURCEAll all source
	SOURCEAll
)

// MarshalJSON parse orders `source` field to string
func (b Source) MarshalJSON() ([]byte, error) {
	var s string
	switch b {
	case WEB:
		s = "WEB"
	case API:
		s = "API"
	case SYSTEM:
		s = "SYSTEM"
	default:
		return nil, errors.New("invalid order source")
	}
	return json.Marshal(s)
}

// UnmarshalJSON parse bytes to orders `source` field
func (b *Source) UnmarshalJSON(s []byte) error {
	var str string
	if err := json.Unmarshal(s, &str); err != nil {
		return err
	}
	switch str {
	case "WEB":
		*b = WEB
	case "API":
		*b = API
	case "SYSTEM":
		*b = SYSTEM
	default:
		return errors.New("invalid order source")
	}
	return nil
}

// Operator represents stop order operator
type Operator string

const (
	// LTE represents less than or equal
	LTE Operator = "LTE"
	// GTE represents great than or equal
	GTE Operator = "GTE"
)

// Order is orders table struct
type Order struct {
	ID           uint64          `gorm:"primary_key" json:"id"`
	CustomerID   uint64          `json:"customer_id"`
	MarketUUID   string          `gorm:"column:market_uuid" json:"asset_pair_uuid"`
	Price        decimal.Decimal `sql:"type:decimal(32,16);" json:"price"`
	Amount       decimal.Decimal `sql:"type:decimal(32,16);" json:"amount"`
	FilledAmount decimal.Decimal `sql:"type:decimal(32,16);" json:"filled_amount"`
	Side         Side            `json:"side"`
	State        State           `json:"state"`
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
	Type         Type            `json:"type"`
	BU           BusinessUnit    `json:"bu"`
	Source       Source          `json:"source"`
	PostOnly     bool            `json:"post_only"`
}

func NewOrder(id, customerID uint64, price, amount string, now time.Time, marketUUID string) *Order {
	p, _ := decimal.NewFromString(price)
	return &Order{ID: id, CustomerID: customerID, MarketUUID: marketUUID, Price: p, Amount: p, Side: Side(id % 2), State: StatePending, Type: TypeLimit, IOC: false, FilledAmount: decimal.Zero, AvgDealPrice: decimal.Zero, InsertedAt: now, UpdatedAt: now}
}
