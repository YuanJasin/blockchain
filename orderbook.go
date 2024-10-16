package main

import (
	"fmt"
	"sort"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}
type Orders []*Order

func (o Orders) Len() int               { return len(o) }
func (o Orders) Swap(i int, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i int, j int) bool { return o[i].Timestamp < o[j].Timestamp }

// NewOrder 创建一个新的订单
func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (o *Order) String() string {
	return fmt.Sprintf("Order{Size: %.2f}", o.Size)
}

type Limit struct {
	// 价格
	Price float64
	// 订单
	Orders Orders
	// 总交易量
	TotalVolume float64
}

type Limits []*Limit

type ByBestAsk struct {
	Limits
}

func (a *ByBestAsk) Len() int               { return len(a.Limits) }
func (a *ByBestAsk) Swap(i int, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a *ByBestAsk) Less(i int, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct {
	Limits
}

func (b *ByBestBid) Len() int               { return len(b.Limits) }
func (b *ByBestBid) Swap(i int, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b *ByBestBid) Less(i int, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

// NewLimit 创建一个限价
func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (l *Limit) String() string {
	return fmt.Sprintf("Limit{Price: %.2f | volume: %.2f}", l.Price, l.TotalVolume)
}

// AddOrder 添加账单
func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.TotalVolume += o.Size
}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
		}
	}
	o.Limit = nil
	l.TotalVolume -= o.Size

	sort.Sort(l.Orders)
}

// OrderBook 账本
type OrderBook struct {
	// 卖出
	Asks []*Limit
	// 买入
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceOrder(price float64, o *Order) []Match {
	// p1: 尝试匹配订单
	// 匹配逻辑
	// p2: 将订单的其余部分添加到书籍中
	if o.Size > 0.0 {
		ob.Add(price, o)
	}
	return []Match{}
}

func (ob *OrderBook) Add(price float64, o *Order) {
	var limit *Limit

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.Bids = append(ob.Bids, limit)
			ob.BidLimits[price] = limit
		} else {
			ob.Asks = append(ob.Asks, limit)
			ob.AskLimits[price] = limit
		}
	}
	limit.AddOrder(o)
}
