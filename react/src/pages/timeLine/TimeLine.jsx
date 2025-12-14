import ButtonBottomOption from "../../components//button/buttonBottomOption/ButtonBottomOption.jsx"
import { useState,useEffect, use } from "react"
import styles from "./TimeLine.module.css"
import {BsEmojiNeutral} from "react-icons/bs"
import { BsEmojiGrinFill } from "react-icons/bs";

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

export default function TimeLine() {
  const [posts,setPosts]=useState([]);
  const [reaction,setReaction]=useState([]);

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

      const initialReaction=data.map((d)=>{
        return{
          postId: d.id,
          isReacted: false,
          count: d.postData.reaction_count,
        };
      });
      setReaction(initialReaction)
    };
    fetchTimeLine();

  },[]);
  const onReactionClick =async (post) => {
    const temp=post.count
    const updatedReaction=reaction.map((r)=>{
      if(r.postId === post.id){
        r.isReacted ? temp -1 : temp + 1
        return {
          postId:r.postId,
          isReacted: !r.isReacted,
          count: r.isReacted ? r.count - 1 : r.count + 1
        }
      }
      return r;
    })
    setReaction(updatedReaction)
    const newPost = {
      id: post.id,
      user_id: post.user_id,
      goal_id: post.goal_id,
      content: post.content,
      image_urls: post.image_urls,
      reaction_count: temp,
      created_at: post.created_at,
      updated_at: post.updated_at
    }

    const res = await fetch(`${API_BASE_URL}/posts/${post.id}`,{
      method: "PUT",
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${MOCK_TOKEN}`,
      },
      body: JSON.stringify(newPost),
    });
    console.log("res:",res)
  }
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
            const reacted = reaction.find((r)=>r.postId === post.id)?.isReacted
            const reactionCount = reaction.find((r)=>r.postId === post.id)?.count
            return(
              <div key={post.id} className={styles.post}>
                  <div>
                    <img src={post.icon} alt={`${post.userData.name}のアイコン`} className={styles.usericon}/>
                    <strong>{post.userData.name}</strong>
                  </div>
                  <div>
                    <div className={styles.option}>{post.goalData.title}</div>
                  <div className={styles.post}>
                    <div>{post.content}<small>({new Date(post.created_at).toLocaleDateString()})</small></div>
                    <button
                      onClick={() =>{
                        console.log('Clicked:', post.id) 
                        onReactionClick(post)}}
                      className={styles.button}
                    >
                      {reacted ? <BsEmojiGrinFill/> : <BsEmojiNeutral />}
                      {reactionCount}
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
