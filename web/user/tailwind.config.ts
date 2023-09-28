import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  theme: {
    extend: {
      inset: {
        '9px': '9px',
      },
      fontFamily: {
        sans: ['Noto Sans JP', 'sans-serif'],
      },
      dropShadow: {
        sm: '2px 2px 8px rgba(0, 0, 0, 0.10)',
      },
      colors: {
        base: '#F9F6EA',
        main: '#604C3F',
        typography: '#707070',
        green: '#7CB342',
        orange: '#F48D26',
        'apple-red': '#E74C3C',
        facebook: '#1877F2',
        line: '#06C755',
      },
      screens: {
        sm: '640px',
        md: '768px',
        lg: '1024px',
        xl: '1280px',
        '2xl': '1536px',
      },
    },
  },
}
