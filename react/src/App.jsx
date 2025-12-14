import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ROUTES } from './pages/pathData/path.js'
import MyGoalsList from './pages/myGoalsList/MyGoalsList.jsx'
import MyProfile from './pages/myProfile/MyProfile.jsx'
import FriendList from './pages/friendList/FriendList.jsx'
import FriendProfile from './pages/friendProfile/FriendProfile.jsx'
import FriendGoalsList from './pages/friendGoalsList/FriendGoalsList.jsx'
import TimeLine from './pages/timeLine/TimeLine.jsx'
import AddPost from './pages/AddPost/AddPost.jsx'
import SettingGoal from './pages/settingGoal/SettingGoal.jsx'
import Goals from './pages/goals/Goals.jsx'
import { useState } from 'react'

function App() {
  const [users,setUsers]=useState();

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path={ROUTES.myList} element={<MyGoalsList users={users}/>}/>
          <Route path={ROUTES.myProfile} element={<MyProfile users={users}/>}/>
          <Route path={ROUTES.friendList} element={<FriendList users={users}/>}/>
          <Route path={ROUTES.friendProfile} element={<FriendProfile users={users}/>}/>
          <Route path={ROUTES.friendGoals} element={<FriendGoalsList users={users}/>}/>
          <Route path={ROUTES.timeLine} element={<TimeLine users={users}/>}/>
          <Route path={ROUTES.settingGoal} element={<SettingGoal users={users}/>}/>
          <Route path={ROUTES.addPost} element={<AddPost users={users}/>}/>
          {/* <Route path="/"
          {ROUTES.goals}
           element={<Goals />}/> */}
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
