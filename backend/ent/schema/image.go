package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Image holds the schema definition for the Image entity.
type Image struct {
	ent.Schema
}

// Fields of the Image.
func (Image) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		// GCS内のオブジェクトパス (例: "images/{uuid}.jpg")
		field.String("object_name").
			NotEmpty(),
		// MIMEタイプ (例: "image/jpeg")
		field.String("content_type").
			NotEmpty(),
		field.Time("created_at").
			Default(time.Now).Immutable(),
	}
}

// Edges of the Image.
func (Image) Edges() []ent.Edge {
	return []ent.Edge{
		// Image -> Post (多対1、任意: 投稿に紐付く前はNULL)
		edge.From("post", Post.Type).
			Ref("images").
			Unique(),
		// Image -> User (多対1: アップロードしたユーザー)
		edge.From("uploaded_by", User.Type).
			Ref("uploaded_images").
			Unique().
			Required(),
	}
}

// Indexes of the Image.
func (Image) Indexes() []ent.Index {
	return []ent.Index{
		// 投稿別画像取得用
		index.Edges("post"),
		// ユーザー別画像取得用
		index.Edges("uploaded_by"),
	}
}
