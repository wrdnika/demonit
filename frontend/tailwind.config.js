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
        sans: ['"DM Sans"', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['"IBM Plex Mono"', 'ui-monospace', 'monospace'],
      },
      colors: {
        surface: {
          50: '#f4f7f6',
          100: '#e8eeec',
          200: '#d1ddd9',
          800: '#1e2a27',
          900: '#121a18',
          950: '#0b110f',
        },
        accent: {
          DEFAULT: '#0d9488',
          soft: '#ccfbf1',
          strong: '#0f766e',
        },
        alert: {
          DEFAULT: '#dc2626',
          soft: '#fee2e2',
        },
        online: {
          DEFAULT: '#059669',
          soft: '#d1fae5',
        },
        offline: {
          DEFAULT: '#dc2626',
          soft: '#fee2e2',
        },
      },
      boxShadow: {
        panel: '0 1px 2px rgb(18 26 24 / 0.06), 0 8px 24px rgb(18 26 24 / 0.06)',
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
        'slide-down': 'slide-down 0.25s ease-out',
        'fade-in': 'fade-in 0.3s ease-out',
        'pulse-dot': 'pulse-dot 1.6s ease-in-out infinite',
      },
    },
  },
  plugins: [],
}
