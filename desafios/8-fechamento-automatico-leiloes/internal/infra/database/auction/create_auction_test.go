package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/allangrds/fullcycle-mba-go-expert/desafios/8-fechamento-automatico-leiloes/internal/entity/auction_entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestCreateAuction_ClosesAutomaticallyAfterAuctionDuration comprova o
// requisito do desafio: um leilão criado com AUCTION_DURATION curto deve
// ser fechado (status Completed) sozinho, sem nenhuma chamada manual,
// assim que a duração expira. Precisa de um MongoDB real acessível via
// MONGODB_URL (suba com `make test`, que já sobe o serviço mongodb).
func TestCreateAuction_ClosesAutomaticallyAfterAuctionDuration(t *testing.T) {
	const auctionDuration = 2 * time.Second
	os.Setenv("AUCTION_DURATION", auctionDuration.String())

	ctx := context.Background()

	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		mongoURL = "mongodb://admin:admin@localhost:27017/auctions?authSource=admin"
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	require.NoError(t, err)
	defer client.Disconnect(ctx)
	require.NoError(t, client.Ping(ctx, nil))

	database := client.Database("auctions_test")
	defer database.Collection("auctions").Drop(ctx)

	repo := NewAuctionRepository(database)

	auction, createErr := auction_entity.CreateAuction(
		"Notebook Dell",
		"Eletrônicos",
		"Notebook seminovo em ótimo estado de conservação",
		auction_entity.Used)
	require.Nil(t, createErr)

	insertErr := repo.CreateAuction(ctx, auction)
	require.Nil(t, insertErr)

	justCreated, findErr := repo.FindAuctionById(ctx, auction.Id)
	require.Nil(t, findErr)
	assert.Equal(t, auction_entity.Active, justCreated.Status)

	time.Sleep(auctionDuration + 2*time.Second)

	closed, findErr := repo.FindAuctionById(ctx, auction.Id)
	require.Nil(t, findErr)
	assert.Equal(t, auction_entity.Completed, closed.Status)
}
