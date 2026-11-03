import { BrowserRouter, Route, Routes } from 'react-router'
import Chat from './pages/chat'
import Dashboard from './pages/dashboard'
import Home from './pages/home'
import SignIn from './pages/signin'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route index element={<Home />} />
        <Route path="/chat" element={<Chat />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/dashboard" element={<Dashboard />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
