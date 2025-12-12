import { useState } from "react"
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import GoalsData from "../../components/GoalsData/GoalsData"
import styles from "./SettingGoal.module.css"

export default function SettingGoal() {

    const [goalDatas, setGoalDatas] = useState([])
    const [addText, setAddText] = useState("")
    const [addDate,setAddDate] = useState("")
    
    const addGoal = () => {
        if(addText === "" || addDate === ""){
            alert("新しい目標名と日付を入力してください")
            return
        }
        setGoalDatas(prev => [...prev,{text : addText, date : addDate}])
        setAddText("")
        setAddDate("")
    }

    return (
        <>
            <div className={styles.global}>
                <title>目標の追加</title>
                    <header className="header">
                        <h1>目標の追加</h1>
                    </header>
                <main className="main">
                    <h2>新しく目標を追加</h2>
                    <label >
                        <input className={styles.inputText} type="text" value={addText} placeholder="ここに入力" onChange={prev => setAddText(prev.target.value)}/>
                        <input className={styles.inputDate}type="date" value={addDate} onChange={prev => setAddDate(prev.target.value)}/>
                    </label>
                    <button className={styles.button} onClick={addGoal}>追加する</button>
                    <h2>現在の目標</h2>
                    <GoalsData goalDatas={goalDatas}/>
                </main>
                <footer>
                    <ButtonBottomOption />
                </footer>
            </div>
        </>
    )
}
