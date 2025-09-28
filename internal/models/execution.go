package models

type Side string
type OrderStatus string
type PositionStatus string
type OrderEventType string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"

	Pending   OrderStatus = "PENDING"
	Cancelled OrderStatus = "CANCELLED"
	Executed  OrderStatus = "EXECUTED"

	Open   PositionStatus = "OPEN"
	Closed PositionStatus = "CLOSED"

	Create OrderEventType = "CREATE"
	Update OrderEventType = "UPDATE"
)

type Order struct {
	BaseModel
	Symbol string
	Side   Side
	Qty    int
	Entry  float64
	Stop   float64
	Target float64
	Reason string // e.g., "ORB_BREAKOUT_UP"
	Status OrderStatus
}

type Fill struct {
	Symbol string
	Side   Side
	Qty    int
	Price  float64
	Ref    string // order id or paper ref
}

type Position struct {
	BaseModel
	Symbol string
	Side   Side
	Qty    int
	Entry  float64
	Stop   float64
	Target float64
	Exit   float64
	Status PositionStatus
}

type OrderEvent struct {
	Type  OrderEventType
	Order *Order
}
