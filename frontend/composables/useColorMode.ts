export function useColorMode() {
  const isDark = useState('color-mode-dark', () => false)

  function apply(dark: boolean) {
    if (!import.meta.client) {
      return
    }
    document.documentElement.classList.toggle('dark', dark)
    localStorage.setItem('demonit-theme', dark ? 'dark' : 'light')
  }

  function toggle() {
    isDark.value = !isDark.value
    apply(isDark.value)
  }

  function init() {
    if (!import.meta.client) {
      return
    }
    const stored = localStorage.getItem('demonit-theme')
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    isDark.value = stored ? stored === 'dark' : prefersDark
    apply(isDark.value)
  }

  return { isDark, toggle, init }
}
