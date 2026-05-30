// API types aligned with docs/api/API_Contract.md; UI helpers in utils/apiNormalize.ts

export interface AuthUser {
  id: string
  phone: string
  role?: string
  name?: string
}

export interface LoginResponse {
  access_token: string
  refresh_token?: string
  expires_in?: number
  user: AuthUser
}

export interface SankeyNode {
  name: string
  value?: number
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
  id?: string
  name: string
  avg_check: number
  visits: number
  purchases?: number
  total: number
  impulse_ratio: number
}

export interface StoresResponse {
  stores: StoreBubble[]
}

export interface CategoryLineItem {
  name: string
  price?: number
  quantity?: number
  total?: number
  amount?: number
}

export interface SubcategoryDetail {
  name: string
  total?: number
  items: CategoryLineItem[]
}

export interface CategoryDetail {
  name: string
  amount: number
  share?: number
  total?: number
  subcategories: SubcategoryDetail[]
}

export interface CategoriesResponse {
  categories: CategoryDetail[]
}

export interface CompareCategory {
  name: string
  amount?: number
  total?: number
  share?: number
}

export interface CompareMonth {
  month?: string
  label?: string
  categories: CompareCategory[]
}

export interface CompareInsights {
  biggest_change?: {
    category: string
    delta: number
    delta_percent: number
  }
}

export interface CompareResponse {
  months: CompareMonth[]
  insights?: CompareInsights
}

export interface TimeMachinePoint {
  month: string
  actual: number
  optimistic: number
}

/** Сырой ответ API GET /dashboard/timemachine и POST /scenarios/simulate */
export interface TimeMachineApiResponse {
  months: string[]
  real_savings: number[]
  optimized_savings: number[]
  difference_final?: number
}

export interface TimeMachineResponse {
  points: TimeMachinePoint[]
  delta: number
  difference_final?: number
}

export interface ForecastResponse {
  dates: string[]
  forecast: number[]
  upper_bound?: number[]
  lower_bound?: number[]
  /** @deprecated legacy mock shape */
  points?: { month: string; amount: number }[]
}

export type InsightSeverity = 'info' | 'warning' | 'success' | 'critical'

export interface InsightItem {
  id?: string
  type?: string
  title: string
  body?: string
  description?: string
  severity: InsightSeverity
  amount?: number
  merchant?: string
  store?: string
}

export interface InsightsResponse {
  insights: InsightItem[]
}

export type AiDiagnosisIndicatorStatus = 'good' | 'warning' | 'critical'

export type AiDiagnosisDifficulty = 'easy' | 'medium' | 'hard'

export interface AiDiagnosisIndicator {
  name: string
  value: number
  norm: string
  status: AiDiagnosisIndicatorStatus
}

export interface AiDiagnosisMainAction {
  title: string
  description: string
  potential_savings: number
  difficulty: AiDiagnosisDifficulty
}

export interface AiDiagnosisResponse {
  score: number
  grade: string
  indicators: AiDiagnosisIndicator[]
  main_action: AiDiagnosisMainAction
  next_check_days: number
}

export interface SimulateScenarioRequest {
  scenario: 'reduce_delivery' | 'reduce_cafe' | 'reduce_entertainment' | 'custom'
  reduction_percent: number
  months?: number
}

export interface SimulateScenarioResponse {
  months: string[]
  real_savings: number[]
  optimized_savings: number[]
  difference_final: number
  scenario?: {
    name: string
    monthly_saving: number
    annual_saving: number
  }
}

export interface CreditItem {
  id: string
  name?: string
  bank?: string
  balance?: number
  remaining?: number
  payment?: number
  monthly_payment?: number
  rate: number
  term_months?: number
  amount?: number
  next_payment?: string
}

export interface CreditsDashboardResponse {
  dti: number
  stress_test_dti?: number
  stress_test_months?: number
  savings?: number
  total_debt?: number
  monthly_payments?: number
  monthly_income: number
  credits: CreditItem[]
}

export interface CreditScanResponse {
  parsed: {
    amount: number
    rate: number
    term_months: number
    monthly_payment: number
    bank: string
  }
  confidence: number
}

export interface ReceiptLineItem {
  name: string
  price: number
  quantity?: number
  category?: string
  impulse?: boolean
}

export interface ReceiptListItem {
  id: string
  store: string
  amount: number
  date: string
  category?: string
  items?: ReceiptLineItem[]
  impulse_count?: number
}

export interface ReceiptsListResponse {
  receipts: ReceiptListItem[]
}

export interface ManualExpenseRequest {
  user_id: string
  raw_text?: string
  amount?: number
  category?: string
  description?: string
  date?: string
  source?: 'manual' | 'voice'
}

export interface ManualExpenseResponse {
  success: boolean
  id: string
  amount: number
  category: string
  parsed: boolean
}

export interface ReceiptManualRequest {
  store: string
  amount: number
  category: string
  date: string
}

export interface ReceiptManualResponse {
  receipt_id: string
  store: string
  amount: number
  category: string
  date: string
  status: string
}

export interface ReceiptVoiceResponse {
  receipt_id: string
  store: string
  items: ReceiptLineItem[]
  total: number
  category: string
  confidence: number
}

export interface ReceiptFnsScanRequest {
  fn: string
  fd: string
  fp: string
}

