import { useEffect, useState } from "react";
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption";
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption";
import styles from "./MyProfile.module.css";

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

export default function MyProfile() {

    const [myProfileDatas, setMyProfileDatas] = useState(null)
    const [myProfileIconData, setMyProfileIconData] = useState("")

    useEffect(() => {
        const fetchMyProfile = async () => {
            const response = await fetch(`${API_BASE_URL}/auth/me`, {
                headers: {
                    'Authorization': `Bearer ${MOCK_TOKEN}`,
                },
            });
            const data = await response.json();
            setMyProfileDatas(data);

            const iconResponse = await fetch(
                `${API_BASE_URL}/users/${data.id}/icon`,
                {
                    headers: {
                        'Authorization': `Bearer ${MOCK_TOKEN}`,
                    },
                }
            );
            const iconData = await iconResponse.json();
            setMyProfileIconData(iconData.url);
        };

        fetchMyProfile();
    }, []);

    
    // 取得したプロフィールを各入力欄に反映
    useEffect(() => {
        if (!myProfileDatas) return;

        setNameData(myProfileDatas.name ?? "");
        setBirthData(myProfileDatas.birthday ?? "");
        setGenreData(myProfileDatas.genres ?? "");
        setBirthPlaceData(myProfileDatas.hometown ?? "");
        setCommentData(myProfileDatas.bio ?? "");
    }, [myProfileDatas]);


    // プロフィール情報のステート
    const [previewData, setPreview] = useState(null);
    const [nameData, setNameData] = useState("")
    const [birthData, setBirthData] = useState("")
    const [genre, setGenreData] = useState("")
    const [birthPlaceData, setBirthPlaceData] = useState("")
    const [commentData, setCommentData] = useState("")

    // アイコン画像の表示
    const previewFile = (e) => {
        const file = e.target.files[0];
        if (!file) return;
        
        const reader = new FileReader();
        reader.onload = () => {
            setPreview(reader.result);
        };
        reader.readAsDataURL(file);
    };
    
    // 変更ボタンの処理
    const changeProInfo = async () => {
        
        // アイコンと名前を書いてたら通る。
        if(previewData === null || nameData === "") {
            alert("アイコンと名前を入力してください")
            return
        }




        // 未実装
        // 新規登録はpostで、二回目以降はputメソッドを使用する
        const Profile = {
            "name": nameData,
            "birthday": birthData, 
            "genres": genre,
            "hometown": birthPlaceData,
            "bio": commentData
        }

        console.log('Profile', Profile);

        // プロフィール情報をpostメソッドで送信
        const res = await fetch(`${API_BASE_URL}/users`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
        body: JSON.stringify(Profile),
        });

        console.log('res', res);


        const Icon = {
            "icon" : previewData ?? myProfileIconData
        }

        // アイコン情報をpostメソッドで送信
        const resIcon = await fetch(`${API_BASE_URL}/users`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
        body: JSON.stringify(Icon),
        });

        console.log("icon",resIcon)
    }

    return (
        <>
            <title>プロフィール</title>
            <header className="header myHeader">
                <h1>プロフィール</h1>
                <ButtonTopMyOption />
            </header>
            <main className="main myMain">
                {/* アイコン */}
                <div className={styles.icon}>
                    {(previewData || myProfileIconData) && (
                        <img className={styles.previewImage} src={previewData ?? myProfileIconData} alt="profile icon"
                        />
                    )}
                    <br />
                    <label htmlFor="icon">アイコンを選択</label>
                    <input type="file" id="icon" onChange={previewFile} />
                </div>
                {/* 名前 */}
                <div className={styles.form}>
                    <label htmlFor="name" >名前</label>
                    <input type="text" name="name" id="name" value={nameData} onChange={e => setNameData(e.target.value)}/><br />
                </div>
                {/* 誕生日 */}
                <div className={styles.form}>
                    <label htmlFor="birthday">誕生日</label>
                    <input type="date" name="birthday" id="birthday" value={birthData} className={styles.date} onChange={e => setBirthData(e.target.value)}/><br />
                </div>
                {/* ジャンル */}
                <div className={styles.form}>
                    <label htmlFor="genre">ジャンル</label>
                    <select name="genre" id="genre" value={genre} className={styles.select} onChange={e => setGenreData(e.target.value)}>
                        <option value="0">全般</option>
                        <option value="1">IT</option>
                        <option value="2">勉強</option>
                        <option value="3">課題</option>
                        <option value="4">筋トレ</option>
                        <option value="5">読書</option>
                        <option value="6">やる気駆動開発</option>
                        <option value="7">その他</option>
                    </select>
                </div>
                {/* 出身 */}
                <div className={styles.form}>
                    <label htmlFor="birthPlace">出身</label>
                    <input type="text" name="birthPlace" id="birthPlace" value={birthPlaceData} onChange={e => setBirthPlaceData(e.target.value)}/><br />
                </div>
                {/* 自由記述欄 */}
                <div className={styles.form}>
                    <label htmlFor="comment">コメント</label>
                    <textarea name="comment" id="comment" value={commentData} onChange={e => setCommentData(e.target.value)}/>
                </div>
                {/* 変更ボタン */}
                <div className={styles.form}>
                    <button className={styles.button} onClick={changeProInfo}>変更する</button>
                </div>
            </main>
            <footer>
                <ButtonBottomOption />
            </footer>
        </>
    );
}