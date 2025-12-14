import { useState,useEffect } from "react";
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"
import styles from "./AddPost.module.css"

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';
const MOCK_USER_ID = '123e4567-e89b-12d3-a456-426614174000';

// const initGoals=[
//   {
//     id:1,
//     title:"個人開発する",
//     progress:[
//       {
//         title:"デザイン作成",
//         percentage: 80,
//         imageUrl:""
//       }
//     ],
//   },
//   {
//     id:2,
//     title:"ハッカソン参加",
//     progress:[
//       {
//         title:"アイデア出し",
//         percentage: 50,
//         imageUrl:""
//       }
//     ],
//   },
// ]

export default function AddPost() {
  const [goals,setGoals]=useState([]);
  const [goalsTitle,setGoalsTitle]=useState([]);
  const [selectedGoalId,setSelectedGoalId]=useState(null);
  const [postContent,setPostContent]=useState("");
  const [postPercent,setPostPercent]=useState(0);
  const [postImageUrl,setImageUrl]=useState("");

  useEffect(()=>{
      const fetchGoalsTitle = async () => {
      const response = await fetch(`${API_BASE_URL}/goals`, {
        headers: {
          'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
      });    

      const data = await response.json();
      setGoalsTitle(data);
    };

    fetchGoalsTitle();
    const fetchPosts=async()=>{
      const res = await fetch(`${API_BASE_URL}/posts`,{
        headers: {
          "Authorization":`Bearer ${MOCK_TOKEN}`,
        },
      });
      const data = await res.json();
      setGoals(data); 
    }
    fetchPosts();
  },[]);
  const addProgress=async(e)=>{
    if(postContent === "" || postPercent === ""){
      alert("新しい進捗を入力してください")
      return
    }
    e.preventDefault();
    const newProgress={
      id:Math.floor(Math.random() * 1e5),
      user_id: MOCK_USER_ID,
      goal_id: selectedGoalId,
      content: postContent,
      // percentage:postPercent,
      image_urls: postImageUrl,
      reaction_count: 5,
      created_at: Date.now(),
      updated_at: Date.now(),
    };
    const res = await fetch(`${API_BASE_URL}/posts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${MOCK_TOKEN}`,
      },
      body: JSON.stringify(newProgress),
    });
    console.log('res', res);
    
    setGoals(
      goals.map((goal)=>{
        if(goal.goal_id===selectedGoalId){
          return {
                ...goal,newProgress         
            };
      }
      return goal;
  }))
  setSelectedGoalId(null); 
  setPostContent("");
  setPostPercent(0);
  setImageUrl("");
  }

  const selectedGoal = goals.find(g => g.goal_id === selectedGoalId)

  const openForm = (id) => {
    setSelectedGoalId(id)
  }

  return (
    <>
        <title>進捗の投稿</title>
        <header className="header">
            <h1>進捗の投稿</h1>
        </header>
        <main className="main addPostMain">
          <div>
              <div className={styles.post}>
                {goals.map((goal)=>{
                  return(
                    <div key={goal.goal_id}>
                      <div className={styles.option}>
                        <a onClick={()=>openForm(goal.goal_id)}>
                        {goalsTitle.map((g)=>{
                          return(
                            <span>
                              {g.id===goal.goal_id && g.title}
                            </span>
                          )
                        })}
                        </a>
                      </div>
                      {/* <small>期限: {new Date(goal.limit).toLocaleDateString()}</small> */}
                    </div>
                  )
                })}
              </div>
          </div>
                      {selectedGoalId !== null && selectedGoal && (
          <div className={styles.test}>
              <div>
                <h2>投稿一覧</h2>
                <ul>
                  <li className={styles.fukidashi}>{selectedGoal.content}</li>
                </ul>
                <p>新規投稿</p>
                <form onSubmit={addProgress}>
                  <input
                    type="text"
                    value={postContent}
                    onChange={(e) => setPostContent(e.target.value)}
                  />
                  <p>進捗率</p>
                  <input
                    type="range"
                    min="0"
                    max="100"
                    value={postPercent}
                    onChange={(e) => setPostPercent(e.target.value)}
                  /><span>{postPercent}%</span><br />
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setImageUrl(e.target.files[0] ? e.target.files[0].name : "")}
                  />
                  <button className={styles.button}>投稿</button>
                </form>
              </div>
          </div>
            )}
        </main>
        <footer>
            <ButtonBottomOption />
        </footer>
    </>
  )
}