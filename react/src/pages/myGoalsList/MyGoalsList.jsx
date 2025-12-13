import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption"
import { useState } from 'react'
import styles from "./MyGoalsList.module.css"
export const goalsList = [
{
  userId:1,
  id:1,
  content:"ハッカソン出場",
  limit:new Date("2025-12-14"),
  progress:[{
    title:"ミーティング",
    date:new Date("2025-12-8"),
  },],
},
{
  userId:1,
  id:2,
  content:"応用情報合格",
  limit:new Date("2026-4-1"),
  progress:[
    {
      title:"過去問200問やった",
      date:new Date("2025-10-1"),
    },
    {
      title:"模擬試験",
      date:new Date("2025-10-12"),
    },
  ],
},
{
  userId:1,
  id:3,
  content:"雀豪",
  limit:new Date("2026-12-31"),
  progress:[],
}
];
const MyGoalsList = () => {
    // 自身の目標一覧(ホーム)

  const [goals,setGoals]=useState(goalsList);
  const deleteGoal=(id)=>{
    const newGoals=goals.filter((goal)=>{
      return goal.id!==id;
    });
    setGoals(newGoals)
  };

  const createGoal=(goals)=>{
    setGoals([...goals,goal]);
  }

  const complete=(id)=>{
    deleteGoal(id)
  }
  return (
    <>
      <title>目標一覧</title>
      <header className="header myHeader">
        <h1>目標一覧</h1>
        <ButtonTopMyOption />
      </header>
      <main className="main myMain">
        <div className={styles.default}>
          <div>         
            {goals.map((goal)=>{
              return(
                <div key={goal.id}>
                            <span className={styles.option}>
                              <input type="checkbox" onChange={() => complete(goal.id)}/>
                              {goal.content}
                              <small>期限: {new Date(goal.limit).toLocaleDateString()}</small>
                            </span>
                              <ul>
                                {goal.progress.map((pro) => (
                                  <li className={styles.list} key={pro.date}>{pro.title} <small>({new Date(pro.date).toLocaleDateString()})</small></li>
                                ))}
                              </ul>
                            
                        </div>
                    )
                  })}
              </div>
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  )
}

export default MyGoalsList
