package floctory

import (
	"io"
)

type SendParams struct {
	HttpCode    int
	Path        string
	HttpMethod  string
	Date        string
	Token       string
	Body        io.Reader
	QueryParams map[string]string
	Response    interface{}
}

type Request struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	From    int64 `json:"from"`
	To      int64 `json:"to"`
}

// ExchangeLeadsResponseData represents a single lead's data
type ExchangeLeadsResponseData struct {
	Email                  string `json:"email"`                     // Email user left when accepting your offer
	CreatedAt              int64  `json:"created_at"`                // Unix timestamp
	FirstName              string `json:"first_name"`                // First name
	FullName               string `json:"full_name"`                 // Full name
	LastExchangeAcceptDate int64  `json:"last_exchange_accept_date"` // Last time this user accepted your offer; unix timestamp
}

// ExchangeLeadsResponse represents the response containing leads and pagination info
type ExchangeLeadsResponse struct {
	HasNextData bool                        `json:"has_next_data"` // Has data
	NextPage    string                      `json:"next_page"`     // Next page
	Data        []ExchangeLeadsResponseData `json:"data"`          // Array of lead data
}

// PhoneLeadsResponseData represents a single phone lead's data
type PhoneLeadsResponseData struct {
	Email     string `json:"email"`      // Email user left when accepting your offer
	Name      string `json:"name"`       // Name user left when accepting your offer
	Phone     string `json:"phone"`      // Phone user left when accepting your offer
	CreatedAt int64  `json:"created_at"` // The time customer accepted an offer; unix timestamp
}

// PhoneLeadsResponse represents the response containing phone leads and pagination info
type PhoneLeadsResponse struct {
	HasNextData bool                     `json:"has_next_data"` // Has data
	NextPage    string                   `json:"next_page"`     // Next page
	Data        []PhoneLeadsResponseData `json:"data"`          // Array of phone lead data
}
