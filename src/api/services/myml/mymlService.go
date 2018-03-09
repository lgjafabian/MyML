package myml

import (
	"encoding/json"
	"fmt"
	"sync"

	addressesDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/addresses"
	itemsDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/items"
	ordersDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/orders"
	paymentsDomain "github.com/emikohmann/itacademy2018-myml/src/api/domain/payments"
	"github.com/emikohmann/itacademy2018-myml/src/api/util/apierrors"
	mymlDomain "github.com/lgjafabian/MyML/src/api/domain/myml"
	"github.com/mercadolibre/go-meli-toolkit/restful/rest"
)

// GetInformation Creates a JSON out of other API
func GetInformation(orderID int64) (*mymlDomain.Sumary, *apierrors.ApiError) {

	var (
		cValues    [3]interface{}
		group      sync.WaitGroup
		order      ordersDomain.Order
		orderAux   mymlDomain.OrderSumary
		addressAux mymlDomain.Address
		paymentAux mymlDomain.Payment
		itemAux    mymlDomain.Item
	)

	if err := GetOrder(&order, &orderAux, orderID); err != nil {
		return nil, err
	}

	c := make(chan interface{})
	group.Add(3)

	go GetItem(order, &group, c)
	go GetPayment(order, &group, c)
	go GetAddress(order, &group, c)

	cValues[0], cValues[1], cValues[2] = <-c, <-c, <-c
	group.Wait()

	GetChannelInfo(&addressAux, &paymentAux, &itemAux, cValues)

	return &mymlDomain.Sumary{
		Title:   fmt.Sprintf("Compra #%d", orderID),
		Order:   orderAux,
		Item:    itemAux,
		Payment: paymentAux,
		Address: addressAux,
	}, nil
}

// GetChannelInfo recovers data from Channel
func GetChannelInfo(addressAux *mymlDomain.Address, paymentAux *mymlDomain.Payment, itemAux *mymlDomain.Item, cValues [3]interface{}) {
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

// GetOrder Get and Format an order from 8081 API
func GetOrder(order *ordersDomain.Order, orderAux *mymlDomain.OrderSumary, orderID int64) *apierrors.ApiError {

	if orderID < 0 || order == nil || orderAux == nil {
		return &apierrors.ApiError{
			StatusCode: 400,
			Error:      "Wrong arguments",
		}
	}

	resp := rest.Get(fmt.Sprintf("http://localhost:8081/orders/%d", orderID))

	json.Unmarshal(resp.Bytes(), &order)

	orderAux.Name = fmt.Sprintf("Pago de %d producto", len(order.Items))
	orderAux.TotalAmount = order.TotalAmount
	return nil
}

// GetItem Format an Item from 8081 API
func GetItem(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{}) {
	defer group.Done()

	var item itemsDomain.Item
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/items/%s", order.Items[0]))
	json.Unmarshal(resp.Bytes(), &item)

	c <- mymlDomain.Item{
		Name:     item.Title,
		Price:    item.Price,
		Quantity: item.Quantity,
		ImgUrl:   item.Pictures[0],
	}
}

// GetPayment Format a Payment from 8081 API
func GetPayment(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{}) {
	defer group.Done()

	var payment paymentsDomain.Payment
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/payments/%s", order.Payments[0]))
	json.Unmarshal(resp.Bytes(), &payment)

	c <- mymlDomain.Payment{
		Method: "Dinero en cuenta de mercado Pago",
		Amount: payment.TransactionAmount,
	}
}

// GetAddress Format an Address from 8081 API
func GetAddress(order ordersDomain.Order, group *sync.WaitGroup, c chan interface{}) {
	defer group.Done()

	var address addressesDomain.Address
	resp := rest.Get(fmt.Sprintf("http://localhost:8081/addresses/%s", order.Address))
	json.Unmarshal(resp.Bytes(), &address)

	c <- mymlDomain.Address{
		StreetAndNumber: fmt.Sprintf("%s %d", address.StreetName, address.StreetNumber),
		City:            address.City,
	}

}
