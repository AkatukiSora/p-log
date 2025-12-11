import ButtonBottomOption from "../../components/buttonBottomOption/h/ButtonBottomOption"
import ButtonTopMyOption from "../../components/buttonTopMyOption/ButtonTopMyOption"

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
