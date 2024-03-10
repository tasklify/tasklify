const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: [
      'internal/web/**/*.templ',
    ]
  },
  theme: {
    extend: {
      colors: {
        primary: colors.blue,
        secondary: colors.yellow,
        neutral: colors.gray,
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),

    require("daisyui"),
  ]
}