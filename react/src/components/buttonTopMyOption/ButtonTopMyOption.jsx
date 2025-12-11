import { Link } from "react-router-dom"
import styles from "./ButtonTopMyOption.module.css"
import { ROUTES } from "../../pages/pathData/path"

export default function ButtonTopMyOption() {
  return (
    <>
      <div className={styles.options}>
        <Link className={styles.link} to={ROUTES.settingGoal}>
            <div className={styles.option}>
                <p>目標の追加</p>
            </div>
        </Link>
        <Link className={styles.link} to={ROUTES.addPost}>
            <div className={styles.option}>
                <p>進捗の投稿</p>
            </div>
        </Link>
      </div>
    </>
  )
}
