import type { Config } from 'tailwindcss'

export default {
  content: [
    './app.vue',
    './layouts/**/*.vue',
    './pages/**/*.vue',
    './components/**/*.{vue,ts}'
  ],
  theme: {
    extend: {}
  },
  plugins: []
} satisfies Config
