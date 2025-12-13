import { Link } from 'react-router-dom';
import { ROUTES } from '../../../pages/pathData/path';
import styles from './ButtonTopFriendOption.module.css';
import { useState } from 'react';

export default function ButtonTopFriendOption() {
  // テスト用
  const [flg, setFlg] = useState(false);
  const [text, setText] = useState('追加');

  // フレンドの追加
  const addFriend = () => {
    setFlg(!flg);
    // まだフレンドを追加していなかったら追加する
    if (flg === true) {
      alert('フレンドを削除');
      setText('追加');
    }
    // すでにフレンドを追加していたら削除する
    else {
      alert('フレンドを追加');
      setText('削除');
    }
  };

  // フレンド機能の上部のボタン
  return (
    <div className={styles.options}>
      <Link className={styles.link} to={ROUTES.friendProfile}>
        <div className={styles.option}>
          <p>プロフィール</p>
        </div>
      </Link>
      <Link className={styles.link} to={ROUTES.friendGoals}>
        <div className={styles.option}>
          <p>目標</p>
        </div>
      </Link>
      <div className={styles.option}>
        <p onClick={addFriend}>フレンド{text}</p>
      </div>
    </div>
  );
}
