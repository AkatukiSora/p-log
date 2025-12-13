import styles from './GoalData.module.css';

export default function GoalData({ goalData }) {
  return (
    <div className={styles.global}>
      <p className={styles.date}>{goalData.date}まで</p>
      <p className={styles.text}>{goalData.text}</p>
    </div>
  );
}
