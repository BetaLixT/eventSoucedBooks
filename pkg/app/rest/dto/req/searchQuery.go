package req

import (
	"eventSourcedBooks/pkg/domain/base"
)

type SearchQuery base.SearchQuery

func MapDtoToSearchQuery(dto *SearchQuery) base.SearchQuery {
	return base.SearchQuery{
		Key:       dto.Key,
		Operation: dto.Operation,
		Value:     dto.Value,
		ValueType: dto.ValueType,
	}
}

func MapDtoToSearchQueryList(dtos []SearchQuery) (sqs []base.SearchQuery) {
	sqs = make([]base.SearchQuery, len(dtos))
	for idx, dto := range dtos {
		sqs[idx] = MapDtoToSearchQuery(&dto)
	}
	return
}
