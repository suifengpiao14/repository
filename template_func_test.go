package repository

import (
	"fmt"
	"testing"
)

func TestJsonSchema2Path(t *testing.T) {
	jsonschema := `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","properties":{"items":{"type":"array","items":{"type":"object","properties":{"sourceId":{"type":"string","src":"PaginateOut.#.api_id"},"sourceType":{"type":"string","src":"PaginateOut.#.name"},"url":{"type":"string","src":"PaginateOut.#.url"},"createTime":{"type":"string","src":"PaginateOut.#.created_at"},"updateTime":{"type":"string","src":"PaginateOut.#.updated_at"}},"required":["sourceId","sourceType","url","createTime","updateTime"]}},"pagination":{"type":"object","properties":{"total":{"type":"integer","src":"PaginateTotalOut"},"pageIndex":{"type":"string","src":"PageIndex"},"pageSize":{"type":"string","src":"PageSize"}},"required":["total","pageIndex","pageSize"]}},"required":["items","pagination"]}`

	paths, err := JsonSchema2Path(jsonschema)
	if err != nil {
		panic(err)
	}

	for _, path := range paths {
		fmt.Println(*path)
	}

}
