package client

import (
	"context"
	"net/http"

	"github.com/rjansen/fivecolors/collection"
)

type Reader struct {
	*graphqlClient
}

func NewReader(client *http.Client) *Reader {
	graphqlClient := newGraphQLClient(method, url, client)
	return &Reader{graphqlClient: graphqlClient}
}

func (c *Reader) Set(ctx context.Context, id string) (*collection.Set, error) {
	request := Request{
		Query: `
			query Set($id: ID!) {
				set(id: $id) {
					id
					name
					alias
				}
			}
		`,
		Variables: Object{"id": id},
	}
	type result struct {
		Response
		Data struct {
			Set collection.Set `json:"set"`
		} `json:"data"`
	}
	var response result

	err := c.Do(ctx, request, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Set, nil
}

func (c *Reader) Card(ctx context.Context, id string) (*collection.Card, error) {
	request := Request{
		Query: `
			query Card($id: ID!) {
				card(id: $id) {
					id
					name
					types
				}
			}
		`,
		Variables: Object{"id": id},
	}
	type result struct {
		Response
		Data struct {
			Card collection.Card `json:"card"`
		} `json:"data"`
	}
	var response result

	err := c.Do(ctx, request, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Card, nil
}

func (c *Reader) SetBy(ctx context.Context, filter collection.SetFilter) ([]collection.Set, error) {
	panic("not implemented")
}

func (c *Reader) CardBy(ctx context.Context, filter collection.CardFilter) ([]collection.Card, error) {
	panic("not implemented")
}
