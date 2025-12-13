import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"
import ButtonTopFriendOption from "../../components/button/buttonTopFriendOption/ButtonTopFriendOption.jsx"
import { useState } from "react"
import styles from "./friendGoalsList.module.css"
export default function friendGoalsList() {
    const goalsList = [
{
  userId:2,
  id:1,
  content:"筋トレ",
  limit:new Date("2026-7-31"),
  progress:[
    {
    title:"プランク10分",
    date:new Date("2025-8-2"),
    },
    {
    title:"スクワット10回",
    date:new Date("2025-9-17"),
    },
],
},
{
  userId:2,
  id:2,
  content:"G検定合格",
  limit:new Date("2026-5-1"),
  progress:[
    {
      title:"黒本100問",
      date:new Date("2025-9-30"),
    },
    {
      title:"あんちょこ作成",
      date:new Date("2025-10-22"),
    },
  ],
},
{
  userId:1,
  id:3,
  content:"課題",
  limit:new Date("2025-12-31"),
  progress:[],
}
];
    const [friendGoals,setFriendGoals]=useState(goalsList);

    // フレンドの目標一覧
  return (
    <>
        <title>ユーザ名の目標一覧</title>
        <header className="header friendHeader">
            <h1>ユーザ名の目標一覧</h1>
            <ButtonTopFriendOption />
        </header>
        <main className="main friendMain">
            <div>
                    {friendGoals.map((goal)=>{
                            return(
                                <div key={goal.id}>
                                    <span className={styles.option}>
                                      {goal.content}
                                      <small>期限: {new Date(goal.limit).toLocaleDateString()}</small>
                                    </span>
                                    {goal.progress && goal.progress.length > 0 ? (
                                      <ul>
                                        {goal.progress.map((pro) => (
                                          <li key={pro.date}>{pro.title} <small>({new Date(pro.date).toLocaleDateString()})</small></li>
                                        ))}
                                      </ul>
                                    ):""
                                    }
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