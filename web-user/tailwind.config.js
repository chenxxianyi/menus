/** @type {import('tailwindcss').Config} */
export default {
  content: [
    './index.html',
    './src/**/*.{vue,js,ts,jsx,tsx}',
  ],
  theme: {
    extend: {
      colors: {
        'glass-border': 'rgba(255,255,255,0.5)',
        'glass-bg': 'rgba(255,255,255,0.65)',
        'glass-bg-light': 'rgba(255,255,255,0.22)',
        navy: '#1e1e2e',
        orange: { DEFAULT: '#ffb347', dark: '#f59e0b' },
        purple: { DEFAULT: '#b8a9e8', dark: '#8b7bc6' },
        blue: { DEFAULT: '#a8d8ea' },
        green: { DEFAULT: '#b8e6c8' },
        pink: { DEFAULT: '#f5d0f0' },
        yellow: { DEFAULT: '#fff3cd' },
        red: { DEFAULT: '#fca5a5' },
      },
      fontFamily: {
        sans: ['Plus Jakarta Sans', '-apple-system', 'system-ui', 'sans-serif'],
      },
      borderRadius: {
        'glass': '20px',
        'glass-lg': '24px',
        'glass-xl': '28px',
      },
      boxShadow: {
        'glass': '0 8px 32px rgba(0,0,0,0.06)',
        'glass-sm': '0 4px 16px rgba(0,0,0,0.05)',
        'glass-lg': '0 8px 40px rgba(0,0,0,0.08)',
        'glass-highlight': 'inset 0 1px 0 rgba(255,255,255,0.8)',
        'glass-highlight-sm': 'inset 0 1px 0 rgba(255,255,255,0.6)',
      },
      backdropBlur: {
        'glass': '20px',
        'glass-sm': '16px',
        'glass-lg': '24px',
      },
      animation: {
        'fade-up': 'fadeUp 0.5s cubic-bezier(0.25,0.46,0.45,0.94) forwards',
      },
      keyframes: {
        fadeUp: {
          '0%': { opacity: '0', transform: 'translateY(16px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
      },
    },
  },
  plugins: [],
}
