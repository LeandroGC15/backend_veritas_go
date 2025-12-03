package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Product holds the schema definition for the Product entity.
type Product struct {
	ent.Schema
}

// Fields of the Product.
func (Product) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("Nombre del producto"),
		field.Text("description").
			Optional().
			Comment("Descripción del producto"),
		field.Float("price").
			Min(0).
			Comment("Precio del producto"),
		field.Int("stock").
			Default(0).
			Min(0).
			Comment("Cantidad en stock"),
		field.String("sku").
			Optional().
			Comment("SKU del producto"),
		field.Int("tenant_id").
			Comment("ID del tenant al que pertenece"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Product.
func (Product) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the Product.
func (Product) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("sku"),
	}
}
