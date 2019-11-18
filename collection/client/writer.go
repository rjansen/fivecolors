package client

import (
	"context"
	"net/http"

	"github.com/rjansen/fivecolors/collection"
)

type Writer struct {
	*graphqlClient
}

func NewWriter(client *http.Client) *Writer {
	graphqlClient := newGraphQLClient(method, url, client)
	return &Writer{graphqlClient: graphqlClient}
}

func (c *Writer) UpsertSet(ctx context.Context, set collection.SetInput) (*collection.Set, error) {
	request := Request{
		Query: `
			mutation UpsertSet($set: SetInput!) {
				upsertSet(set: $set) {
					id
					name
					alias
				}
			}
		`,
		Variables: Object{"set": set},
	}
	type result struct {
		Response
		Data struct {
			Set collection.Set `json:"upsertSet"`
		} `json:"data"`
	}
	var response result

	err := c.Do(ctx, request, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Set, nil
}

func (c *Writer) UpsertCards(ctx context.Context, cards []collection.CardInput) (*collection.UpsertCards, error) {
	request := Request{
		Query: `
			mutation UpsertCards($cards: [CardsInput!]!) {
				upsertCards(cards: $cards) {
					id
					name
					types
				}
			}
		`,
		Variables: Object{"cards": cards},
	}
	type result struct {
		Response
		Data struct {
			UpsertCards collection.UpsertCards `json:"upsertCards"`
		} `json:"data"`
	}
	var response result

	err := c.Do(ctx, request, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.UpsertCards, nil
}
