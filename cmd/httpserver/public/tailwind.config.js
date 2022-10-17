const colors = require("tailwindcss/colors")

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{vue,html,js,ts,md}"],
  theme: {
    colors: {
      transparent: "transparent",
      current: "currentColor",
      black: colors.black,
      white: colors.white,
      channel: {
        50: "#f5f5f5",
        100: "#eaeaeb",
        200: "#cbcccd",
        300: "#acadaf",
        400: "#6d6f72",
        500: "#2f3136",
        600: "#2a2c31",
        700: "#232529",
        800: "#1c1d20",
        900: "#17181a",
      },
      chat: {
        200: "#cdcecf",
        300: "#afb0b2",
        400: "#727479",
        500: "#36393f",
        600: "#313339",
        700: "#292b2f",
        800: "#202226",
        900: "#1a1c1f",
      },
      gray: {
        50: "#fafbfb",
        100: "#f5f7f8",
        200: "#e6eaed",
        300: "#d6dde1",
        400: "#b8c4cb",
        500: "#99aab5",
        600: "#8a99a3",
        700: "#738088",
        800: "#5c666d",
        900: "#4b5359",
      },
      skin: {
        500: "#f9c9a9",
        600: "#e0b598",
        700: "#bb977f",
        800: "#957965",
        900: "#7a6253",
      },
      nitro: {
        300: "#ffc7fd",
        400: "#ff9dfc",
        500: "#ff73fa",
        600: "#e668e1",
        700: "#bf56bc",
        800: "#994596",
        900: "#7d387b",
      },
      high: {
        300: "#fbc9ad",
        400: "#f8a06f",
        500: "#f57731",
        600: "#dd6b2c",
        700: "#b85925",
        800: "#93471d",
        900: "#783a18",
      },
      idle: {
        300: "#fddba3",
        400: "#fcc15f",
        500: "#faa61a",
        600: "#e19517",
        700: "#bc7d14",
        800: "#966410",
        900: "#7b510d",
      },
      brilliance: {
        300: "#fbcac3",
        400: "#f7a395",
        500: "#f47b68",
        600: "#dc6f5e",
        700: "#b75c4e",
        800: "#924a3e",
        900: "#783c33",
      },
      dnd: {
        300: "#f9b5b5",
        400: "#f57e7e",
        500: "#f04747",
        600: "#d84040",
        700: "#b43535",
        800: "#902b2b",
        900: "#762323",
      },
      balance: {
        300: "#b4f1e5",
        400: "#7ce7d2",
        500: "#44ddbf",
        600: "#3dc7ac",
        700: "#33a68f",
        800: "#298573",
        900: "#216c5e",
      },
      online: {
        300: "#aee1cd",
        400: "#71cba7",
        500: "#34b581",
        600: "#2fa374",
        700: "#278861",
        800: "#1f6d4d",
        900: "#19593f",
      },
      bravery: {
        300: "#d7cef8",
        400: "#b9a9f3",
        500: "#9b84ee",
        600: "#8c77d6",
        700: "#7463b3",
        800: "#5d4f8f",
        900: "#4c4175",
      },
      discord: {
        300: "#c7d0f0",
        400: "#9cace5",
        500: "#7289da",
        600: "#677bc4",
        700: "#5667a4",
        800: "#445283",
        900: "#38436b",
      },
    },
    extend: {},
  },
  plugins: [require("daisyui"), require("@tailwindcss/forms"), require("@tailwindcss/typography")],
  daisyui: {
    themes: [
      {
        default: {
          primary: "#9b84ee",
          "primary-focus": "#8c77d6",
          secondary: "#e668e1",
          "secondary-focus": "#bf56bc",
          accent: "#3dc7ac",
          "accent-focus": "#33a68f",
          neutral: "#1c1d20",
          "neutral-focus": "#17181a",
          "base-100": "#232529",
          "base-200": "#1c1d20",
          "base-300": "#17181a",
          info: "#7289da",
          "info-focus": "#677bc4",
          success: "#2fa374",
          "success-focus": "#278861",
          warning: "#e19517",
          "warning-focus": "#bc7d14",
          error: "#b43535",
          "error-focus": "#902b2b",

          "--rounded-box": "0.25rem",
          "--rounded-btn": "0.25rem",
          "--rounded-badge": "0.25rem",
          "--tab-radius": "0.25rem",
        },
      },
    ],
    logs: true,
  },
  safelist: [
    "flex items-center gap-2", // guild icon in tables
    "text-bravery-500 text-idle-500 text-dnd-500 text-nitro-500", // status in tables
    "flex-col badge badge-secondary", // metadata in tables
    "prose max-w-none", // markdown
  ],
}
