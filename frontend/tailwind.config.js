/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: 'class',
  content: [
    './components/**/*.{vue,js,ts}',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './composables/**/*.{js,ts}',
    './plugins/**/*.{js,ts}',
    './app.vue',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Outfit"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['"IBM Plex Mono"', 'ui-monospace', 'monospace'],
        display: ['"Outfit"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        ink: 'rgb(var(--ink) / <alpha-value>)',
        paper: 'rgb(var(--paper) / <alpha-value>)',
        surface: 'rgb(var(--surface) / <alpha-value>)',
        // Google palette — stay vivid in both themes
        gblue: '#4285F4',
        gred: '#EA4335',
        gyellow: '#FBBC05',
        ggreen: '#34A853',
        lime: '#34A853',
        pink: '#EA4335',
        cyan: '#4285F4',
        sun: '#FBBC05',
        online: {
          DEFAULT: 'rgb(var(--ink) / <alpha-value>)',
          soft: '#34A853',
        },
        offline: {
          DEFAULT: 'rgb(var(--ink) / <alpha-value>)',
          soft: '#EA4335',
        },
        alert: {
          DEFAULT: 'rgb(var(--ink) / <alpha-value>)',
          soft: '#EA4335',
        },
        accent: {
          DEFAULT: '#4285F4',
          soft: 'rgb(var(--accent-soft) / <alpha-value>)',
          strong: 'rgb(var(--ink) / <alpha-value>)',
        },
      },
      borderRadius: {
        brutal: '1rem',
        'brutal-lg': '1.5rem',
      },
      boxShadow: {
        brutal: '4px 4px 0 0 rgb(var(--ink))',
        'brutal-sm': '3px 3px 0 0 rgb(var(--ink))',
        'brutal-lg': '6px 6px 0 0 rgb(var(--ink))',
        'brutal-blue': '4px 4px 0 0 #4285F4',
        'brutal-red': '4px 4px 0 0 #EA4335',
        'brutal-yellow': '4px 4px 0 0 #FBBC05',
        'brutal-green': '4px 4px 0 0 #34A853',
      },
      borderWidth: {
        3: '3px',
      },
      keyframes: {
        'slide-down': {
          from: { opacity: '0', transform: 'translateY(-0.75rem)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'pulse-dot': {
          '0%': { opacity: '1' },
          '50%': { opacity: '0.45' },
          '100%': { opacity: '1' },
        },
      },
      animation: {
        'slide-down': 'slide-down 0.2s ease-out',
        'fade-in': 'fade-in 0.25s ease-out',
        'pulse-dot': 'pulse-dot 1.4s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
