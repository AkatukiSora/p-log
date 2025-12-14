package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Post holds the schema definition for the Post entity.
type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		field.String("content").
			NotEmpty().
			MaxLen(1000),
		field.Time("created_at").
			Default(time.Now).Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		// Post -> User (多対1)
		edge.From("user", User.Type).
			Ref("posts").
			Unique().
			Required(),
		// Post -> Goal (多対1、任意)
		edge.From("goal", Goal.Type).
			Ref("posts").
			Unique().
			Required(),
		// Post -> Image (1対多)
		// いったんバケットとの整合は考慮しない
		edge.To("images", Image.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// Post -> Reaction (1対多)
		edge.To("reactions", Reaction.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}

// Indexes of the Post.
func (Post) Indexes() []ent.Index {
	return []ent.Index{
		// ユーザー別投稿一覧取得用
		index.Edges("user"),
		// ゴール別投稿フィルタ用
		index.Edges("goal"),
		// タイムライン降順取得用
		index.Fields("created_at"),
	}
}
