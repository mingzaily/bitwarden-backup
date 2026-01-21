/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['"Inter"', 'system-ui', '-apple-system', 'sans-serif'],
      },
      colors: {
        brutalist: {
          cream: '#FAF7F0',
          charcoal: '#2A2A2A',
          blue: '#2563EB',
          green: '#10B981',
          orange: '#F97316',
          red: '#EF4444',
        }
      },
      boxShadow: {
        'brutalist': '4px 4px 0px 0px rgba(0, 0, 0, 1)',
        'brutalist-sm': '2px 2px 0px 0px rgba(0, 0, 0, 1)',
        'brutalist-hover': '2px 2px 0px 0px rgba(0, 0, 0, 1)',
      }
    },
  },
  plugins: [],
}
