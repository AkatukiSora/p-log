import { useEffect, useState } from "react"
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import GoalsData from "../../components/goalsData/GoalsData"
import styles from "./SettingGoal.module.css"

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

export default function SettingGoal() {

    const [goalDatas, setGoalDatas] = useState([])
    const [addText, setAddText] = useState("")
    const [addDate,setAddDate] = useState("")
    


    // 自身のプロフィール情報を取得
    useEffect(() => {
        const fetchGoals = async () => {
            const response = await fetch(`${API_BASE_URL}/goals`, {
                headers: {
                    'Authorization': `Bearer ${MOCK_TOKEN}`,
                },
            });
            const datas = await response.json();
            datas.map(data => {
                setGoalDatas(prev => [...prev,{text : data.title, date : data.deadline}])
            })
        };
        fetchGoals();
    }, []);

        // 追加ボタンを押したら追加する
    const addGoal = async () => {
        if(addText === "" || addDate === ""){
            alert("新しい目標名と日付を入力してください")
            return
        }

        setGoalDatas(prev => [...prev,{text : addText, date : addDate}])

        const newGoalPayload = {
            "title" : addText,
            "deadline" : addDate
        }

        console.log('Goals Payload', newGoalPayload);

        const res = await fetch(`${API_BASE_URL}/goals`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
        body: JSON.stringify(newGoalPayload),
        });
        
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
