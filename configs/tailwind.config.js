const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: {
    files: [
      'internal/web/**/*.templ',
    ]
  },
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["cupcake"],
          ".btn-twitter": {
            "background-color": "#1EA1F1",
            "border-color": "#1EA1F1",
          },
          ".login-button": {
            "padding-top": "0.25rem", /* Adjust as needed */
            "padding-bottom": "0.25rem", /* Adjust as needed */
            "height": "2rem",
          },
          ".text-xxl": {
            "font-size": "2.00rem",
            "line-height": "1.75rem",
            "font-weight": "700",
          },
          ".tab": {
            "border-color": "white",
          },
        },
      },
    ],
  },
  plugins: [
    require('@tailwindcss/typography'),

    require("daisyui"),
  ]
}