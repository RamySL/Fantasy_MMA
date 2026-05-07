/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,jsx}"],
  theme: {
    extend: {
      boxShadow: {
        soft: "0 18px 50px -30px rgba(15, 23, 42, 0.35)",
      },
    },
  },
  plugins: [],
};
