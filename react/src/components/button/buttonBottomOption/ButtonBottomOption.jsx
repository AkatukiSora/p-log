import { Link } from 'react-router-dom';
import styles from './ButtonBottomOption.module.css';
import { ROUTES } from '../../../pages/pathData/path';

export default function ButtonBottomOption() {
  // 下部の4つのボタン
  return (
    <div className={styles.options}>
      <Link className={styles.link} to={ROUTES.myList}>
        <div className={styles.option}>
          <p>目標一覧</p>
        </div>
      </Link>
      <Link className={styles.link} to={ROUTES.friendList}>
        <div className={styles.option}>
          <p>フレンド</p>
        </div>
      </Link>
      <Link className={styles.link} to={ROUTES.timeLine}>
        <div className={styles.option}>
          <p>タイムライン</p>
        </div>
      </Link>
      <Link className={styles.link} to={ROUTES.myProfile}>
        <div className={styles.option}>
          <p>プロフィール</p>
        </div>
      </Link>
    </div>
  );
}
