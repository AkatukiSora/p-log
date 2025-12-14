import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption"
import ButtonTopFriendOption from "../../components/button/buttonTopFriendOption/ButtonTopFriendOption"
import styles from "./FriendProfile.module.css"
import { useEffect, useState } from "react"

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';
const friend_id = "223e4567-e89b-12d3-a456-426614174001"

export default function FriendProfile({users}) {
    
    const [fp,setFriendProfile]=useState([]);
    useEffect(()=>{
            const fetchProfile=async()=>{
              const res = await fetch(`${API_BASE_URL}/users/${friend_id}`,{
                headers: {
                  "Authorization":`Bearer ${MOCK_TOKEN}`,
                },
              });

              const iconData=await fetch(`${API_BASE_URL}/users/${friend_id}/icon`,{
                headers: {
                  "Authorization":`Bearer ${MOCK_TOKEN}`,
                },
              })
              const data = await res.json();
              const finalData = {...data,iconData}
              setFriendProfile(finalData);
              console.log(finalData)
            }
            fetchProfile();
          },[]);
    // フレンドのプロフィール
  return (
    <>
        <title>ユーザ名のプロフィール</title>
        <header className="header friendHeader">
            <h1>ユーザ名のプロフィール</h1>
            <ButtonTopFriendOption />
        </header>
            <main className="main friendMain">
                <div className={styles.icon}>
                    {/* フレンドのアイコンを挿入予定 */}
                    <img src={fp.icon} alt={`${fp.name}アイコンの写真`} />
                </div>
                <div className={styles.profile}>
                    <div className={styles.row}>
                        <span className={styles.label}>名前</span>
                        <span className={styles.value}>{fp.name}</span>
                    </div>
                    <div className={styles.row}>
                        <span className={styles.label}>誕生日</span>
                        <span className={styles.value}>{fp.birthday}</span>
                    </div>
                    <div className={styles.row}>
                        <span className={styles.label}>ジャンル</span>
                        <span className={styles.value}>{fp.genres}</span>
                    </div>
                    <div className={styles.row}>
                        <span className={styles.label}>出身</span>
                        <span className={styles.value}>{fp.hometown}</span>
                    </div>
                    <div className={styles.row}>
                        <span className={styles.label}>コメント</span>
                        <span className={styles.value}>{fp.bio}</span>
                    </div>
                </div>
            </main>
        <footer>
            <ButtonBottomOption />
        </footer>
    </>
  )
}