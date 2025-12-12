import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption"

export default function MyList() {
    // 自身の目標一覧(ホーム)
  return (
    <>
      <title>目標一覧</title>
      <header className="header myHeader">
        <h1>目標一覧</h1>
        <ButtonTopMyOption />
      </header>
      <main className="main myMain">
        aa
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  )
}