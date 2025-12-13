import { useState } from "react";
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption.jsx"
import styles from "./AddPost.module.css"

// const goalsList = [
// {
//   userId:1,
//   id:1,
//   content:"ハッカソン出場",
//   limit:new Date("2025-12-14"),
//   progress:[{
//     title:"ミーティング",
//     date:new Date("2025-12-8"),
//   },],
// },
// {
//   userId:1,
//   id:2,
//   content:"応用情報合格",
//   limit:new Date("2026-4-1"),
//   progress:[
//     {
//       title:"過去問200問やった",
//       date:new Date("2025-10-1"),
//     },
//     {
//       title:"模擬試験",
//       date:new Date("2025-10-12"),
//     },
//   ],
// },
// {
//   userId:1,
//   id:3,
//   content:"雀豪",
//   limit:new Date("2026-12-31"),
//   progress:[],
// }
// ];

const initGoals=[
  {
    id:1,
    title:"個人開発する",
    progress:[
      {
        title:"デザイン作成",
        percentage: 80,
        imageUrl:""
      }
    ],
  },
  {
    id:2,
    title:"ハッカソン参加",
    progress:[
      {
        title:"アイデア出し",
        percentage: 50,
        imageUrl:""
      }
    ],
  },
]

export default function AddPost() {
  const [goals,setGoals]=useState(initGoals);
  const [selectedGoalId,setSelectedGoalId]=useState(null);
  const [postContent,setPostContent]=useState("");
  const [postPercent,setPostPercent]=useState(0);
  const [postImageUrl,setImageUrl]=useState("");

  const addProgress=(e)=>{
    
    e.preventDefault();
    const newProgress={
      title:postContent,
      percentage:postPercent,
      imageUrl: postImageUrl,
    };
    setGoals(
      goals.map((goal)=>{
        if(goal.id===selectedGoalId){
          return {
            ...goal,
            progress:
              [
                ...goal.progress,newProgress
              ]  ,
          };
      }
      return goal;
  }))
  setSelectedGoalId(null); 
  setPostContent("");
  setPostPercent(0);
  setImageUrl("");
  }

  const selectedGoal = goals.find(g => g.id === selectedGoalId)

  const openForm = (id) => {
    setSelectedGoalId(id)
  }

  return (
    <>
        <title>進捗の投稿</title>
        <header className="header">
            <h1>進捗の投稿</h1>
        </header>
        <main className="main">
          <div>
              <div className={styles.post}>
                {initGoals.map((goal)=>{
                  return(
                    <div key={goal.id}>
                      <div className={styles.option}>
                        <a onClick={()=>openForm(goal.id)}>
                        {goal.title}
                        </a>
                      </div>
                      {/* <small>期限: {new Date(goal.limit).toLocaleDateString()}</small> */}
                    </div>
                  )
                })}
              </div>
          </div>
                      {selectedGoalId !== null && selectedGoal && (
          <div>
            <h2 className={styles.option}>{selectedGoal.title}</h2>
              <div>
                <p>投稿文</p>
                <form onSubmit={addProgress}>
                  <input
                    type="text"
                    value={postContent}
                    onChange={(e) => setPostContent(e.target.value)}
                  />
                  <p>進捗率</p>
                  <input
                    type="number"
                    min="0"
                    max="100"
                    value={postPercent}
                    onChange={(e) => setPostPercent(e.target.value)}
                  /><span>%</span><br />
                  <input
                    type="file"
                    accept="image/*"
                    onChange={(e) => setImageUrl(e.target.files[0] ? e.target.files[0].name : "")}
                  />
                  <button>投稿</button>
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