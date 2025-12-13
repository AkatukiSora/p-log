import { useState } from 'react';
import ButtonBottomOption from '../../components/button/buttonBottomOption/ButtonBottomOption';
import ButtonTopMyOption from '../../components/button/buttonTopMyOption/ButtonTopMyOption';
import styles from './MyProfile.module.css';

export default function MyProfile() {
  // プロフィール情報のステート
  const [previewData, setPreview] = useState(null);
  const [nameData, setNameData] = useState('');
  const [birthData, setBirthData] = useState('');
  const [genre, setGenreData] = useState('');
  const [birthPlaceData, setBirthPlaceData] = useState('');
  const [commentData, setCommentData] = useState('');

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
  const changeProInfo = () => {
    // アイコンと名前を書いてたら通る。
    if (previewData === null || nameData === '') {
      alert('アイコンと名前を入力してください');
      return;
    }
  };

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
          {previewData && (
            <img className={styles.previewImage} src={previewData} alt="previewData" />
          )}
          <br />
          <label htmlFor="icon">アイコンを選択</label>
          <input type="file" name="icon" id="icon" onChange={previewFile} required />
        </div>
        {/* 名前 */}
        <div className={styles.form}>
          <label htmlFor="name">名前</label>
          <input
            type="text"
            name="name"
            id="name"
            value={nameData}
            onChange={(e) => setNameData(e.target.value)}
          />
          <br />
        </div>
        {/* 誕生日 */}
        <div className={styles.form}>
          <label htmlFor="birthday">誕生日</label>
          <input
            type="date"
            name="birthday"
            id="birthday"
            value={birthData}
            className={styles.date}
            onChange={(e) => setBirthData(e.target.value)}
          />
          <br />
        </div>
        {/* ジャンル */}
        <div className={styles.form}>
          <label htmlFor="genre">ジャンル</label>
          <select
            name="genre"
            id="genre"
            value={genre}
            className={styles.select}
            onChange={(e) => setGenreData(e.target.value)}
          >
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
          <input
            type="text"
            name="birthPlace"
            id="birthPlace"
            value={birthPlaceData}
            onChange={(e) => setBirthPlaceData(e.target.value)}
          />
          <br />
        </div>
        {/* 自由記述欄 */}
        <div className={styles.form}>
          <label htmlFor="comment">コメント</label>
          <textarea
            name="comment"
            id="comment"
            value={commentData}
            onChange={(e) => setCommentData(e.target.value)}
          />
        </div>
        {/* 変更ボタン */}
        <div className={styles.form}>
          <button className={styles.button} onClick={changeProInfo}>
            変更する
          </button>
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  );
}
