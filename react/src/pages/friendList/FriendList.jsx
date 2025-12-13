import { Link } from "react-router-dom"
import { ROUTES } from "../pathData/path"
import styles from "./FriendList.module.css"
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"

export default function FriendList() {
    // フレンドの選択画面
  return (
    <>
      <title>フレンド</title>
      <header className="header">
        <h1>フレンド</h1>
        <div className={styles.search}>
          <input type="text"></input>
          <button>検索</button>
        </div>
      </header>
      <main className="main">
        <Link to={ROUTES.friendProfile}>アイコンの代わり</Link>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  )
}
