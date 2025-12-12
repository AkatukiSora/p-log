import { useState } from "react";
import ButtonBottomOption from "../../components/button/buttonBottomOption/h/ButtonBottomOption.jsx"
import styles from "./AddPost.module.css"

const goalsList = [
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

export default function AddPost() {
  const [goals,setGoals]=useState(goalsList);
  const createProgress=(e)=>{
    setSelected(null);
    setPostContent("")
    setImage(null);
    setGoals([...goals,e.target.value]);
  }
  const [selected,setSelected]=useState(null);
  const [postContent,setPostContent]=useState("");
  const [image,setImage]=useState(null);

  const selectedGoal = goalsList.find(g => g.id === selected)

  const openForm = (id) => {
    setSelected(id)
  }


  const PostForm = () => {
    if (!selectedGoal) return null;
    return (
      <div>
        <h2 className={styles.option}>{selectedGoal.content}</h2>
        <div>
            <div>
              <p>投稿文</p>
              <textarea
                            value={postContent}
                            onChange={(e) => setPostContent(e.target.value)}
              />
              <input
                            type="file"
                            accept="image/*"
                            onChange={(e) => setImage(e.target.files[0])}
              />
            {image && (
                            <p>添付ファイル: {image.name}</p>
              )} 
            <button onClick={(e) => createProgress(e.target)}>投稿</button>
            </div>
        </div>
      </div>

    )
  }
  return (
    <>
        <title>進捗の投稿</title>
        <header className="header">
            <h1>進捗の投稿</h1>
        </header>
        <main className="main">
            <div className={styles.post}>
              {goalsList.map((goal)=>{
                return(
                  <div key={goal.id}>
                    <div className={styles.option}>
                      <a onClick={()=>openForm(goal.id)}>
                      {goal.content}
                      </a>
                    </div>
                    <small>期限: {new Date(goal.limit).toLocaleDateString()}</small>
                  </div>
                )
              })}
            </div>
        {selected !== null && <PostForm />}
        </main>
        <footer>
            <ButtonBottomOption />
        </footer>
    </>
  )
}
