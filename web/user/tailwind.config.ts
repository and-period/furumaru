import type { Config } from 'tailwindcss'

export default <Partial<Config>>{
  theme: {
    extend: {
      animation: {
        'tracking-in-expand': 'tracking-in-expand 0.7s cubic-bezier(0.215, 0.610, 0.355, 1.000)   both',
      },
      keyframes: {
        'tracking-in-expand': {
          '0%': {
            'letter-spacing': '-.5em',
            'opacity': '0',
          },
          '40%': {
            opacity: '.6',
          },
          'to': {
            opacity: '1',
          },
        },
      },
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
        'base': '#F9F6EA',
        'main': '#604C3F',
        'typography': '#707070',
        'green': '#7CB342',
        'orange': '#F48D26',
        'apple-red': '#E74C3C',
        'facebook': '#1877F2',
        'line': '#06C755',
        'success': '#66bb6a',
        'error': '#f44336',
      },
      screens: {
        'sm': '640px',
        'md': '768px',
        'lg': '1024px',
        'xl': '1280px',
        '2xl': '1536px',
      },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      typography: (theme: any) => ({
        DEFAULT: {
          css: {
            'font-size': '14px',
            'line-height': '32px',
            'color': theme('colors.main'),
            '@media (min-width: theme("screens.md"))': {
              'font-size': '20px',
              'line-height': '48px',
            },
            'h1': {
              'font-size': '18px',
              'font-weight': 700,
              'letter-spacing': '0.1em',
              'margin-bottom': '40px',
              'line-height': '32px',
              'color': theme('colors.main'),
              '@media (min-width: theme("screens.md"))': {
                'font-size': '24px',
                'margin-bottom': '56px',
                'line-height': '48px',
              },
            },
            'h2': {
              'font-size': '16px',
              'font-weight': 700,
              'line-height': '24px',
              'color': theme('colors.main'),
              '@media (min-width: theme("screens.md"))': {
                'font-size': '22px',
                'line-height': '48px',
              },
            },
            'p': {
              'font-size': '14px',
              'line-height': '32px',
              'font-weight': 500,
              'color': theme('colors.main'),
              '@media (min-width: theme("screens.md"))': {
                'font-size': '20px',
                'line-height': '48px',
                'font-weight': 500,
              },
            },
            'a': {
              color: theme('colors.main'),
            },
            'li': {
              'font-size': '14px',
              'line-height': '32px',
              'font-weight': 500,
              'color': theme('colors.main'),
              '@media (min-width: theme("screens.md"))': {
                'font-size': '20px',
                'line-height': '48px',
                'font-weight': 500,
              },
            },
          },
        },
      }),
    },
  },
  plugins: [require('@tailwindcss/typography')],
}
