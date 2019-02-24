// +build !firestore

package model

type Resolver struct {
	queryResolver QueryResolver
}

func NewResolver(queryResolver QueryResolver) *Resolver {
	return &Resolver{
		queryResolver: queryResolver,
	}
}

func (r *Resolver) Query() QueryResolver {
	return r.queryResolver
}
