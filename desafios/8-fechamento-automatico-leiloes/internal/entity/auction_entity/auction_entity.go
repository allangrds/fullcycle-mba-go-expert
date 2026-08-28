package auction_entity

import (
	"context"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/internal/internal_error"
	"github.com/google/uuid"
)

func CreateAuction(
	productName, category, description string,
	condition ProductCondition) (*Auction, *internal_error.InternalError) {
	auction := &Auction{
		Id:          uuid.New().String(),
		ProductName: productName,
		Category:    category,
		Description: description,
		Condition:   condition,
		Status:      Active,
		Timestamp:   time.Now(),
	}

	if err := auction.Validate(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (au *Auction) Validate() *internal_error.InternalError {
	if len(au.ProductName) <= 1 ||
		len(au.Category) <= 2 ||
		len(au.Description) <= 10 && (au.Condition != New &&
			au.Condition != Refurbished &&
			au.Condition != Used) {
		return internal_error.NewBadRequestError("invalid auction object")
	}

	return nil
}

type Auction struct {
	Id          string
	ProductName string
	Category    string
	Description string
	Condition   ProductCondition
	Status      AuctionStatus
	Timestamp   time.Time
}

type ProductCondition int
type AuctionStatus int

// Active é o leilão em andamento. Completed é o estado "Closed" (fechado)
// pedido pelo desafio: é para esse status que a rotina de fechamento
// automático atualiza o leilão quando o AUCTION_DURATION expira.
const (
	Active AuctionStatus = iota
	Completed
)

const (
	New ProductCondition = iota + 1
	Used
	Refurbished
)

type AuctionRepositoryInterface interface {
	CreateAuction(
		ctx context.Context,
		auctionEntity *Auction) *internal_error.InternalError

	FindAuctions(
		ctx context.Context,
		status AuctionStatus,
		category, productName string) ([]Auction, *internal_error.InternalError)

	FindAuctionById(
		ctx context.Context, id string) (*Auction, *internal_error.InternalError)
}
