/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,jsx,ts,tsx}'],
  options: {
    // https://purgecss.com/safelisting.html#patterns
    safelist: {
      standard: [/^bg-/, /^text-/],
    },
  },
  theme: {
    extend: {
      colors: {
        gray_333: '#333333',
      },
    },
  },
  plugins: [],
}
