package myml

type Sumary struct{
    Title           string          `json:"title"`
    Order           OrderSumary           `json:"order_sumary"`
    Item            Item            `json:"item"`
    Payment         Payment         `json:"payment"`
    Address         Address         `json:"address"`

}

type OrderSumary struct {
    Name      string   `json:"name"`
    TotalAmount float32  `json:"total_amount"`
}

type Item struct {
    Name        string      `json:"name"`
    Price       float32     `json:"price"`
    Quantity    int32       `json:"quantity"`
    ImgUrl      string      `json:"img_url"`
}

type Payment struct {
    Method      string      `json:"method"`
    Amount      float32      `json:"amount"`
}

type Address struct {
    StreetAndNumber     string      `json:"street_and_number"`
    City                string      `json:"city"`
}

