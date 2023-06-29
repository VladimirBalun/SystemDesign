package main

type Order struct {
	username string
	goods    []string
}

func NewOrder(u string, g []string) Order {
	return Order{
		username: u,
		goods:    g,
	}
}

func (o *Order) Username() string {
	return o.username
}

func (o *Order) Goods() []string {
	return o.goods
}
