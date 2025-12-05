package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// InvoiceItem holds the schema definition for the InvoiceItem entity.
type InvoiceItem struct {
	ent.Schema
}

// Fields of the InvoiceItem.
func (InvoiceItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("invoice_id").
			Comment("ID de la factura"),
		field.Int("product_id").
			Comment("ID del producto"),
		field.Int("quantity").
			Min(1).
			Comment("Cantidad vendida"),
		field.Float("unit_price").
			Min(0).
			Comment("Precio unitario al momento de la venta"),
		field.Float("subtotal").
			Min(0).
			Comment("Subtotal (quantity * unit_price)"),
	}
}

// Edges of the InvoiceItem.
func (InvoiceItem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("invoice", Invoice.Type).
			Ref("items").
			Field("invoice_id").
			Required().
			Unique(),
	}
}

// Indexes of the InvoiceItem.
func (InvoiceItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("invoice_id"),
		index.Fields("product_id"),
	}
}

