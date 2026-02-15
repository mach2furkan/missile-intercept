/** @type {import('tailwindcss').Config} */
export default {
    content: [
        "./index.html",
        "./src/**/*.{js,ts,jsx,tsx}",
    ],
    theme: {
        extend: {
            colors: {
                radar: {
                    bg: '#001a00',
                    grid: '#003300',
                    beam: '#00ff00',
                    blip: '#ccffcc',
                    text: '#00cc00',
                    alert: '#ff0000',
                }
            },
            fontFamily: {
                mono: ['"Courier New"', 'Courier', 'monospace'],
            },
            animation: {
                'radar-sweep': 'radar-sweep 4s linear infinite',
                'blink': 'blink 1s cubic-bezier(0.4, 0, 0.6, 1) infinite',
            },
            keyframes: {
                'radar-sweep': {
                    '0%': { transform: 'rotate(0deg)' },
                    '100%': { transform: 'rotate(360deg)' },
                },
                blink: {
                    '0%, 100%': { opacity: 1 },
                    '50%': { opacity: .5 },
                }
            }
        },
    },
    plugins: [],
}
