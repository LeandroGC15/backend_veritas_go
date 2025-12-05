package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Invoice holds the schema definition for the Invoice entity.
type Invoice struct {
	ent.Schema
}

// Fields of the Invoice.
func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.Float("total").
			Min(0).
			Comment("Total de la factura"),
		field.String("status").
			Default("pending").
			Comment("Estado de la factura (pending, paid, cancelled)"),
		field.Int("tenant_id").
			Comment("ID del tenant"),
		field.Int("user_id").
			Comment("ID del usuario que creó la factura"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Invoice.
func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("items", InvoiceItem.Type),
	}
}

// Indexes of the Invoice.
func (Invoice) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("user_id"),
		index.Fields("status"),
	}
}
