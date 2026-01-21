package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// SupplierPayment holds the schema definition for the SupplierPayment entity.
type SupplierPayment struct {
	ent.Schema
}

// Fields of the SupplierPayment.
func (SupplierPayment) Fields() []ent.Field {
	return []ent.Field{
		field.Int("purchase_invoice_id").
			Comment("ID de la factura de compra"),
		field.Int("supplier_id").
			Comment("ID del proveedor"),
		field.Float("amount").
			Min(0).
			Comment("Monto del pago"),
		field.Time("payment_date").
			Default(time.Now).
			Comment("Fecha del pago"),
		field.String("payment_method").
			NotEmpty().
			Comment("Método de pago (cash, transfer, check, credit_card)"),
		field.String("reference").
			Optional().
			Comment("Referencia del pago (número de cheque, transferencia, etc.)"),
		field.Text("notes").
			Optional().
			Comment("Notas adicionales"),
		field.Int("tenant_id").
			Comment("ID del tenant"),
		field.Int("user_id").
			Comment("ID del usuario que registró el pago"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the SupplierPayment.
func (SupplierPayment) Edges() []ent.Edge {
	return nil
}

// Indexes of the SupplierPayment.
func (SupplierPayment) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("supplier_id"),
		index.Fields("purchase_invoice_id"),
		index.Fields("payment_date"),
	}
}