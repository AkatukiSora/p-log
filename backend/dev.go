package main

import (
	"context"
	"log"
	"time"

	"backend/ent"
	"backend/ent/genre"
	"backend/internal/jwt"

	"github.com/google/uuid"
)

// dev は開発用の初期データ投入関数です。
func dev(client *ent.Client, jwtHandler *jwt.JwtHandler) error {
	ctx := context.Background()

	// === ジャンルの作成 ===
	genres := []string{"筋トレ", "プログラミング", "読書", "料理", "語学学習", "音楽"}
	createdGenres := make([]*ent.Genre, 0, len(genres))

	log.Println("Creating genres...")
	for _, genreName := range genres {
		g, err := client.Genre.
			Create().
			SetName(genreName).
			Save(ctx)
		if err != nil {
			// すでに存在する場合は取得
			g, err = client.Genre.Query().Where(genre.NameEQ(genreName)).Only(ctx)
			if err != nil {
				log.Printf("failed to create or get genre %s: %v", genreName, err)
				continue
			}
		}
		createdGenres = append(createdGenres, g)
		log.Printf("Genre created: %s", genreName)
	}

	// === ユーザーの作成 ===
	type UserData struct {
		ID       uuid.UUID
		Name     string
		Email    string
		Bio      string
		Hometown string
		Genres   []*ent.Genre
	}

	users := []UserData{
		{
			ID:       uuid.MustParse("3fa85f64-5717-4562-b3fc-2c963f66afa6"),
			Name:     "田中太郎",
			Email:    "tanaka@example.com",
			Bio:      "健康的な生活を目指して日々トレーニング中です！",
			Hometown: "東京都",
			Genres:   createdGenres[0:2], // 筋トレ、プログラミング
		},
		{
			ID:       uuid.New(),
			Name:     "佐藤花子",
			Email:    "sato@example.com",
			Bio:      "読書と料理が趣味です。毎日新しいレシピに挑戦しています。",
			Hometown: "大阪府",
			Genres:   createdGenres[2:4], // 読書、料理
		},
		{
			ID:       uuid.New(),
			Name:     "鈴木一郎",
			Email:    "suzuki@example.com",
			Bio:      "フルスタックエンジニア。趣味でギターも弾きます。",
			Hometown: "神奈川県",
			Genres:   []*ent.Genre{createdGenres[1], createdGenres[5]}, // プログラミング、音楽
		},
		{
			ID:       uuid.New(),
			Name:     "山田美咲",
			Email:    "yamada@example.com",
			Bio:      "英語とスペイン語を勉強中。海外旅行が夢です。",
			Hometown: "京都府",
			Genres:   []*ent.Genre{createdGenres[4], createdGenres[2]}, // 語学学習、読書
		},
	}

	createdUsers := make([]*ent.User, 0, len(users))
	log.Println("Creating users...")

	for _, userData := range users {
		// ユーザーが存在しない場合は作成する
		user, err := client.User.Get(ctx, userData.ID)
		if err != nil {
			// Get で見つからない場合は作成を試みる
			birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
			user, err = client.User.
				Create().
				SetID(userData.ID).
				SetName(userData.Name).
				SetEmail(userData.Email).
				SetBio(userData.Bio).
				SetHometown(userData.Hometown).
				SetBirthday(birthday).
				AddGenres(userData.Genres...).
				Save(ctx)
			if err != nil {
				log.Printf("failed to create user %s: %v", userData.Name, err)
				continue
			}
		}
		createdUsers = append(createdUsers, user)
		log.Printf("User created: %s (%s)", userData.Name, userData.Email)
	}

	// === ゴールの作成 ===
	log.Println("Creating goals...")
	type GoalData struct {
		Title    string
		Deadline *time.Time
		UserIdx  int
	}

	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	threeMonths := now.AddDate(0, 3, 0)

	goals := []GoalData{
		{Title: "ベンチプレス100kg達成", Deadline: &threeMonths, UserIdx: 0},
		{Title: "毎日10km走る習慣をつける", Deadline: &nextMonth, UserIdx: 0},
		{Title: "Go言語でWebアプリ開発", Deadline: &threeMonths, UserIdx: 0},
		{Title: "月に10冊本を読む", Deadline: &nextMonth, UserIdx: 1},
		{Title: "フランス料理のコース料理を作る", Deadline: &threeMonths, UserIdx: 1},
		{Title: "Reactでポートフォリオサイト作成", Deadline: &nextMonth, UserIdx: 2},
		{Title: "ギターで好きな曲を10曲マスター", Deadline: nil, UserIdx: 2},
		{Title: "TOEIC900点突破", Deadline: &threeMonths, UserIdx: 3},
		{Title: "スペイン語検定3級取得", Deadline: &threeMonths, UserIdx: 3},
	}

	createdGoals := make([]*ent.Goal, 0, len(goals))
	for _, goalData := range goals {
		if goalData.UserIdx >= len(createdUsers) {
			continue
		}

		goalBuilder := client.Goal.
			Create().
			SetTitle(goalData.Title).
			SetUser(createdUsers[goalData.UserIdx])

		if goalData.Deadline != nil {
			goalBuilder = goalBuilder.SetDeadline(*goalData.Deadline)
		}

		goal, err := goalBuilder.Save(ctx)
		if err != nil {
			log.Printf("failed to create goal %s: %v", goalData.Title, err)
			continue
		}
		createdGoals = append(createdGoals, goal)
		log.Printf("Goal created: %s (User: %s)", goalData.Title, createdUsers[goalData.UserIdx].Name)
	}

	// === 投稿の作成 ===
	log.Println("Creating posts...")
	type PostData struct {
		Content string
		UserIdx int
		GoalIdx int
	}

	posts := []PostData{
		{Content: "今日からジム通い開始！頑張るぞー💪", UserIdx: 0, GoalIdx: 0},
		{Content: "ベンチプレス80kgまで上げられるようになった！", UserIdx: 0, GoalIdx: 0},
		{Content: "朝ランニング5km完走。気持ちいい！", UserIdx: 0, GoalIdx: 1},
		{Content: "Go言語の基礎文法を学習中。シンプルで書きやすい！", UserIdx: 0, GoalIdx: 2},
		{Content: "村上春樹の新作を読了。深い物語だった。", UserIdx: 1, GoalIdx: 3},
		{Content: "今月5冊目の本を読了！順調に進んでいる。", UserIdx: 1, GoalIdx: 3},
		{Content: "手作りパスタに挑戦！思ったより美味しくできた。", UserIdx: 1, GoalIdx: 4},
		{Content: "Reactのhooksの使い方を勉強中。便利！", UserIdx: 2, GoalIdx: 5},
		{Content: "ポートフォリオサイトのデザイン案完成。明日からコーディング。", UserIdx: 2, GoalIdx: 5},
		{Content: "ギターの練習1時間。指が痛い...でも楽しい！", UserIdx: 2, GoalIdx: 6},
		{Content: "TOEIC模試で850点取れた！あと少し。", UserIdx: 3, GoalIdx: 7},
		{Content: "毎日30分英語のポッドキャストを聞くようにしている。", UserIdx: 3, GoalIdx: 7},
		{Content: "スペイン語の基本的な挨拶をマスターした！¡Hola!", UserIdx: 3, GoalIdx: 8},
	}

	for _, postData := range posts {
		if postData.UserIdx >= len(createdUsers) || postData.GoalIdx >= len(createdGoals) {
			continue
		}

		post, err := client.Post.
			Create().
			SetContent(postData.Content).
			SetUser(createdUsers[postData.UserIdx]).
			SetGoal(createdGoals[postData.GoalIdx]).
			Save(ctx)
		if err != nil {
			log.Printf("failed to create post: %v", err)
			continue
		}
		log.Printf("Post created: %s (User: %s)", post.Content[:20]+"...", createdUsers[postData.UserIdx].Name)
	}

	// === 最初のユーザーのトークン生成 ===
	if len(createdUsers) > 0 {
		firstUser := createdUsers[0]
		accessToken, refreshToken, err := jwtHandler.GenerateTokens(firstUser.ID, firstUser.Email, ctx)
		if err != nil {
			log.Fatalf("failed to generate tokens: %v", err)
		}
		log.Println("\n=== Dev User Credentials ===")
		log.Printf("User: %s (%s)", firstUser.Name, firstUser.Email)
		log.Printf("User ID: %s", firstUser.ID)
		log.Printf("AccessToken: %s", accessToken)
		log.Printf("RefreshToken: %s", refreshToken)
		log.Println("============================")
	}

	log.Println("\n✅ Dev data setup completed!")
	log.Printf("Created: %d genres, %d users, %d goals, %d posts",
		len(createdGenres), len(createdUsers), len(createdGoals), len(posts))

	return nil
}
