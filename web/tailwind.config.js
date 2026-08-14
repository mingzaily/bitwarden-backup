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
        ink: 'rgb(var(--color-bg) / <alpha-value>)',
        panel: 'rgb(var(--color-surface) / <alpha-value>)',
        'panel-soft': 'rgb(var(--color-surface-raised) / <alpha-value>)',
        line: 'rgb(var(--color-border) / <alpha-value>)',
        accent: 'rgb(var(--color-accent) / <alpha-value>)',
        violet: '#8B5CF6',
        cyan: '#22D3EE',
        danger: 'rgb(var(--color-danger) / <alpha-value>)',
        brutalist: {
          cream: 'rgb(var(--color-bg) / <alpha-value>)',
          charcoal: 'rgb(var(--color-text) / <alpha-value>)',
          blue: 'rgb(var(--color-accent) / <alpha-value>)',
          green: 'rgb(var(--color-accent) / <alpha-value>)',
          orange: 'rgb(var(--color-warning) / <alpha-value>)',
          red: 'rgb(var(--color-danger) / <alpha-value>)',
        }
      },
      boxShadow: {
        card: '0 18px 45px rgba(0, 0, 0, 0.18)',
        glow: '0 0 30px rgba(30, 215, 96, 0.18)',
        'brutalist': '0 14px 34px rgba(0, 0, 0, 0.20)',
        'brutalist-sm': '0 10px 24px rgba(0, 0, 0, 0.18)',
        'brutalist-hover': '0 16px 38px rgba(0, 0, 0, 0.24)',
      }
    },
  },
  plugins: [],
}
