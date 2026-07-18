/** @type {import('tailwindcss').Config} */
module.exports = {
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
        ink: '#202124',
        paper: '#F8F9FA',
        // Google palette
        gblue: '#4285F4',
        gred: '#EA4335',
        gyellow: '#FBBC05',
        ggreen: '#34A853',
        lime: '#34A853',
        pink: '#EA4335',
        cyan: '#4285F4',
        sun: '#FBBC05',
        online: {
          DEFAULT: '#202124',
          soft: '#34A853',
        },
        offline: {
          DEFAULT: '#202124',
          soft: '#EA4335',
        },
        alert: {
          DEFAULT: '#202124',
          soft: '#EA4335',
        },
        accent: {
          DEFAULT: '#4285F4',
          soft: '#D2E3FC',
          strong: '#202124',
        },
        surface: {
          50: '#F8F9FA',
          100: '#F1F3F4',
          200: '#202124',
          800: '#202124',
          900: '#202124',
          950: '#202124',
        },
      },
      borderRadius: {
        brutal: '1rem',
        'brutal-lg': '1.5rem',
      },
      boxShadow: {
        brutal: '4px 4px 0 0 #202124',
        'brutal-sm': '3px 3px 0 0 #202124',
        'brutal-lg': '6px 6px 0 0 #202124',
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
