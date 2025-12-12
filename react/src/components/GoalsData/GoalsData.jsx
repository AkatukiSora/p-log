import GoalData from "./GoalData"

export default function GoalsData({ goalDatas }) {

    const reverseGoalsDatas = [...goalDatas].reverse();

    // ユーザの目標だけをとれるようにしたい
    return (
        <>
            {reverseGoalsDatas.map((goalData, idx) => (
                <GoalData key={idx} goalData={goalData}/>
            ))}
        </>
    )
}
