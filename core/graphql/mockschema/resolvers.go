package mockschema

import (
	"context"
	"time"

	"github.com/rjansen/fivecolors/core/util"
)

type Resolver struct{}

func NewResolver() *Resolver {
	return new(Resolver)
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Me(ctx context.Context) (MeResponse, error) {
	return MeResponse{
		Tid: util.NewUUID(),
		User: User{
			ID:   "fivecolorsd",
			Name: "Fivecolors D",
		},
	}, nil
}
func (r *queryResolver) MockEntity(ctx context.Context) (MockEntityResponse, error) {
	return MockEntityResponse{
		Tid: util.NewUUID(),
		Entity: MockEntity{
			ID:       util.NewUUID(),
			String:   "string field",
			Integer:  999,
			Float:    999.99,
			DateTime: time.Now().UTC(),
			Boolean:  false,
			Object: Object{
				"key_string":   "string value",
				"key_integer":  int64(999),
				"key_float":    float64(999.99),
				"key_boolean":  true,
				"key_datetime": time.Now().UTC(),
			},
		},
	}, nil
}
