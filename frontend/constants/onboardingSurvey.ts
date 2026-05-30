import type { OnboardingParseStep } from '~/types/api'

export interface OnboardingVoiceQuestion {
  id: OnboardingParseStep
  title: string
  prompt: string
  placeholder: string
  examples: string[]
}

export const ONBOARDING_SURVEY_STEP_TO_PARSE: Record<number, OnboardingParseStep> = {
  1: 'income',
  2: 'cushion',
  3: 'goal',
  4: 'expenses'
}

export const ONBOARDING_VOICE_QUESTIONS: OnboardingVoiceQuestion[] = [
  {
    id: 'income',
    title: 'Доходы',
    prompt: 'Сколько получаете в месяц? Можно вслух — зарплата, пенсия, аренда.',
    placeholder: 'Зарплата 180 тысяч, аренда 25 тысяч',
    examples: ['180 тысяч в месяц', '150 000 зарплата и 20 000 с аренды']
  },
  {
    id: 'cushion',
    title: 'Запас',
    prompt: 'Сколько денег отложено «на чёрный день»? Или скажите «пока нет».',
    placeholder: 'Накопил 300 тысяч на вкладе',
    examples: ['300 тысяч на вкладе', 'пока нет запаса']
  },
  {
    id: 'goal',
    title: 'Цель',
    prompt: 'На что копите и сколько нужно? Отпуск, машина, ремонт — как удобно.',
    placeholder: 'Отпуск на 250 тысяч через год',
    examples: ['отпуск 250 тысяч', 'коплю на машину 1,2 млн']
  },
  {
    id: 'expenses',
    title: 'Постоянные платежи',
    prompt: 'Есть ли каждый месяц аренда, кредит, ипотека? Или скажите «пропустить».',
    placeholder: 'Аренда 45 тысяч, кредит 18',
    examples: ['аренда 45 тысяч', 'пропустить']
  }
]

export function getVoiceQuestionForSurveyStep(surveyStep: number): OnboardingVoiceQuestion | undefined {
  const id = ONBOARDING_SURVEY_STEP_TO_PARSE[surveyStep]
  if (!id) return undefined
  return ONBOARDING_VOICE_QUESTIONS.find((q) => q.id === id)
}
