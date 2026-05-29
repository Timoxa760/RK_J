export interface SankeyNode {
  name: string
  category?: 'income' | 'category' | 'savings'
}

export interface SankeyLink {
  source: string
  target: string
  value: number
}

export interface SankeyResponse {
  nodes: SankeyNode[]
  links: SankeyLink[]
}

export interface StoreBubble {
  id: string
  name: string
  avg_check: number
  visits: number
  total: number
  impulse_ratio: number
}

export interface StoresResponse {
  stores: StoreBubble[]
}

export interface CategoryItem {
  name: string
  amount: number
  share: number
}

export interface SubcategoryDetail {
  name: string
  items: { name: string; amount: number }[]
}

export interface CategoryDetail {
  name: string
  amount: number
  share: number
  subcategories: SubcategoryDetail[]
}

export interface CategoriesResponse {
  categories: CategoryDetail[]
}

export interface CompareMonth {
  month: string
  categories: CategoryItem[]
}

export interface CompareResponse {
  months: CompareMonth[]
}

export interface TimeMachinePoint {
  month: string
  actual: number
  optimistic: number
}

export interface TimeMachineResponse {
  points: TimeMachinePoint[]
  delta: number
}

export interface ForecastPoint {
  month: string
  amount: number
}

export interface ForecastResponse {
  points: ForecastPoint[]
}

export interface InsightItem {
  id: string
  title: string
  body: string
  severity: 'info' | 'warning' | 'success'
}

export interface InsightsResponse {
  insights: InsightItem[]
}

export interface CreditItem {
  id: string
  name: string
  balance: number
  payment: number
  rate: number
}

export interface CreditsDashboardResponse {
  dti: number
  stress_test_dti: number
  monthly_income: number
  credits: CreditItem[]
}

export interface ReceiptItem {
  id: string
  store: string
  amount: number
  date: string
  category: string
}

export interface ChallengeItem {
  id: string
  title: string
  participants: number
}

export interface LeaderboardEntry {
  id: string
  display_name: string
  relative_score: number
  rank: number
}

export const chartThemeLight = {
  backgroundColor: 'transparent',
  textStyle: { color: '#6d6760' },
  axisColor: '#f5dcc8',
  splitLine: '#f2ebe3'
}
