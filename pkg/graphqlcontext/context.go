package graphqlcontext

import (
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/httphelper"
	"net/http"
)

func UseGraphQl(schema graphql.Schema) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		query := r.URL.Query().Get("query")

		result := ExecuteQuery(ctx, schema, query)
		_ = httphelper.WriteResponse(r, w, result)
	}

	return fn
}
