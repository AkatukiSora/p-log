package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Reaction holds the schema definition for the Reaction entity.
type Reaction struct {
	ent.Schema
}

// Fields of the Reaction.
func (Reaction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		field.Time("created_at").
			Default(time.Now).Immutable(),
	}
}

// Edges of the Reaction.
func (Reaction) Edges() []ent.Edge {
	return []ent.Edge{
		// Reaction -> User (多対1)
		edge.From("user", User.Type).
			Ref("reactions").
			Unique().
			Required(),
		// Reaction -> Post (多対1)
		edge.From("post", Post.Type).
			Ref("reactions").
			Unique().
			Required(),
	}
}

// Indexes of the Reaction.
func (Reaction) Indexes() []ent.Index {
	return []ent.Index{
		// 1ユーザー1投稿につき1リアクションの複合ユニーク制約
		index.Edges("user", "post").
			Unique(),
	}
}
