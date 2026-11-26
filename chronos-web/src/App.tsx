import { BrowserRouter, Route, Routes } from 'react-router'
import { ThemeProvider } from '@/contexts/theme'
import About from './pages/about'
import Chat from './pages/chat'
import Dashboard from './pages/dashboard'
import Home from './pages/home'
import Profile from './pages/profile'
import Recurring from './pages/recurring'
import SignIn from './pages/signin'
import Transactions from './pages/transactions'

function App() {
  return (
    <ThemeProvider>
    <BrowserRouter>
      <Routes>
        <Route index element={<Home />} />
        <Route path="/about" element={<About />} />
        <Route path="/chat" element={<Chat />} />
        <Route path="/signin" element={<SignIn />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/recurring" element={<Recurring />} />
        <Route path="/transactions" element={<Transactions />} />
        <Route path="/profile" element={<Profile />} />
      </Routes>
    </BrowserRouter>
    </ThemeProvider>
  )
}

export default App
