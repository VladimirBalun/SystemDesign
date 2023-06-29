package main

type Storage struct {
}

func NewStorage() Storage {
	return Storage{}
}

func CreateOrder(order Order) (int, error) {
	return 0, nil
}

func RemoveOrder(orderID int) error {
	return nil
}

func SubmitOrder(orderID int) error {
	return nil
}
