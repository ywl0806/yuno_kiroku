import colors from './src/colors'

export default {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  safelist: [{ pattern: /^bg-/ }, { pattern: /^text-/ }],
  theme: {
    extend: {
      colors,
    },
  },
  plugins: [],
}
