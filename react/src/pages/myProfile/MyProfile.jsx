import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption"

export default function MyProfile() {
    // 自身のプロフィール画面
  return (
    <>
        <title>プロフィール</title>
        <header className="header myHeader">
            <h1>プロフィール</h1>
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
