package queries

import (
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
)

func (q *EsQuery) Build() *types.Query {
	queries := []types.Query{}
	filters := []types.Query{}

	if q.SearchText != nil && *q.SearchText != "" {
		queries = append(queries, types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query: *q.SearchText,
				Fields: []string{"title", "description"},
			},
		})
	}

	if q.MinPrice != nil || q.MaxPrice != nil{
		priceRange := types.NumberRangeQuery{} 
		var convert types.Float64
		
		if q.MinPrice != nil {
			convert = types.Float64(*q.MinPrice)
			priceRange.Gte = &convert 
		}

		if q.MaxPrice != nil {
			convert = types.Float64(*q.MaxPrice)
			priceRange.Lte = &convert
		}

		filters = append(filters, types.Query{
			Range: map[string]types.RangeQuery{
				"price": priceRange,
			},
		})
	}

	if q.AvailableQuantity != nil {
		quantity := types.Float64(*q.AvailableQuantity)
		filters = append(filters, types.Query{
			Range: map[string]types.RangeQuery{
				"available_quantity": types.NumberRangeQuery{
					Gte: &quantity,
				},
			},
		})
	}

	if q.Status != nil && *q.Status != "" {
		filters = append(filters, types.Query{
			Term: map[string]types.TermQuery{
				"status": {Value: *q.Status},
			},
		})
	}

	if q.Seller != nil {
		filters = append(filters, types.Query{
			Term: map[string]types.TermQuery{
				"seller": {Value: *q.Seller},
			},
		})
	}

	return &types.Query{
		Bool: &types.BoolQuery{
			Must: queries,
			Filter: filters,
		},
	}
}