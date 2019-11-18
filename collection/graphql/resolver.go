package graphql

import (
	"github.com/rjansen/fivecolors/collection"
)

type resolver struct {
	reader collection.Reader
	writer collection.Writer
}

func NewResolver(reader collection.Reader, writer collection.Writer) ResolverRoot {
	return &resolver{reader: reader, writer: writer}
}

func (r *resolver) Mutation() MutationResolver {
	return r.writer
}

func (r *resolver) Query() QueryResolver {
	return r.reader
}
