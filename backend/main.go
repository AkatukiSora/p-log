package main

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean ../docs/api.yaml

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"backend/ent"
	"backend/ent/goal"
	"backend/ent/image"
	"backend/ent/user"

	"math/rand"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func main() {
	// 環境変数から接続情報を取得（デフォルト値を設定）
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "p-log")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbName, dbPassword)

	log.Printf("Connecting to database at %s:%s", dbHost, dbPort)

	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	// ランダムなテストデータを生成・挿入
	if err := InsertTestData(context.Background(), client); err != nil {
		log.Fatalf("failed inserting test data: %v", err)
	}

	log.Println("✅ All test data inserted successfully!")
}

// InsertTestData はランダムなテストデータを全エンティティに挿入
func InsertTestData(ctx context.Context, client *ent.Client) error {
	log.Println("📝 Creating test data...")

	// 1. ジャンルを作成
	genres, err := createGenres(ctx, client)
	if err != nil {
		return err
	}
	log.Printf("✅ Created %d genres", len(genres))

	// 2. ユーザーを作成（画像は後で設定）
	users, err := createUsers(ctx, client, 5)
	if err != nil {
		return err
	}
	log.Printf("✅ Created %d users", len(users))

	// 3. ユーザーにジャンルを関連付け
	if err := associateUsersWithGenres(ctx, client, users, genres); err != nil {
		return err
	}
	log.Println("✅ Associated users with genres")

	// 4. フォロー関係を作成
	if err := createFollowRelationships(ctx, client, users); err != nil {
		return err
	}
	log.Println("✅ Created follow relationships")

	// 5. 目標を作成
	goals, err := createGoals(ctx, client, users)
	if err != nil {
		return err
	}
	log.Printf("✅ Created %d goals", len(goals))

	// 6. 投稿を作成（一部は目標に関連付け）
	posts, err := createPosts(ctx, client, users, goals)
	if err != nil {
		return err
	}
	log.Printf("✅ Created %d posts", len(posts))

	// 7. 画像を作成（投稿に関連付け、プロフィール画像も含む）
	images, err := createImages(ctx, client, users, posts)
	if err != nil {
		return err
	}
	log.Printf("✅ Created %d images", len(images))

	// 8. プロフィール画像を設定
	if err := setProfilePictures(ctx, client, users, images); err != nil {
		return err
	}
	log.Println("✅ Set profile pictures")

	// 9. リアクションを作成
	if err := createReactions(ctx, client, users, posts); err != nil {
		return err
	}
	log.Println("✅ Created reactions")

	return nil
}

// createGenres はサンプルジャンルを作成
func createGenres(ctx context.Context, client *ent.Client) ([]*ent.Genre, error) {
	genreNames := []string{"音楽", "スポーツ", "料理", "プログラミング", "読書", "旅行", "映画", "ゲーム"}
	genres := make([]*ent.Genre, 0, len(genreNames))

	for _, name := range genreNames {
		g, err := client.Genre.Create().SetName(name).Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating genre %s: %w", name, err)
		}
		genres = append(genres, g)
	}

	return genres, nil
}

// createUsers はランダムなユーザーを作成
func createUsers(ctx context.Context, client *ent.Client, count int) ([]*ent.User, error) {
	users := make([]*ent.User, 0, count)

	for i := 0; i < count; i++ {
		name := RandomString(8)
		email := fmt.Sprintf("%s@example.com", RandomString(12))

		builder := client.User.Create().
			SetName(name).
			SetEmail(email)

		// 一部のユーザーにはオプショナル情報を追加
		if rand.Float32() > 0.5 {
			birthday := time.Now().AddDate(-rand.Intn(30)-20, -rand.Intn(12), -rand.Intn(28))
			builder = builder.SetBirthday(birthday)
		}
		if rand.Float32() > 0.5 {
			builder = builder.SetHometown(RandomString(10) + "県")
		}
		if rand.Float32() > 0.5 {
			builder = builder.SetBio(RandomString(50))
		}

		u, err := builder.Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed creating user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// associateUsersWithGenres はユーザーとジャンルを関連付け
func associateUsersWithGenres(ctx context.Context, client *ent.Client, users []*ent.User, genres []*ent.Genre) error {
	for _, u := range users {
		// 各ユーザーにランダムに2-4個のジャンルを関連付け
		numGenres := rand.Intn(3) + 2
		selectedGenres := make([]*ent.Genre, 0, numGenres)

		for i := 0; i < numGenres && i < len(genres); i++ {
			selectedGenres = append(selectedGenres, genres[rand.Intn(len(genres))])
		}

		_, err := client.User.UpdateOneID(u.ID).
			AddGenres(selectedGenres...).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed associating genres with user: %w", err)
		}
	}

	return nil
}

// createFollowRelationships はフォロー関係を作成
func createFollowRelationships(ctx context.Context, client *ent.Client, users []*ent.User) error {
	for _, u := range users {
		// 各ユーザーがランダムに1-3人をフォロー
		numFollowing := rand.Intn(3) + 1

		for i := 0; i < numFollowing && i < len(users); i++ {
			followTarget := users[rand.Intn(len(users))]
			// 自分自身はフォローしない
			if followTarget.ID == u.ID {
				continue
			}

			_, err := client.User.UpdateOneID(u.ID).
				AddFollowing(followTarget).
				Save(ctx)
			if err != nil {
				// すでにフォロー済みの場合はスキップ
				continue
			}
		}
	}

	return nil
}

// createGoals は目標を作成
func createGoals(ctx context.Context, client *ent.Client, users []*ent.User) ([]*ent.Goal, error) {
	goalTitles := []string{"毎日運動する", "英語を習得する", "新しいスキルを学ぶ", "健康的な生活を送る", "本を10冊読む"}
	goals := make([]*ent.Goal, 0)

	for _, u := range users {
		// 各ユーザーに1-2個の目標を作成
		numGoals := rand.Intn(2) + 1

		for i := 0; i < numGoals; i++ {
			title := goalTitles[rand.Intn(len(goalTitles))] + " - " + RandomString(5)
			builder := client.Goal.Create().
				SetTitle(title).
				SetUser(u)

			// 一部の目標には期限を設定
			if rand.Float32() > 0.5 {
				deadline := time.Now().AddDate(0, rand.Intn(6)+1, 0)
				builder = builder.SetDeadline(deadline)
			}

			g, err := builder.Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed creating goal: %w", err)
			}
			goals = append(goals, g)
		}
	}

	return goals, nil
}

