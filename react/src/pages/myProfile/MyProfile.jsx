import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption.jsx"

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
