import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"
import ButtonTopFriendOption from "../../components/button/buttonTopFriendOption/ButtonTopFriendOption.jsx"
import { useEffect, useState } from "react"
import styles from "./friendGoalsList.module.css"

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';
const friend_id = "223e4567-e89b-12d3-a456-426614174001"

export default function friendGoalsList() {

    const [friendGoals,setFriendGoals]=useState([]);
    useEffect(()=>{
        const fetchGoals=async()=>{
          const res = await fetch(`${API_BASE_URL}/users/${friend_id}/goals`,{
            headers: {
              "Authorization":`Bearer ${MOCK_TOKEN}`,
            },
          });
          const data = await res.json();
          setFriendGoals(data);
        }
        fetchGoals();
      },[]);
    // フレンドの目標一覧
  return (
    <>
        <title>ユーザ名の目標一覧</title>
        <header className="header friendHeader">
            <h1>ユーザ名の目標一覧</h1>
            <ButtonTopFriendOption />
        </header>
        <main className="main friendMain">
            <div className={styles.default}>
                    {friendGoals.map((goal)=>{
                            return(
                                <div key={goal.id}>
                                    <span className={styles.option}>
                                      {goal.title}
                                      <small>期限: {new Date(goal.deadline).toLocaleDateString()}</small>
                                    </span>
                                    {/* <ul>
                                      {goal.progress.map((pro) => (
                                        <li key={pro.date}>{pro.title} <small>({new Date(pro.date).toLocaleDateString()})</small></li>
                                      ))}
                                    </ul> */}
                                </div>
                            )
                    })}
                    </div>
        </main>
        <footer>
            <ButtonBottomOption />
        </footer>
    </>
  )
}