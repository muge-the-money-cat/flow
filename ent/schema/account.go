package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Account holds the schema definition for the Account entity.
type Account struct {
	ent.Schema
}

// Fields of the Account.
func (Account) Fields() (fields []ent.Field) {
	fields = append(fields,
		field.String("name").
			Unique(),
	)

	return
}

// Edges of the Account.
func (Account) Edges() (edges []ent.Edge) {
	edges = append(edges,
		edge.To("subtotal", Subtotal.Type).
			Unique().
			Required(),
	)

	return
}
