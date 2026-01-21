package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// PurchaseInvoice holds the schema definition for the PurchaseInvoice entity.
type PurchaseInvoice struct {
	ent.Schema
}

// Fields of the PurchaseInvoice.
func (PurchaseInvoice) Fields() []ent.Field {
	return []ent.Field{
		field.String("invoice_number").
			NotEmpty().
			Comment("Número de factura del proveedor"),
		field.Float("total").
			Min(0).
			Comment("Total de la factura"),
		field.String("status").
			Default("pending").
			Comment("Estado de pago (pending, partial, paid, cancelled)"),
		field.String("payment_method").
			Optional().
			Comment("Método de pago (cash, credit, partial)"),
		field.Time("due_date").
			Optional().
			Comment("Fecha de vencimiento"),
		field.Float("paid_amount").
			Default(0).
			Min(0).
			Comment("Monto pagado hasta el momento"),
		field.Int("supplier_id").
			Comment("ID del proveedor"),
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

// Edges of the PurchaseInvoice.
func (PurchaseInvoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("supplier", Supplier.Type).
			Field("supplier_id").
			Required().
			Unique(),
	}
}

// Indexes of the PurchaseInvoice.
func (PurchaseInvoice) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("supplier_id"),
		index.Fields("status"),
		index.Fields("due_date"),
		index.Fields("invoice_number").Unique(),
	}
}