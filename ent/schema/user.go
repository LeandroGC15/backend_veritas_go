package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("email").
			Unique().
			NotEmpty().
			Comment("Email del usuario"),
		field.String("password").
			Sensitive().
			NotEmpty().
			Comment("Password hasheado"),
		field.String("name").
			NotEmpty().
			Comment("Nombre del usuario"),
		field.String("role").
			Default("user").
			Comment("Rol del usuario (admin, manager, user)"),
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

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("tenant_id"),
	}
}
