/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["../../templates/*.{html,js,tmpl}", "./**/*.{html,js}" ],
  theme: {
    extend: {
      minWidth: {
        '300': '300px',
        '500': '500px',
        '700': '700px'
      },
      minHeight: {
        '300': '300px',
        '500': '500px',
        '700': '700px'
      },
      fontFamily: {
        sans: ['Montserrat', 'sans-serif'],
      },
      colors: {
        bblue: {
          '500': '#177CB1',
        },
        gred: {
          '500': '#DB4437'
        }
      },
    },
  },
  plugins: [],
}

