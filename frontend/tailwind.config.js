/** @type {import('tailwindcss').Config} */
module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  content: [],
  theme: {
    screens: {
      xs: '400px',
      sm: '540px',
      md: '768px',
      lg: '1300px'
    },
    extend: {
      backgroundImage: {
        'JNR185': "url('/assets/img/background3.jpg')"
      }
    },
  },
  plugins: [],
}

