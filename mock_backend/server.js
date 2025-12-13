const express = require('express');
const cors = require('cors');
const multer = require('multer');
const { v4: uuidv4 } = require('uuid');
const mockData = require('./mock-data');

const app = express();
const PORT = process.env.PORT || 8080;
const API_PREFIX = '/api/v1';

// ミドルウェア設定
app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// multer設定（画像アップロード用）
const storage = multer.memoryStorage();
const upload = multer({ storage: storage });

// リクエストログ
app.use((req, res, next) => {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.path}`);
  next();
});

// 簡易的なバリデーションヘルパー
const validateUUID = (id) => {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
  return uuidRegex.test(id);
};

const validateRequired = (fields, body) => {
  const missing = fields.filter(field => !body[field]);
  return missing.length === 0 ? null : `必須フィールドが不足しています: ${missing.join(', ')}`;
};

// 認証ミドルウェア（簡易版）
// モック用のため、実際のトークン検証は行いません
const authMiddleware = (req, res, next) => {
  const authHeader = req.headers.authorization;
  if (!authHeader || !authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ message: '認証が必要です' });
  }
  // モックとして固定ユーザーIDを設定（実際の検証はなし）
  req.userId = mockData.MOCK_USER_ID;
  next();
};

// ===== Auth エンドポイント =====
app.get(`${API_PREFIX}/auth/login`, (req, res) => {
  // OIDCプロバイダーへのリダイレクトを模擬
  res.redirect(301, `http://localhost:${PORT}${API_PREFIX}/auth/callback?code=mock_code&state=mock_state`);
});

app.get(`${API_PREFIX}/auth/callback`, (req, res) => {
  const { code, state } = req.query;
  
  if (!code || !state) {
    return res.status(400).json({ message: 'codeとstateパラメータが必要です' });
  }
  
  res.json(mockData.authToken);
});

app.post(`${API_PREFIX}/auth/logout`, (req, res) => {
  res.status(204).send();
});

app.get(`${API_PREFIX}/auth/me`, authMiddleware, (req, res) => {
  const user = mockData.users[req.userId];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  res.json(user);
});

// ===== Genres エンドポイント =====
app.get(`${API_PREFIX}/genres`, (req, res) => {
  res.json(mockData.genres);
});

// ===== Users エンドポイント =====
app.post(`${API_PREFIX}/users`, (req, res) => {
  const error = validateRequired(['name'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const newUser = {
    id: uuidv4(),
    name: req.body.name,
    birthday: req.body.birthday || null,
    genres: req.body.genres || [],
    hometown: req.body.hometown || null,
    bio: req.body.bio || null,
  };
  
  res.status(201).json(newUser);
});

app.get(`${API_PREFIX}/users/:user_id`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  res.json(user);
});

app.put(`${API_PREFIX}/users/:user_id`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const error = validateRequired(['name'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  const updatedUser = {
    ...user,
    name: req.body.name,
    birthday: req.body.birthday || user.birthday,
    genres: req.body.genres || user.genres,
    hometown: req.body.hometown || user.hometown,
    bio: req.body.bio || user.bio,
  };
  
  res.json(updatedUser);
});

app.delete(`${API_PREFIX}/users/:user_id`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  res.status(204).send();
});

// ===== User Icon エンドポイント =====
app.post(`${API_PREFIX}/users/:user_id/icon`, authMiddleware, upload.single('icon'), (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  if (!req.file) {
    return res.status(400).json({ message: 'アイコンファイルが必要です' });
  }
  
  res.status(204).send();
});

app.get(`${API_PREFIX}/users/:user_id/icon`, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  // ダミー画像を返す（1x1の透明なPNG）
  const dummyImage = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==', 'base64');
  res.contentType('image/png');
  res.send(dummyImage);
});

app.delete(`${API_PREFIX}/users/:user_id/icon`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  res.status(204).send();
});

