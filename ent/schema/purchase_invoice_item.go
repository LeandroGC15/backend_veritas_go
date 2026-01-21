package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PurchaseInvoiceItem holds the schema definition for the PurchaseInvoiceItem entity.
type PurchaseInvoiceItem struct {
	ent.Schema
}

// Fields of the PurchaseInvoiceItem.
func (PurchaseInvoiceItem) Fields() []ent.Field {
	return []ent.Field{
		field.Int("purchase_invoice_id").
			Comment("ID de la factura de compra"),
		field.Int("product_id").
			Comment("ID del producto comprado"),
		field.Int("quantity").
			Min(1).
			Comment("Cantidad comprada"),
		field.Float("unit_cost").
			Min(0).
			Comment("Costo unitario al momento de la compra"),
		field.Float("subtotal").
			Min(0).
			Comment("Subtotal (quantity × unit_cost)"),
	}
}

// Edges of the PurchaseInvoiceItem.
func (PurchaseInvoiceItem) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the PurchaseInvoiceItem.
func (PurchaseInvoiceItem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("purchase_invoice_id"),
		index.Fields("product_id"),
	}
}