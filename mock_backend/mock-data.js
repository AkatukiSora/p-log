const { v4: uuidv4 } = require('uuid');

// 固定のUUID
const MOCK_USER_ID = '123e4567-e89b-12d3-a456-426614174000';
const MOCK_USER_ID_2 = '223e4567-e89b-12d3-a456-426614174001';
const MOCK_GOAL_ID = '323e4567-e89b-12d3-a456-426614174002';
const MOCK_GOAL_ID_2 = '423e4567-e89b-12d3-a456-426614174003';
const MOCK_POST_ID = '523e4567-e89b-12d3-a456-426614174004';
const MOCK_POST_ID_2 = '623e4567-e89b-12d3-a456-426614174005';
const MOCK_IMAGE_ID = '723e4567-e89b-12d3-a456-426614174006';
const MOCK_GENRE_ID_1 = '823e4567-e89b-42d3-a456-426614174007';
const MOCK_GENRE_ID_2 = '923e4567-e89b-42d3-a456-426614174008';

// ジャンルマスターデータ
const genres = [
  { id: MOCK_GENRE_ID_1, name: 'プログラミング' },
  { id: MOCK_GENRE_ID_2, name: '語学学習' },
  { id: 'a23e4567-e89b-42d3-a456-426614174009', name: '資格取得' },
  { id: 'b23e4567-e89b-42d3-a456-42661417400a', name: 'フィットネス' },
  { id: 'c23e4567-e89b-42d3-a456-42661417400b', name: 'アート・デザイン' },
];

// ユーザーデータ
const users = {
  [MOCK_USER_ID]: {
    id: MOCK_USER_ID,
    name: 'テストユーザー',
    birthday: '1995-05-15',
    genres: [MOCK_GENRE_ID_1, MOCK_GENRE_ID_2],
    hometown: '東京都',
    bio: 'プログラミングと語学学習を頑張っています！',
  },
  [MOCK_USER_ID_2]: {
    id: MOCK_USER_ID_2,
    name: 'サンプルユーザー',
    birthday: '1998-03-20',
    genres: [MOCK_GENRE_ID_2],
    hometown: '大阪府',
    bio: '毎日コツコツ継続中です',
  },
};

// 目標データ
const goals = {
  [MOCK_GOAL_ID]: {
    id: MOCK_GOAL_ID,
    user_id: MOCK_USER_ID,
    title: 'Go言語をマスターする',
    created_at: '2024-01-01T00:00:00Z',
    deadline: '2025-12-31',
  },
  [MOCK_GOAL_ID_2]: {
    id: MOCK_GOAL_ID_2,
    user_id: MOCK_USER_ID,
    title: '英語でプレゼンできるようになる',
    created_at: '2024-01-15T00:00:00Z',
    deadline: '2025-06-30',
  },
};

// 投稿データ
const posts = {
  [MOCK_POST_ID]: {
    id: MOCK_POST_ID,
    user_id: MOCK_USER_ID,
    goal_id: MOCK_GOAL_ID,
    content: '今日はGoの並行処理について学習しました。goroutineとchannelの使い方がだんだん理解できてきました！',
    image_urls: [],
    reaction_count: 5,
    created_at: '2024-12-10T10:00:00Z',
    updated_at: '2024-12-10T10:00:00Z',
  },
  [MOCK_POST_ID_2]: {
    id: MOCK_POST_ID_2,
    user_id: MOCK_USER_ID,
    goal_id: MOCK_GOAL_ID_2,
    content: '英語のプレゼン資料を作成中。まだまだですが頑張ります！',
    image_urls: [],
    reaction_count: 3,
    created_at: '2024-12-11T15:30:00Z',
    updated_at: '2024-12-11T15:30:00Z',
  },
};

// リアクションデータ
const reactions = {
  [MOCK_POST_ID]: [
    { user_id: MOCK_USER_ID_2, created_at: '2024-12-10T11:00:00Z' },
    { user_id: '333e4567-e89b-12d3-a456-426614174010', created_at: '2024-12-10T12:00:00Z' },
  ],
  [MOCK_POST_ID_2]: [
    { user_id: MOCK_USER_ID_2, created_at: '2024-12-11T16:00:00Z' },
  ],
};

// フレンドデータ（ユーザーIDをキーにしたフレンドIDの配列）
const friends = {
  [MOCK_USER_ID]: [MOCK_USER_ID_2],
  [MOCK_USER_ID_2]: [MOCK_USER_ID],
};

// 認証トークン
const authToken = {
  access_token: 'mock_access_token_' + Date.now(),
  token_type: 'Bearer',
  refresh_token: 'mock_refresh_token_' + Date.now(),
  expires_in: new Date(Date.now() + 3600000).toISOString(), // 1時間後のISO 8601形式
};

// 画像データ（実際のバイナリではなくURLのみ）
const images = {
  [MOCK_IMAGE_ID]: {
    id: MOCK_IMAGE_ID,
    url: `http://localhost:8080/images/${MOCK_IMAGE_ID}`,
  },
};

module.exports = {
  MOCK_USER_ID,
  MOCK_USER_ID_2,
  MOCK_GOAL_ID,
  MOCK_GOAL_ID_2,
  MOCK_POST_ID,
  MOCK_POST_ID_2,
  MOCK_IMAGE_ID,
  MOCK_GENRE_ID_1,
  MOCK_GENRE_ID_2,
  genres,
  users,
  goals,
  posts,
  reactions,
  friends,
  authToken,
  images,
};
