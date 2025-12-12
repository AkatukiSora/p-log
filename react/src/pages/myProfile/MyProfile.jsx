import { useState } from "react";
import ButtonBottomOption from "../../components/button/buttonBottomOption/ButtonBottomOption";
import ButtonTopMyOption from "../../components/button/buttonTopMyOption/ButtonTopMyOption";
import styles from "./MyProfile.module.css";

export default function MyProfile() {

    const [preview, setPreview] = useState(null);

    const previewFile = (e) => {
        const file = e.target.files[0];
        if (!file) return;

        const reader = new FileReader();
        reader.onload = () => {
            setPreview(reader.result);
        };
        reader.readAsDataURL(file);
    };

    const changeProInfo = () => {

    }
    // 名前だけ書いてたら通る。
    return (
        <>
            <title>プロフィール</title>
            <header className="header myHeader">
                <h1>プロフィール</h1>
                <ButtonTopMyOption />
            </header>
            <main className="main myMain">
                {/* 送信はしません */}
                {/* アイコン */}
                <div className={styles.icon}>
                    {preview && (
                        <img className={styles.previewImage} src={preview} alt="preview"/>
                    )}
                    <br ></br>
                    <label htmlFor="icon">アイコンを選択</label>
                    <input type="file" name="icon" id="icon" onChange={previewFile} required />
                </div>
                {/* 名前 */}
                <div className={styles.form}>
                    <label htmlFor="name">名前</label>
                    <input id="name" type="text" name="name" /><br />
                </div>
                {/* 誕生日 */}
                <div className={styles.form}>
                    <label htmlFor="birthday">誕生日</label>
                    <input id="birthday" type="date" name="birthday" className={styles.date}/><br />
                </div>
                {/* ジャンル */}
                <div className={styles.form}>
                    <label htmlFor="genre">ジャンル</label>
                    <select id="genre" name="genre" className={styles.select}>
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
                    <input id="birthPlace" type="text" name="birthPlace" /><br />
                </div>
                {/* 自由記述欄 */}
                <div className={styles.form}>
                    <label htmlFor="free">自由記述欄</label>
                    <textarea/>
                </div>
                {/* 変更ボタン */}
                <div className={styles.form}>
                    <button onClick={changeProInfo}>変更する</button>
                </div>
            </main>
            <footer>
                <ButtonBottomOption />
            </footer>
        </>
    );
}
