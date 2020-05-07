package traits

import (
	"github.com/graphql-go/graphql"
)

// GraphQLTrait returns a schema
type GraphQLTrait interface {
	GetSchema() graphql.Schema
	SetSchema(graphql.Schema)
}

// RequestModel holds the request
type GraphQLModel struct {
	schema graphql.Schema
}

func (m GraphQLModel) GetSchema() graphql.Schema {
	return m.schema
}

func (m *GraphQLModel) SetSchema(s graphql.Schema) {
	m.schema = s
}

func AddGraphQL(trait GraphQLTrait, s graphql.Schema) {
	trait.SetSchema(s)
}