// createPosts は投稿を作成
func createPosts(ctx context.Context, client *ent.Client, users []*ent.User, goals []*ent.Goal) ([]*ent.Post, error) {
	posts := make([]*ent.Post, 0)

	for _, u := range users {
		// 各ユーザーに2-5個の投稿を作成
		numPosts := rand.Intn(4) + 2

		for i := 0; i < numPosts; i++ {
			content := RandomString(50) + "という内容の投稿です。"
			builder := client.Post.Create().
				SetContent(content).
				SetUser(u)

			// 一部の投稿は目標に関連付け
			if rand.Float32() > 0.6 && len(goals) > 0 {
				// このユーザーの目標を探す
				userGoals, err := client.Goal.Query().
					Where(goal.HasUserWith(user.ID(u.ID))).
					All(ctx)
				if err == nil && len(userGoals) > 0 {
					builder = builder.SetGoal(userGoals[rand.Intn(len(userGoals))])
				}
			}

			p, err := builder.Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed creating post: %w", err)
			}
			posts = append(posts, p)
		}
	}

	return posts, nil
}

// createImages は画像を作成
func createImages(ctx context.Context, client *ent.Client, users []*ent.User, posts []*ent.Post) ([]*ent.Image, error) {
	images := make([]*ent.Image, 0)

	// 各投稿に0-3枚の画像を追加
	for _, p := range posts {
		numImages := rand.Intn(4) // 0-3枚

		postUser, err := p.QueryUser().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying post user: %w", err)
		}

		for i := 0; i < numImages; i++ {
			objectName := fmt.Sprintf("images/%s.jpg", uuid.New().String())
			img, err := client.Image.Create().
				SetObjectName(objectName).
				SetContentType("image/jpeg").
				SetUploadedBy(postUser).
				SetPost(p).
				Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed creating image: %w", err)
			}
			images = append(images, img)
		}
	}

	// プロフィール画像用の画像を追加（投稿に紐付かない）
	for _, u := range users {
		if rand.Float32() > 0.3 { // 70%のユーザーにプロフィール画像を作成
			objectName := fmt.Sprintf("images/profile_%s.jpg", uuid.New().String())
			img, err := client.Image.Create().
				SetObjectName(objectName).
				SetContentType("image/jpeg").
				SetUploadedBy(u).
				Save(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed creating profile image: %w", err)
			}
			images = append(images, img)
		}
	}

	return images, nil
}

// setProfilePictures はプロフィール画像を設定
func setProfilePictures(ctx context.Context, client *ent.Client, users []*ent.User, images []*ent.Image) error {
	for _, u := range users {
		// このユーザーがアップロードした投稿に紐付かない画像を取得
		profileImages, err := client.Image.Query().
			Where(
				image.HasUploadedByWith(user.ID(u.ID)),
				image.Not(image.HasPost()),
			).
			All(ctx)

		if err != nil {
			return fmt.Errorf("failed querying profile images: %w", err)
		}

		if len(profileImages) > 0 {
			// 最初の画像をプロフィール画像に設定
			_, err := client.User.UpdateOneID(u.ID).
				SetProfilePictureID(profileImages[0].ID).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed setting profile picture: %w", err)
			}
		}
	}

	return nil
}

// createReactions はリアクションを作成
func createReactions(ctx context.Context, client *ent.Client, users []*ent.User, posts []*ent.Post) error {
	for _, u := range users {
		// 各ユーザーがランダムに3-7個の投稿にリアクション
		numReactions := rand.Intn(5) + 3

		for i := 0; i < numReactions && i < len(posts); i++ {
			targetPost := posts[rand.Intn(len(posts))]

			// 自分の投稿にはリアクションしない確率を高く
			if rand.Float32() > 0.8 {
				postUser, err := targetPost.QueryUser().Only(ctx)
				if err == nil && postUser.ID == u.ID {
					continue
				}
			}

			_, err := client.Reaction.Create().
				SetUser(u).
				SetPost(targetPost).
				Save(ctx)
			if err != nil {
				// すでにリアクション済みの場合はスキップ（ユニーク制約違反）
				continue
			}
		}
	}

	return nil
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// getEnv は環境変数を取得し、存在しない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
