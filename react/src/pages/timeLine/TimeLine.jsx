import ButtonBottomOption from '../../components/button/buttonBottomOption/ButtonBottomOption.jsx';
import { useState } from 'react';
import styles from './TimeLine.module.css';
import { BsEmojiNeutral } from 'react-icons/bs';
import { BsEmojiGrinFill } from 'react-icons/bs';
export default function TimeLine() {
  const postsList = [
    {
      userId: 1,
      id: 1,
      title: '応用情報合格',
      content: '模擬試験受験',
      date: new Date('2025-10-12'),
    },
    {
      userId: 1,
      id: 2,
      title: '応用情報合格',
      content: '過去問200問やった',
      date: new Date('2025-10-5'),
    },
    {
      userId: 2,
      id: 3,
      title: 'AtCoder入茶',
      content: 'C問題解いた',
      date: new Date('2025-10-1'),
    },
    {
      userId: 3,
      id: 4,
      title: '貯金50万円',
      content: '5万円入金',
      date: new Date('2025-9-30'),
    },
  ];

  const users = {
    1: {
      name: 'たま',
      icon: './azarashi.png',
    },
    2: {
      name: 'フリテン',
      icon: './ika.png',
    },
    3: {
      name: 'もっさり',
      icon: './pig.png',
    },
  };
  const initialLikeCounts = { 1: 5, 2: 12, 3: 2, 4: 0 };
  const [likeCounts, setLikeCounts] = useState(initialLikeCounts);

  const initialUserLikes = { 1: false, 2: true, 3: false, 4: false };
  const [userLiked, setUserLiked] = useState(initialUserLikes);

  const handleLikeToggle = (id) => {
    const isCurrentlyLiked = userLiked[id];

    setUserLiked((prev) => ({
      ...prev,
      [id]: !isCurrentlyLiked,
    }));

    setLikeCounts((prevCounts) => ({
      ...prevCounts,
      [id]: isCurrentlyLiked ? prevCounts[id] - 1 : prevCounts[id] + 1,
    }));
  };
  // タイムライン
  return (
    <>
      <title>タイムライン</title>
      <header className="header">
        <h1>タイムライン</h1>
      </header>
      <main className="main">
        <div className={styles.default}>
          {postsList.map((post) => {
            const user = users[post.userId];
            const Liked = userLiked[post.id];
            return (
              <div key={post.id} className={styles.post}>
                <div>
                  <img src={user.icon} alt={`${user.name}のアイコン`} className={styles.usericon} />
                  <strong>{user.name}</strong>
                </div>
                <div>
                  <div className={styles.option}>{post.title}</div>
                  <div className={styles.post}>
                    <div className={styles.p}>
                      {post.content}
                      <small>({new Date(post.date).toLocaleDateString()})</small>
                    </div>
                    <button onClick={() => handleLikeToggle(post.id)} className={styles.button}>
                      {Liked ? <BsEmojiGrinFill /> : <BsEmojiNeutral />}
                      {likeCounts[post.id]}
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  );
}