// ===== Goals エンドポイント =====
app.post(`${API_PREFIX}/goals`, authMiddleware, (req, res) => {
  const error = validateRequired(['title'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const newGoal = {
    id: uuidv4(),
    user_id: req.userId,
    title: req.body.title,
    created_at: new Date().toISOString(),
    deadline: req.body.deadline || null,
  };
  
  res.status(201).json(newGoal);
});

app.get(`${API_PREFIX}/goals`, authMiddleware, (req, res) => {
  const { page = 1, limit = 20 } = req.query;
  
  // ユーザーの目標一覧を返す
  const userGoals = Object.values(mockData.goals).filter(
    goal => goal.user_id === req.userId
  );
  
  res.json(userGoals);
});

app.get(`${API_PREFIX}/goals/:goal_id`, (req, res) => {
  const { goal_id } = req.params;
  
  if (!validateUUID(goal_id)) {
    return res.status(400).json({ message: 'goal_idの形式が不正です' });
  }
  
  const goal = mockData.goals[goal_id];
  if (!goal) {
    return res.status(404).json({ message: '目標が見つかりません' });
  }
  
  res.json(goal);
});

app.put(`${API_PREFIX}/goals/:goal_id`, authMiddleware, (req, res) => {
  const { goal_id } = req.params;
  
  if (!validateUUID(goal_id)) {
    return res.status(400).json({ message: 'goal_idの形式が不正です' });
  }
  
  const error = validateRequired(['title'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const goal = mockData.goals[goal_id];
  if (!goal) {
    return res.status(404).json({ message: '目標が見つかりません' });
  }
  
  const updatedGoal = {
    ...goal,
    title: req.body.title,
    deadline: req.body.deadline || goal.deadline,
  };
  
  res.json(updatedGoal);
});

app.delete(`${API_PREFIX}/goals/:goal_id`, authMiddleware, (req, res) => {
  const { goal_id } = req.params;
  
  if (!validateUUID(goal_id)) {
    return res.status(400).json({ message: 'goal_idの形式が不正です' });
  }
  
  const goal = mockData.goals[goal_id];
  if (!goal) {
    return res.status(404).json({ message: '目標が見つかりません' });
  }
  
  res.status(204).send();
});

app.get(`${API_PREFIX}/users/:user_id/goals`, (req, res) => {
  const { user_id } = req.params;
  const { page = 1, limit = 20 } = req.query;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  const userGoals = Object.values(mockData.goals).filter(
    goal => goal.user_id === user_id
  );
  
  res.json(userGoals);
});

// ===== Posts エンドポイント =====
app.post(`${API_PREFIX}/posts`, authMiddleware, (req, res) => {
  const error = validateRequired(['goal_id', 'content'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  if (!validateUUID(req.body.goal_id)) {
    return res.status(400).json({ message: 'goal_idの形式が不正です' });
  }
  
  const newPost = {
    id: uuidv4(),
    user_id: req.userId,
    goal_id: req.body.goal_id,
    content: req.body.content,
    image_urls: req.body.image_ids ? req.body.image_ids.map(id => `http://localhost:${PORT}${API_PREFIX}/images/${id}`) : [],
    reaction_count: 0,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };
  
  res.status(201).json(newPost);
});

app.get(`${API_PREFIX}/posts`, authMiddleware, (req, res) => {
  const { goal_id, page = 1, limit = 20 } = req.query;
  
  let userPosts = Object.values(mockData.posts).filter(
    post => post.user_id === req.userId
  );
  
  if (goal_id) {
    if (!validateUUID(goal_id)) {
      return res.status(400).json({ message: 'goal_idの形式が不正です' });
    }
    userPosts = userPosts.filter(post => post.goal_id === goal_id);
  }
  
  res.json(userPosts);
});

app.get(`${API_PREFIX}/posts/:post_id`, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  res.json(post);
});

app.put(`${API_PREFIX}/posts/:post_id`, authMiddleware, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const error = validateRequired(['goal_id', 'content'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  const updatedPost = {
    ...post,
    goal_id: req.body.goal_id,
    content: req.body.content,
    image_urls: req.body.image_ids ? req.body.image_ids.map(id => `http://localhost:${PORT}${API_PREFIX}/images/${id}`) : post.image_urls,
    updated_at: new Date().toISOString(),
  };
  
  res.json(updatedPost);
});

app.delete(`${API_PREFIX}/posts/:post_id`, authMiddleware, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  res.status(204).send();
});

app.get(`${API_PREFIX}/users/:user_id/posts`, (req, res) => {
  const { user_id } = req.params;
  const { goal_id, page = 1, limit = 20 } = req.query;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  let userPosts = Object.values(mockData.posts).filter(
    post => post.user_id === user_id
  );
  
  if (goal_id) {
    if (!validateUUID(goal_id)) {
      return res.status(400).json({ message: 'goal_idの形式が不正です' });
    }
    userPosts = userPosts.filter(post => post.goal_id === goal_id);
  }
  
  res.json(userPosts);
});

// ===== Timeline エンドポイント =====
app.get(`${API_PREFIX}/timeline`, authMiddleware, (req, res) => {
  const { goal_id, page = 1, limit = 20 } = req.query;
  
  // 自分とフレンドの投稿を取得
  const friendIds = mockData.friends[req.userId] || [];
  const relevantUserIds = [req.userId, ...friendIds];
  
  let timelinePosts = Object.values(mockData.posts).filter(
    post => relevantUserIds.includes(post.user_id)
  );
  
  if (goal_id) {
    if (!validateUUID(goal_id)) {
      return res.status(400).json({ message: 'goal_idの形式が不正です' });
    }
    timelinePosts = timelinePosts.filter(post => post.goal_id === goal_id);
  }
  
  // 新しい順にソート
  timelinePosts.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  
  res.json(timelinePosts);
});

// ===== Images エンドポイント =====
app.post(`${API_PREFIX}/images`, authMiddleware, upload.single('image'), (req, res) => {
  if (!req.file) {
    return res.status(400).json({ message: '画像ファイルが必要です' });
  }
  
  const imageId = uuidv4();
  const newImage = {
    id: imageId,
    url: `http://localhost:${PORT}${API_PREFIX}/images/${imageId}`,
  };
  
  res.status(201).json(newImage);
});

app.get(`${API_PREFIX}/images/:image_id`, (req, res) => {
  const { image_id } = req.params;
  
  if (!validateUUID(image_id)) {
    return res.status(400).json({ message: 'image_idの形式が不正です' });
  }
  
  // ダミー画像を返す（1x1の透明なPNG）
  const dummyImage = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==', 'base64');
  res.contentType('image/png');
  res.send(dummyImage);
});

// ===== Reactions エンドポイント =====
app.post(`${API_PREFIX}/posts/:post_id/reactions`, authMiddleware, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  res.status(201).send();
});

app.delete(`${API_PREFIX}/posts/:post_id/reactions`, authMiddleware, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  res.status(204).send();
});

app.get(`${API_PREFIX}/posts/:post_id/reactions`, (req, res) => {
  const { post_id } = req.params;
  
  if (!validateUUID(post_id)) {
    return res.status(400).json({ message: 'post_idの形式が不正です' });
  }
  
  const post = mockData.posts[post_id];
  if (!post) {
    return res.status(404).json({ message: '投稿が見つかりません' });
  }
  
  const postReactions = mockData.reactions[post_id] || [];
  res.json(postReactions);
});

// ===== Friends エンドポイント =====
app.get(`${API_PREFIX}/friends`, authMiddleware, (req, res) => {
  const userFriends = mockData.friends[req.userId] || [];
  res.json(userFriends);
});

app.post(`${API_PREFIX}/friends`, authMiddleware, (req, res) => {
  const error = validateRequired(['user_id'], req.body);
  if (error) {
    return res.status(400).json({ message: error });
  }
  
  const { user_id } = req.body;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  res.status(201).send();
});

app.get(`${API_PREFIX}/users/:user_id/friends`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  const userFriends = mockData.friends[user_id] || [];
  res.json(userFriends);
});

app.delete(`${API_PREFIX}/friends/:user_id`, authMiddleware, (req, res) => {
  const { user_id } = req.params;
  
  if (!validateUUID(user_id)) {
    return res.status(400).json({ message: 'user_idの形式が不正です' });
  }
  
  const user = mockData.users[user_id];
  if (!user) {
    return res.status(404).json({ message: 'ユーザーが見つかりません' });
  }
  
  res.status(204).send();
});

// ===== エラーハンドリング =====
app.use((req, res) => {
  res.status(404).json({ message: 'リソースが見つかりません' });
});

app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({ message: '内部サーバーエラー' });
});

// サーバー起動
app.listen(PORT, () => {
  console.log(`🚀 Mock API server is running on http://localhost:${PORT}`);
  console.log(`📚 API documentation: http://localhost:${PORT}/docs`);
});
