import { createContext, useState, useEffect } from 'react'
import type { Theme, ThemeContextValue } from '@/types'

export const ThemeContext = createContext<ThemeContextValue | null>(null)

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem('theme') as Theme) ?? 'light'
  )

  useEffect(() => {
    document.documentElement.classList.toggle('dark', theme === 'dark')
    localStorage.setItem('theme', theme)
  }, [theme])

  function setTheme(theme: Theme) {
    setThemeState(theme)
  }

  return (
    <ThemeContext value={{ theme, setTheme }}>
      {children}
    </ThemeContext>
  )
}
