import GoalData from "./GoalData"

export default function GoalsData({ goalsData }) {
    return (
        <>
            {goalsData.map((goalData, idx) => (
                <GoalData key={idx} goalData={goalData}/>
            ))}
        </>
    )
}
