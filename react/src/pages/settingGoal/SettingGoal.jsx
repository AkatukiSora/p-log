import { useState } from "react"
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import GoalsData from "../../components/GoalsData/GoalsData"

export default function SettingGoal() {

    const [goalDatas, setGoalsDatas] = useState(["aa"])
    
    return (
        <>
            <title>目標の追加</title>
                <header className="header">
                    <h1>目標の追加</h1>
                </header>
            <main className="main">
                <h2>現在の目標</h2>
                <GoalsData goalsData={goalDatas}/>
            </main>
            <footer>
                <ButtonBottomOption />
            </footer>
        </>
    )
}
