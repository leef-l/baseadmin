/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#ff6a00',
        primaryDark: '#e55c00',
        accent: '#ffb74d',
      },
    },
  },
  corePlugins: {
    preflight: false,
  },
  plugins: [],
};
