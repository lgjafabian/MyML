package myml

import (
    mymlDomain "github.com/lgjafabian/MyML/src/api/domain/myml"
    addressesDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/addresses"
	paymentsDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/payments"
	itemsDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/items"
	ordersDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/orders"
	"github.com/emikohmann/itacademy2018-myml/src/api/util/apierrors"
	"github.com/mercadolibre/go-meli-toolkit/restful/rest"
	"fmt"
	"encoding/json"
	"sync"
)



func GetInformation(orderID int64) (*mymlDomain.Sumary, *apierrors.ApiError) {
	
	var (
		cValues [3] interface {}
		group sync.WaitGroup
		order 	ordersDomain.Order
		orderAux 	mymlDomain.OrderSumary
		addressAux mymlDomain.Address
		paymentAux mymlDomain.Payment
		itemAux    mymlDomain.Item
	)
	
	GetOrder(&order, &orderAux, orderID)
	
	c := make(chan interface{})
	group.Add(3)

	go GetItem(order, &group, c)
	go GetPayment(order, &group, c)
	go GetAddress(order, &group, c)
	
	cValues[0], cValues[1], cValues[2] = <-c, <-c, <-c
	group.Wait()

	GetChannelInfo(&addressAux, &paymentAux, &itemAux, cValues)

	return &mymlDomain.Sumary{
        Title:          fmt.Sprintf("Compra #%d", orderID),
		Order: 			orderAux,
		Item:			itemAux,
		Payment:		paymentAux,
		Address:		addressAux,
    }, nil
}

func GetChannelInfo(addressAux *mymlDomain.Address, paymentAux *mymlDomain.Payment,itemAux   * mymlDomain.Item, cValues [3] interface {}) {
	for i := 0; i < 3; i++ {
		switch v := cValues[i].(type) {
			case mymlDomain.Item:
				*itemAux = v
			case mymlDomain.Address:
				*addressAux = v
			case mymlDomain.Payment:
				*paymentAux = v
			default:
		}
	}
}

func GetOrder(order *ordersDomain.Order, orderAux *mymlDomain.OrderSumary, orderID int64)() {
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/orders/%d", orderID))
	json.Unmarshal(resp.Bytes(), &order)
	
	orderAux.Name = fmt.Sprintf("Pago de %d producto", len(order.Items))
	orderAux.TotalAmount = order.TotalAmount
}

func GetItem(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{})() {
	defer group.Done()

	var item itemsDomain.Item
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/items/%s", order.Items[0]))
	json.Unmarshal(resp.Bytes(), &item)
	
	c <- mymlDomain.Item{
		Name: item.Title,
		Price: item.Price,
		Quantity: item.Quantity,
		ImgUrl: item.Pictures[0],
	}	
}

func GetPayment(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{})() {
	defer group.Done()

	var payment paymentsDomain.Payment
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/payments/%s", order.Payments[0]))
	json.Unmarshal(resp.Bytes(), &payment)

	c <- mymlDomain.Payment{
		Method: "Dinero en cuenta de mercado Pago",
    	Amount: payment.TransactionAmount,
	}
}

func GetAddress(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{})() {
	defer group.Done()

	var address addressesDomain.Address
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/addresses/%s", order.Address))
	json.Unmarshal(resp.Bytes(), &address)
	
	c <- mymlDomain.Address {
		StreetAndNumber: fmt.Sprintf("%s %d", address.StreetName, address.StreetNumber),
    	City: address.City,
	}
	
}

