package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		field.String("name").
			NotEmpty(),
		field.String("email").
			NotEmpty().
			Unique(),
		field.Time("birthday").
			Optional().
			Nillable(),
		field.String("hometown").
			Optional().
			Nillable(),
		field.String("bio").
			Optional().
			Nillable(),
		field.UUID("profile_picture_id", uuid.UUID{}).
			Optional().
			Nillable(),
		field.Time("created_at").
			Default(time.Now).Immutable(),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// User -> Genre (多対多)
		edge.To("genres", Genre.Type),
		// User -> Goal (1対多)
		edge.To("goals", Goal.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// User -> Post (1対多)
		edge.To("posts", Post.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// User -> Reaction (1対多)
		edge.To("reactions", Reaction.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// User -> Image (アップロードした画像、1対多)
		edge.To("uploaded_images", Image.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// User -> RefreshToken (1対多)
		edge.To("refresh_tokens", RefreshToken.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
		// フォロー関係 (一方通行)
		edge.To("following", User.Type).
			From("followers"),
	}
}
