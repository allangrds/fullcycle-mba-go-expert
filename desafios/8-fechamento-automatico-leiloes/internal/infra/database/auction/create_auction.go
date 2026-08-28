package auction

import (
	"context"
	"os"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/configuration/logger"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/internal/entity/auction_entity"
	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection      *mongo.Collection
	auctionDuration time.Duration
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection:      database.Collection("auctions"),
		auctionDuration: getAuctionDuration(),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go ar.scheduleAuctionClosure(auctionEntity.Id)

	return nil
}

// scheduleAuctionClosure roda em background assim que o leilão é criado e,
// depois de auctionDuration, fecha o leilão automaticamente (status
// Completed/"Closed"). Não recebe o ctx da request: usa seu próprio
// context.Background() para não ser cancelada quando a requisição HTTP que
// criou o leilão termina.
func (ar *AuctionRepository) scheduleAuctionClosure(auctionId string) {
	time.Sleep(ar.auctionDuration)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": auctionId}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	if _, err := ar.Collection.UpdateOne(ctx, filter, update); err != nil {
		logger.Error("Error trying to close auction automatically", err)
	}
}

func getAuctionDuration() time.Duration {
	duration := os.Getenv("AUCTION_DURATION")
	parsed, err := time.ParseDuration(duration)
	if err != nil {
		return 5 * time.Minute
	}

	return parsed
}
