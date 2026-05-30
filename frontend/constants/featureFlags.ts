/** Временные флаги продукта — включить обратно, когда фича готова. */
export const APP_FEATURES = {
  creditsNav: true
} as const

export type AppFeatureKey = keyof typeof APP_FEATURES

export function isAppFeatureEnabled(key: AppFeatureKey) {
  return APP_FEATURES[key]
}
