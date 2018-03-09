package myml

import (
	"sync"
	"testing"

	ordersDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/orders"
	mymlDomain "github.com/lgjafabian/MyML/src/api/domain/myml"
)

func TestGetInformation(t *testing.T) {

	if _, err := GetInformation(-1); err == nil {
		t.Error(err.Error)
	}

	if _, err := GetInformation(1); err != nil {
		t.Error(err.Error)
	}

}

func TestGetOrder(t *testing.T) {
	var (
		order    ordersDomain.Order
		orderAux mymlDomain.OrderSumary
	)

	if err := GetOrder(&order, &orderAux, 1); err != nil {
		t.Error(err)
	}

	if err := GetOrder(&order, &orderAux, -1); err == nil {
		t.Error(err)
	}

	if err := GetOrder(&order, &orderAux, 1000000000000000000); err != nil {
		t.Error(err)
	}

	if err := GetOrder(nil, nil, 1); err == nil {
		t.Error(err)
	}
}

func TestGetItem(t *testing.T) {
	var (
		cValue interface{}
		group  sync.WaitGroup
		order  ordersDomain.Order
	)

	order.Items = []string{"1"}
	c := make(chan interface{})

	group.Add(1)
	go GetItem(order, &group, c)
	cValue = <-c
	group.Wait()

	if cValue == nil {
		t.Error("Returned null value")
	}
}

func TestGetPayment(t *testing.T) {
	var (
		cValue interface{}
		group  sync.WaitGroup
		order  ordersDomain.Order
	)

	order.Payments = []string{"1"}
	c := make(chan interface{})

	group.Add(1)
	go GetPayment(order, &group, c)
	cValue = <-c
	group.Wait()

	if cValue == nil {
		t.Error("Returned null value")
	}
}

func TestGetAddress(t *testing.T) {
	var (
		cValue interface{}
		group  sync.WaitGroup
		order  ordersDomain.Order
	)

	order.Address = "1"
	c := make(chan interface{})

	group.Add(1)
	go GetAddress(order, &group, c)
	cValue = <-c
	group.Wait()

	if cValue == nil {
		t.Error("Returned null value")
	}
}

func BenchmarkGetOrder(b *testing.B) {
	var (
		order    ordersDomain.Order
		orderAux mymlDomain.OrderSumary
	)

	GetOrder(&order, &orderAux, 1)
}

func BenchmarkGetInformation(b *testing.B) {
	// BenchmarkGetInformation-4   	2000000000	         0.00 ns/op	       0 B/op	       0 allocs/op
	GetInformation(1)
}

func BenchmarkGetInformationFor10(b *testing.B) {
	// BenchmarkGetInformationFor1000-4   	       1	1391976065 ns/op	45765672 B/op	  592007 allocs/op
	for n := 0; n < 10; n++ {
		GetInformation(1)
	}
}

func BenchmarkGetInformationFor1000(b *testing.B) {
	// BenchmarkGetInformationFor1000-4   	       1	1391976065 ns/op	45765672 B/op	  592007 allocs/op
	for n := 0; n < 1000; n++ {
		GetInformation(1)
	}
}

func BenchmarkGetInformationFor10000(b *testing.B) {
	// BenchmarkGetInformationFor10000-4   	       1	12375727561 ns/op	457815336 B/op	 5928422 allocs/op
	for n := 0; n < 10000; n++ {
		GetInformation(1)
	}
}
