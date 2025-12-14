import { useState, useEffect } from 'react';
import ButtonBottomOption from '../../components/button/buttonBottomOption/ButtonBottomOption';
import ButtonTopMyOption from '../../components/button/buttonTopMyOption/ButtonTopMyOption';
import styles from './Goals.module.css';

const API_BASE_URL = 'http://localhost:8080/api/v1';
const MOCK_TOKEN = 'mock_access_token';

/**
 * APIから取得した目標一覧を表示するページ
 */
const Goals = () => {
  const [goals, setGoals] = useState([]);
  const [reaction, setReaction] = useState([]);

  useEffect(() => {
    // APIから目標一覧を非同期で取得する
    const fetchGoals = async () => {
      const response = await fetch(`${API_BASE_URL}/goals`, {
        headers: {
          'Authorization': `Bearer ${MOCK_TOKEN}`,
        },
      }); 

      const data = await response.json();
      setGoals(data);

      const tempGoals = data; 

      const initialReaction = tempGoals.map((goal) => {
        return {
          goalId: goal.id,
          isReacted: false,
        };
      });
      setReaction(initialReaction)
    };

    fetchGoals();
  }, []);

  const onClick = async () => {
    const newGoal = {
      title: "新しい目標",
      deadline: "2024-12-31", 
    }

    console.log('newGoal', newGoal);

    const res = await fetch(`${API_BASE_URL}/goals`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${MOCK_TOKEN}`,
      },
      body: JSON.stringify(newGoal),
    });

    console.log('res', res);
  }

    /**
   * @param {string | null} dateString
   * @returns {string}
   */
  const formatDate = (dateString) => {
    if (!dateString) {
      return '未設定';
    }
    return new Date(dateString).toLocaleDateString('ja-JP');
  };

  const onReactionClick = async (clickedGoalId) => {
    
    // await fetch(`${API_BASE_URL}/goals/${clickedGoalId}/reaction`, {
    //   method: 'POST',
    //   headers: {
    //     'Content-Type': 'application/json',
    //     'Authorization': `Bearer ${MOCK_TOKEN}`,
    //   },
    // });

    const updatedReaction = reaction.map((reaction) => {
      if (reaction.goalId === clickedGoalId) {
        return {
          goalId: reaction.goalId,
          isReacted: !reaction.isReacted,
        };
      }
      return reaction;
    });
    setReaction(updatedReaction);  
  }

  return (
    <>
      <title>目標一覧（API）</title>
      <header className="header myHeader">
        <h1>目標一覧</h1>
        <ButtonTopMyOption />
      </header>
      <main className="main myMain">
        <div className={styles.container}>
          <ul className={styles.goalList}>
            {goals.map((goal) => {
              const isReacted = reaction.find((reaction) => reaction.goalId === goal.id)?.isReacted

              return (
              <li key={goal.id} className={styles.goalItem}>
                <div className={styles.goalTitle}>{goal.title}</div>
                <div className={styles.goalMeta}>
                  <span>期限: {formatDate(goal.deadline)}</span>
                  <span>作成日: {formatDate(goal.created_at)}</span>
                </div>
                <button type='button'  onClick={() => onReactionClick(goal.id)}>                 
                {isReacted ? "❤️" : "♡"}
                </button>
              </li>
            )}
            )}
          </ul>
        </div>
      </main>
      <footer>
        <ButtonBottomOption />
      </footer>
    </>
  );
};

export default Goals;
