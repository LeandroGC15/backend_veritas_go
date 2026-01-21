package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Supplier holds the schema definition for the Supplier entity.
type Supplier struct {
	ent.Schema
}

// Fields of the Supplier.
func (Supplier) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Comment("Nombre del proveedor"),
		field.String("email").
			Optional().
			Comment("Email de contacto"),
		field.String("phone").
			Optional().
			Comment("Teléfono de contacto"),
		field.Text("address").
			Optional().
			Comment("Dirección del proveedor"),
		field.String("ruc_nit").
			Optional().
			Comment("Número de identificación fiscal (RUC/NIT)"),
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

// Edges of the Supplier.
func (Supplier) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("purchase_invoices", PurchaseInvoice.Type).
			Ref("supplier"),
	}
}

// Indexes of the Supplier.
func (Supplier) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("tenant_id"),
		index.Fields("ruc_nit").Unique(),
	}
}
