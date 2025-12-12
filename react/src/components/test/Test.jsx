import styles from "./Test.module.css"

export default function Test({name}) {
  return (
    <div className={styles.test}>
        {name}
    </div>
  )
}
