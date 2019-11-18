package firestore

import (
	"context"
	"fmt"
	"time"

	"github.com/rjansen/fivecolors/collection"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel/firestore"
)

var (
	timeNow = time.Now
)

type writer struct {
	logger l.Logger
	client firestore.Client
}

func NewWriter(logger l.Logger, client firestore.Client) collection.Writer {
	return &writer{
		logger: logger,
		client: client,
	}
}

func (w *writer) UpsertSet(ctx context.Context, set collection.SetInput) (*collection.Set, error) {
	w.logger.Debug(ctx, "resolve.upsertset.firestore", l.NewValue("arguments", set))
	setRef := fmt.Sprintf(collectionFmt, fmt.Sprintf("sets/%s", set.ID))
	type setDoc struct {
		collection.SetInput
		UpdatedAt time.Time `firestore:"serverTimestamp"`
	}
	doc := setDoc{SetInput: set}
	err := w.client.Doc(setRef).Set(ctx, &doc)
	if err != nil {
		return nil, err
	}

	result := collection.Set{
		ID:        doc.ID,
		Name:      doc.Name,
		Alias:     doc.Alias,
		Asset:     doc.Asset,
		UpdatedAt: doc.UpdatedAt,
	}
	w.logger.Debug(ctx, "resolve.upsertset.firestore.set", l.NewValue("set", result))
	return &result, err
}

func (w *writer) UpsertCards(ctx context.Context, cards []collection.CardInput) (*collection.UpsertCards, error) {
	w.logger.Debug(ctx, "resolve.upsertcard.firestore", l.NewValue("arguments", len(cards)))

	batch := w.client.Batch()
	type cardDoc struct {
		collection.CardInput
		UpdatedAt time.Time `firestore:"serverTimestamp"`
	}
	docs := 0
	for _, card := range cards {
		cardRef := w.client.Doc(fmt.Sprintf(collectionFmt, fmt.Sprintf("cards/%s", card.ID)))
		doc := cardDoc{CardInput: card}
		batch.Set(cardRef, &doc)
		docs++
	}
	err := batch.Commit(ctx)
	if err != nil {
		return nil, err
	}
	// TODO: Read commit time from batch write result
	committedAt := timeNow()

	result := collection.UpsertCards{AffectedRecords: docs, CommittedAt: committedAt}
	w.logger.Debug(ctx, "resolve.upsertcard.firestore.set", l.NewValue("result", result))
	return &result, err
}
