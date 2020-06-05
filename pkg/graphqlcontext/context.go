package graphqlcontext

import (
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"net/http"
)

func UseGraphQl(schema graphql.Schema) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var reqBody map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := reqBody["query"]
		variables := reqBody["variables"]

		if variables == nil {
			variables = make(map[string]interface{})
		}

		result := ExecuteQuery(ctx, schema, query.(string), variables.(map[string]interface{}))

		_ = httphelper.WriteResponse(r, w, result)
	}

	return fn
}
