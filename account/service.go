package account

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
)

type Service struct {
	client *alpaca.Client
}

func NewAccountService(client *alpaca.Client) *Service {
	return &Service{client: client}
}

func (s *Service) GetAccount() (account *alpaca.Account, err error) {

	return s.client.GetAccount()
}
