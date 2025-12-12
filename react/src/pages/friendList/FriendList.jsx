import { Link } from "react-router-dom"
import { ROUTES } from "../pathData/path"
import styles from "./FriendList.module.css"
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"

export default function FriendList() {

    // テストデータ
    const datas = [
      {id : "1",name : "a"},
      {id : "2",name : "b"},
      {id : "3",name : "c"},
      {id : "4",name : "d"},
      {id : "5",name : "e"},
      {id : "6",name : "f"}
    ]

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
                {/* 5列x行で表示 */}
                <div className={styles.gridContainer}>
                  {datas.map(test => (
                    <Link to={ROUTES.friendProfile} key={test.id} className={styles.link}>
                      <div className={styles.gritItem}>
                        {/* <img src="data.url" alt="アイコン"> */}
                        ここにアイコン
                      </div>
                    </Link>
                  ))}
                </div>
            </main>
            <footer>
                <ButtonBottomOption />
            </footer>
        </>
    )
}
