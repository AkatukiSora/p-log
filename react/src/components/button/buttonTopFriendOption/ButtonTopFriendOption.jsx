import { Link } from "react-router-dom"
import { ROUTES } from "../../../pages/pathData/path"
import styles from "./ButtonTopFriendOption.module.css"

export default function ButtonTopFriendOption() {
    // 開発途中
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
            <p>フレンド追加</p>
        </div>
    </div>
  )
}