import ButtonBottomOption from "../../components//button/buttonBottomOption/ButtonBottomOption.jsx"
import { useState,useEffect, use } from "react"
import styles from "./TimeLine.module.css"
import {BsEmojiNeutral} from "react-icons/bs"
import { BsEmojiGrinFill } from "react-icons/bs";

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

export default function TimeLine() {
  // const postsList = [
  //   {
  //     userId:1,
  //     id:1,
  //     title:"応用情報合格",
  //     content:"模擬試験受験",
  //     date:new Date("2025-10-12"),
  //   },
  //   {
  //     userId:1,
  //     id:2,
  //     title:"応用情報合格",
  //     content:"過去問200問やった",
  //     date:new Date("2025-10-5"),
  //   },
  //   {
  //     userId:2,
  //     id:3,
  //     title:"AtCoder入茶",
  //     content:"C問題解いた",
  //     date:new Date("2025-10-1"),
  //   },
  //   {
  //     userId:3,
  //     id:4,
  //     title:"貯金50万円",
  //     content:"5万円入金",
  //     date:new Date("2025-9-30"),
  //   },
  // ]

  // const users ={
  //   1: {
  //     name: "たま",
  //     icon: "./azarashi.png",
  //   },
  //   2: {
  //     name: "フリテン",
  //     icon: "./ika.png",
  //   },
  //   3: {
  //     name: "もっさり",
  //     icon: "./pig.png",
  //   },
  // }

  const [userLiked, setUserLiked] = useState({});

  const [posts,setPosts]=useState([]);

  useEffect(()=>{
    const fetchTimeLine=async()=>{
      const timelineRes = await fetch(`${API_BASE_URL}/timeline`,{
        headers: {
          'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
      });

      const res = await timelineRes.json();
      const userIds = res.map(item => item.user_id);

      const userDataArray=await Promise.all(userIds.map(async (user_id)=>{
        const userResponse = await fetch(`${API_BASE_URL}/users/${user_id}`,{
          headers: {
            'Authorization': `Bearer ${MOCK_TOKEN}`,
          },
        });
        return userResponse.json();
      }))

      const iconArray=await Promise.all(userIds.map(async (user_id)=>{
        const iconResponse=await fetch(`${API_BASE_URL}/users/${user_id}/icon`,{
          headers: {
            'Authorization': `Bearer ${MOCK_TOKEN}`,
          },
        });
        return iconResponse;
      }

      ))

      const goalIds=res.map(item => item.goal_id);
      const titleArray = await Promise.all(goalIds.map(async (goalId)=>{
        const postResponse = await fetch(`${API_BASE_URL}/goals/${goalId}`,{
          headers: {
            'Authorization': `Bearer ${MOCK_TOKEN}`,
          },
        });
        return postResponse.json();

      }))

      const ids=res.map(item => item.id);
      const postArray= await Promise.all(ids.map(async (id)=>{
        const postReasponse = await fetch(`${API_BASE_URL}/posts/${id}`,{
          headers: {
            'Authorization': `Bearer ${MOCK_TOKEN}`,
          },
        });
        return postReasponse.json();
      }))
      console.log(iconArray)
      console.log(userDataArray)
      console.log(titleArray)
      console.log(res)
      const data = res.map((item)=>({
        ...item,
        userData: userDataArray.find((data)=> data.id === item.user_id),
        goalData: titleArray.find((data)=> data.id === item.goal_id),
        postData: postArray.find((data)=> data.id === item.id),
        icon: iconArray.find((data)=>data.user_id === item.user_id),
      }));
      setPosts(data);
      console.log(data)
    };
    fetchTimeLine();

  },[]);
  const handleLikeToggle = (id) => {

    const newLiked= {
      id:id,
      isReacted:!isReacted,
      count: isReacted ? count - 1 : count + 1,
    }
    setUserLiked(newLiked)
  };
    // タイムライン
  return (
    <>
      <title>タイムライン</title>
      <header className="header">
        <h1>タイムライン</h1>
      </header>
      <main className="main">
        <div className={styles.default}>
          {posts.map((post)=>{
            const reaction = 
            return(
              <div key={post.id} className={styles.post}>
                  <div>
                    <img src={post.icon} alt={`${post.userData.name}のアイコン`} className={styles.usericon}/>
                    <strong>{post.userData.name}</strong>
                  </div>
                  <div>
                    <div className={styles.option}>{post.goalData.title}</div>
                  <div className={styles.post}>
                    <div className={styles.p}>{post.content}<small>({new Date(post.created_at).toLocaleDateString()})</small></div>
                    <button
                      onClick={() => handleLikeToggle(post.id)}
                      className={styles.button}
                    >
                      {reaction.isReacted ? <BsEmojiGrinFill/> : <BsEmojiNeutral />}{post.postData.reaction_count}
                    </button>
                  </div>
                  </div>
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
