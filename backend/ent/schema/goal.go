package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Goal holds the schema definition for the Goal entity.
type Goal struct {
	ent.Schema
}

// Fields of the Goal.
func (Goal) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		field.String("title").
			NotEmpty(),
		field.Time("deadline").
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Goal.
func (Goal) Edges() []ent.Edge {
	return []ent.Edge{
		// Goal -> User (多対1)
		edge.From("user", User.Type).
			Ref("goals").
			Unique().
			Required(),
		// Goal -> Post (1対多)
		edge.To("posts", Post.Type),
	}
}

// Indexes of the Goal.
func (Goal) Indexes() []ent.Index {
	return []ent.Index{
		// ユーザー別目標一覧取得用
		index.Edges("user"),
	}
}
