package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Genre holds the schema definition for the Genre entity.
type Genre struct {
	ent.Schema
}

// Fields of the Genre.
func (Genre) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		field.String("name").
			NotEmpty().
			Unique(),
		field.Time("created_at").
			Default(time.Now).Immutable(),
	}
}

// Edges of the Genre.
func (Genre) Edges() []ent.Edge {
	return []ent.Edge{
		// Genre -> User (多対多、逆参照)
		edge.From("users", User.Type).
			Ref("genres"),
	}
}
