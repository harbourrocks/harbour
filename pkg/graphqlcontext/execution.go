package graphqlcontext

import (
	"context"
	"github.com/graphql-go/graphql"
	"github.com/harbourrocks/harbour/pkg/logconfig"
)

func ExecuteQuery(ctx context.Context, schema graphql.Schema, query string, variables map[string]interface{}) *graphql.Result {
	log := logconfig.GetLogCtx(ctx)

	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
		Context:        ctx,
	})

	if len(result.Errors) > 0 {
		log.Errorf("wrong result, unexpected errors: %v", result.Errors)
	}

	return result
}
