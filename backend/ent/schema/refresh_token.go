package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// RefreshToken holds the schema definition for the RefreshToken entity.
type RefreshToken struct {
	ent.Schema
}

// Fields of the RefreshToken.
func (RefreshToken) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable().Unique(),
		// トークンのSHA512ハッシュ値
		field.String("token_hash").
			NotEmpty().
			Unique(),
		// トークンの有効期限
		field.Time("expires_at"),
		// トークンが無効化されたかどうか
		field.Bool("revoked").
			Default(false),
		// トークン作成日時
		field.Time("created_at").
			Default(time.Now).Immutable(),
	}
}

// Edges of the RefreshToken.
func (RefreshToken) Edges() []ent.Edge {
	return []ent.Edge{
		// RefreshToken -> User (多対1、必須)
		edge.From("user", User.Type).
			Ref("refresh_tokens").
			Unique().
			Required(),
	}
}

// Indexes of the RefreshToken.
func (RefreshToken) Indexes() []ent.Index {
	return []ent.Index{
		// token_hashにユニークインデックス
		index.Fields("token_hash").
			Unique(),
		// revoked と expires_at の複合インデックス（有効なトークン検索用）
		index.Fields("revoked", "expires_at"),
		// user_id へのインデックス（Edgeで自動生成されるが明示的に定義）
		index.Edges("user"),
	}
}
