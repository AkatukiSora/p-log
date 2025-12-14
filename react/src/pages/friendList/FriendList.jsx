import { Link } from 'react-router-dom';
import { ROUTES } from '../pathData/path';
import styles from './FriendList.module.css';
import ButtonBottomOption from '../../components/button/buttonBottomOption/ButtonBottomOption';
import { useEffect, useState } from 'react';

export default function FriendList() {
  const API_BASE_URL = 'http://localhost:8080/api/v1';
  const MOCK_TOKEN = 'mock_access_token';

  // 写真を入れてないです...

  // フレンドのプロフィール情報
  const [friendProfileDatas, setFriendProfileDatas] = useState([]);
  // 検索に一致したフレンドのプロフィール情報
  const [searchedData, setSearchedData] = useState([]);
  // 検索値
  const [inpSearch, setInpSearch] = useState('');

  useEffect(() => {
    // APIから目標一覧を非同期で取得する
    const fetchFriends = async () => {
      const response = await fetch(`${API_BASE_URL}/friends`, {
        headers: {
          Authorization: `Bearer ${MOCK_TOKEN}`,
        },
      });

      const friendDatas = await response.json();
      const userIds = friendDatas.map((friendData) => friendData);

      const userDataArray = await Promise.all(
        userIds.map(async (userId) => {
          const userResponse = await fetch(`${API_BASE_URL}/users/${userId}`, {
            headers: {
              Authorization: `Bearer ${MOCK_TOKEN}`,
            },
          });
          return userResponse.json();
        })
      );

      // フレンドのプロフィール情報を取得
      const finalData = friendDatas.map((friendData, index) => ({
        userData: userDataArray.find((data) => data.id === friendData),
      }));

      // 名前を取得
      finalData.map((item) => {
        setFriendProfileDatas((prev) => [...prev, item.userData.name]);
      });
    };
    fetchFriends();
  }, []);

  // 検索地に関連するデータを抽出
  const search = () => {
    if (inpSearch.trim() === '') {
      setSearchedData([]);
      return;
    }

    const filteredData = friendProfileDatas.filter((data) =>
      data.toLowerCase().includes(inpSearch.toLowerCase())
    );

    setSearchedData(filteredData);
  };

  // 検索結果のデータ
  const displayDatas = friendProfileDatas.filter((data) =>
    data.toLowerCase().includes(inpSearch.toLowerCase())
  );

  // フレンドの選択画面
  return (
    <>
      <title>フレンド</title>
      <header className="header">
        <h1>フレンド</h1>
        <div className={styles.search}>
          <input type="text" onChange={(e) => setInpSearch(e.target.value)}></input>
          <button onClick={search}>検索</button>
        </div>
      </header>
      <main className="main">
        {/* 5列x行で表示 */}
        <div className={styles.gridContainer}>
          {displayDatas.map((data, index) => (
            <Link to={ROUTES.friendProfile} key={index} className={styles.link}>
              <div className={styles.gritItem}>{data}</div>
            </Link>
          ))}
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  );
}
