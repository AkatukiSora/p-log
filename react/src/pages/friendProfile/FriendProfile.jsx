import ButtonBottomOption from '../../components/button/buttonBottomOption/ButtonBottomOption';
import ButtonTopFriendOption from '../../components/button/buttonTopFriendOption/ButtonTopFriendOption';
import styles from './FriendProfile.module.css';

export default function FriendProfile() {
  // フレンドのプロフィール
  return (
    <>
      <title>ユーザ名のプロフィール</title>
      <header className="header friendHeader">
        <h1>ユーザ名のプロフィール</h1>
        <ButtonTopFriendOption />
      </header>
      <main className="main friendMain">
        <div className={styles.icon}>
          {/* フレンドのアイコンを挿入予定 */}
          <img src="" alt="アイコンの写真" />
        </div>
        <div className={styles.profile}>
          <div className={styles.row}>
            <span className={styles.label}>名前</span>
            <span className={styles.value}>yamada ryousei</span>
          </div>
          <div className={styles.row}>
            <span className={styles.label}>誕生日</span>
            <span className={styles.value}>2006/10/8</span>
          </div>
          <div className={styles.row}>
            <span className={styles.label}>ジャンル</span>
            <span className={styles.value}>全般</span>
          </div>
          <div className={styles.row}>
            <span className={styles.label}>出身</span>
            <span className={styles.value}>愛知</span>
          </div>
          <div className={styles.row}>
            <span className={styles.label}>コメント</span>
            <span className={styles.value}>
              各値は直で書きました。アイコンは自分のプロフィールのアイコンのようになるはずです。
            </span>
          </div>
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  );
}
