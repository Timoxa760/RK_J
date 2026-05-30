/**
 * Палитра бренда — оттенки одного hue (#96C8F9, заливка HeroFlowWord).
 * CSS: :root в assets/css/main.css. Менять здесь и там вместе.
 */
export const brandPalette = {
  primary: '#96C8F9',
  primaryHover: '#72B5F6',
  primaryLight: '#B8DBFA',
  primarySoft: '#EBF5FE',
  primaryMuted: '#C5E3FC',
  primaryDeep: '#4A9DE8',
  bg: '#F5FAFF',
  bgElevated: '#FFFFFF',
  bgMuted: '#EAF4FD',
  border: '#D0E8FA',
  borderSubtle: '#E3F2FC',
  text: '#1C1C1A',
  textMuted: '#5A6B78',
  textSoft: '#8496A3',
  flowTrack: '#D8EDFA',
  flowAccent: '#A8D4F7',
  flowPulse: '#6EB5F0'
} as const

export const brandPaletteRgb = {
  primary: '150 200 249',
  primarySoft: '235 245 254',
  primaryMuted: '197 227 252'
} as const

/** HSL-части для Tailwind/shadcn: hsl(var(--token)). */
export const brandHsl = {
  background: '207 55% 98%',
  foreground: '60 4% 11%',
  card: '0 0% 100%',
  cardForeground: '60 4% 11%',
  popover: '0 0% 100%',
  popoverForeground: '60 4% 11%',
  primary: '207 88% 78%',
  primaryForeground: '207 50% 20%',
  secondary: '207 48% 95%',
  secondaryForeground: '60 4% 11%',
  muted: '207 48% 94%',
  mutedForeground: '207 14% 42%',
  accent: '207 62% 92%',
  accentForeground: '207 42% 30%',
  destructive: '0 72% 51%',
  destructiveForeground: '0 0% 100%',
  border: '207 52% 87%',
  input: '207 52% 87%',
  ring: '207 78% 72%',
  chart1: '207 88% 78%',
  chart2: '207 72% 62%',
  chart3: '207 58% 48%',
  chart4: '200 62% 68%',
  chart5: '214 55% 54%',
  sidebarBackground: '207 55% 98%',
  sidebarForeground: '60 4% 11%',
  sidebarPrimary: '207 88% 78%',
  sidebarPrimaryForeground: '207 50% 20%',
  sidebarAccent: '207 62% 92%',
  sidebarAccentForeground: '207 42% 30%',
  sidebarBorder: '207 52% 87%',
  sidebarRing: '207 78% 72%'
} as const
