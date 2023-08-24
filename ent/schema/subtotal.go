package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subtotal holds the schema definition for the Subtotal entity.
type Subtotal struct {
	ent.Schema
}

// Fields of the Subtotal.
func (Subtotal) Fields() (fields []ent.Field) {
	fields = append(fields,
		field.String("name").
			Unique(),
	)

	return
}

// Edges of the Subtotal.
func (Subtotal) Edges() (edges []ent.Edge) {
	edges = append(edges,
		edge.To("children", Subtotal.Type).
			From("parent").
			Unique(),
	)

	return
}
