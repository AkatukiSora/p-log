import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption"
import { useState,useEffect } from 'react'
import styles from "./MyGoalsList.module.css"

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

const MyGoalsList = () => {
    // 自身の目標一覧(ホーム)
  const [goals,setGoals]=useState([]);

  useEffect(()=>{
    const fetchGoals=async()=>{
      const res = await fetch(`${API_BASE_URL}/goals`,{
        headers: {
          "Authorization":`Bearer ${MOCK_TOKEN}`,
        },
      });
      const data = await res.json();
      setGoals(data);
    }
    fetchGoals();
  },[]);
  const deleteGoal=async(id)=>{
    const res = await fetch(`${API_BASE_URL}/goals/${id}`,{
      method: "DELETE",
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${MOCK_TOKEN}`,
      },
    });
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
                              {goal.title}
                              <small>期限: {new Date(goal.deadline).toLocaleDateString()}</small>
                            </span>
                              {/* <ul>
                                {goal.progress.map((pro) => (
                                  <li className={styles.list} key={pro.date}>{pro.title} <small>({new Date(pro.date).toLocaleDateString()})</small></li>
                                ))}
                              </ul> */}
                            
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