export interface ReceiptFnsScanResponse {
  receipt_id: string
  store: string
  inn?: string
  date?: string
  total: number
  items?: ReceiptLineItem[]
  category: string
}

export interface ProfilePatchRequest {
  active_income?: number
  passive_income?: number
  emergency_fund?: number
  emergency_breakdown?: EmergencyFundBreakdown
  fixed_expenses?: FixedExpense[]
}

export interface OnboardingCompleteResponse {
  onboarding_completed: boolean
}

export interface FnsTicketRequest {
  qr: string
}

export interface FnsTicketResponse {
  success?: boolean
  message?: string
  receipt_id?: string
}

export interface FnsMcoSyncResponse {
  synced: number
  message?: string
}

export type ProviderId =
  | 'x5club'
  | 'magnit'
  | 'lenta'
  | 'vkusvill'
  | 'ozon'
  | 'wb'
  | 'email'
  | 'fns'

export interface ProviderInfo {
  id: ProviderId
  label: string
}

export interface ProviderConnectResponse {
  message: string
  provider: string
  status: string
}

export type ChallengeType = 'least_spend' | 'most_saved' | 'streak'

export interface ChallengeItem {
  id: string
  type?: ChallengeType
  title: string
  status?: string
  participants: number
  invite_token?: string
  duration_days?: number
  created_at?: string
}

export interface LeaderboardEntry {
  position: number
  username: string
  avatar?: string | null
  relative_score: number
  /** legacy */
  id?: string
  display_name?: string
  rank?: number
}

export interface LeaderboardResponse {
  challenge_id: string
  type?: ChallengeType
  leaderboard: LeaderboardEntry[]
  my_position?: { position: number; total_participants: number }
}

export interface DigestPeriod {
  from: string
  to: string
}

export interface DigestCategoryRow {
  name: string
  total: number
  percent: number
  trend?: string
}

export interface DigestStoreRow {
  name: string
  total: number
  visits: number
}

export interface DigestResponse {
  period: DigestPeriod
  total_spent: number
  total_income: number
  saved: number
  by_category: DigestCategoryRow[]
  word_cloud: string[]
  top_stores: DigestStoreRow[]
  mindfulness_rating: number
  ai_advice: string
  insights_summary: string
}

export interface Goal {
  id: string
  title: string
  target_amount: number
  current_amount: number
  progress_percent: number
  target_date?: string
  auto_save_percent?: number
}

export interface GoalsListResponse {
  goals: Goal[]
}

export interface CreateGoalRequest {
  title: string
  target_amount: number
  target_date?: string
  auto_save_percent?: number
}

export interface FinancialProfile {
  active_income: number
  passive_income: number
  emergency_fund: number
  emergency_breakdown?: EmergencyFundBreakdown
  fixed_expenses?: FixedExpense[]
  onboarding_completed?: boolean
}

export type GoalKind = 'save' | 'purchase' | 'cushion' | 'other'

export interface EmergencyFundBreakdown {
  cash: number
  deposit: number
  investments: number
}

export interface FixedExpense {
  title: string
  amount: number
}

export interface OnboardingDraft {
  active_income: number
  passive_income: number
  emergency_fund: number
  emergency_breakdown: EmergencyFundBreakdown
  goal_kind: GoalKind
  goal_title: string
  goal_amount: number
  fixed_expenses: FixedExpense[]
  skipped_expenses: boolean
}

/** Шаг опроса онбординга — совпадает с id в `ONBOARDING_VOICE_QUESTIONS`. */
export type OnboardingParseStep = 'income' | 'cushion' | 'goal' | 'expenses'

/** @deprecated API v3 — парсер опроса только локально на front. */
export interface OnboardingParseRequest {
  step: OnboardingParseStep
  raw_text: string
  locale?: string
}

/** @deprecated API v3 — парсер опроса только локально на front. */
export interface OnboardingParseResponse {
  parsed: boolean
  step: OnboardingParseStep
  patch: Partial<OnboardingDraft>
  confidence?: number
  message?: string
}

export type MortgageApprovalLevel = 'high' | 'medium' | 'low'

export interface BankOffer {
  id: string
  bank: string
  rate: number
  monthly_payment: number
  total_overpayment: number
  term_months: number
}

export interface MortgageBreakdownResponse {
  approval_level: MortgageApprovalLevel
  approval_reason: string
  safe_mortgage_amount: number
  comfortable_payment: number
  load_risk: string
  scenario_now: string
  scenario_wait: string
  wait_months: number
  banks: BankOffer[]
  optimal_bank_id: string
}

export interface MortgageAnalyzeRequest {
  mortgage_amount: number
  monthly_income?: number
  savings?: number
  existing_dti?: number
  stress_test_months?: number
}

export { chartThemeLight } from '~/utils/chartTheme'

export const PROVIDERS: ProviderInfo[] = [
  { id: 'x5club', label: 'X5 Club' },
  { id: 'magnit', label: 'Магнит' },
  { id: 'lenta', label: 'Лента' },
  { id: 'vkusvill', label: 'ВкусВилл' },
  { id: 'ozon', label: 'Ozon' },
  { id: 'wb', label: 'Wildberries' }
]

export const SCENARIO_OPTIONS = [
  { value: 'reduce_delivery', label: 'Доставка' },
  { value: 'reduce_cafe', label: 'Кафе' },
  { value: 'reduce_entertainment', label: 'Развлечения' },
  { value: 'custom', label: 'Своя категория' }
] as const
