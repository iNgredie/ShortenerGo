package shortening

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"shortener/internal/model"
	"time"
)

type mgo struct {
	db *mongo.Database
}

func NewMongoDB(client *mongo.Client) *mgo {
	return &mgo{db: client.Database("url-shortener")}
}

func (m *mgo) col() *mongo.Collection {
	return m.db.Collection("shortenings")
}

func (m *mgo) Put(ctx context.Context, shortening model.Shortening) (*model.Shortening, error) {
	const op = "shortening.mgo.Put"

	shortening.CreatedAt = time.Now().UTC()

	// 1. Check, that collection doesn't have document with the same identifier
	count, err := m.col().CountDocuments(ctx, bson.M{"_id": shortening.Identifier})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return nil, fmt.Errorf("%s: %w", op, model.ErrIdentifierExists)
	}
	return &shortening, nil
}

func (m *mgo) Get(ctx context.Context, shorteningID string) (*model.Shortening, error) {
	const op = "shortening.mgo.Get"

	var shortening mgoShortening
	if err := m.col().FindOne(ctx, bson.M{"_id": shorteningID}).Decode(&shortening); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%s: %w", op, model.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return modelShorteningFromMgo(shortening), nil
}

func (m *mgo) IncrementVisits(ctx context.Context, shorteningID string) error {
	const op = "shortening.mgo.IncrementVisits"

	var (
		filter = bson.M{"_id": shorteningID}
		update = bson.M{
			"$inc": bson.M{"visit": 1},
			"$set": bson.M{"update_at": time.Now().UTC()},
		}
	)
	_, err := m.col().UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type mgoShortening struct {
	Identifier  string    `json:"_id"`
	OriginalUrl string    `json:"originalUrl"`
	Visits      int64     `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func mgoShorteningFromModel(shortening model.Shortening) mgoShortening {
	return mgoShortening{
		Identifier:  shortening.Identifier,
		OriginalUrl: shortening.OriginalUrl,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}

func modelShorteningFromMgo(shortening mgoShortening) *model.Shortening {
	return &model.Shortening{
		Identifier:  shortening.Identifier,
		OriginalUrl: shortening.OriginalUrl,
		Visits:      shortening.Visits,
		CreatedAt:   shortening.CreatedAt,
		UpdatedAt:   shortening.UpdatedAt,
	}
}
