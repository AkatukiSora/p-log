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

function App() {

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path={ROUTES.myList} element={<MyGoalsList />}/>
          <Route path={ROUTES.myProfile} element={<MyProfile />}/>
          <Route path={ROUTES.friendList} element={<FriendList />}/>
          <Route path={ROUTES.friendProfile} element={<FriendProfile />}/>
          <Route path={ROUTES.friendGoals} element={<FriendGoalsList />}/>
          <Route path={ROUTES.timeLine} element={<TimeLine />}/>
          <Route path={ROUTES.settingGoal} element={<SettingGoal />}/>
          <Route path={ROUTES.addPost} element={<AddPost />}/>
          <Route path={ROUTES.goals} element={<Goals />}/>
        </Routes>
      </BrowserRouter>
    </>
  )
}

export default App
