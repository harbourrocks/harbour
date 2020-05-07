package handler

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	traits2 "github.com/harbourrocks/harbour/pkg/harbourgateway/traits"
	"github.com/harbourrocks/harbour/pkg/httphandler/traits"
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}

// GraphQLModel is specific for one handler
type GraphQLModel struct {
	traits.HttpModel
	traits.IdTokenModel
	traits2.GraphQLModel
}

func (h GraphQLModel) Handle() (err error) {
	r := h.GetRequest()
	w := h.GetResponse()
	s := h.GetSchema()

	result := executeQuery(r.URL.Query().Get("query"), s)
	return json.NewEncoder(w).Encode(result)
}
