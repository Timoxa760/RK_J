<script setup lang="ts">
import { Check, Circle, Loader2 } from 'lucide-vue-next'
import {
  FINANCIAL_REPORT_LOADING_COPY,
  FINANCIAL_REPORT_LOADING_STAGES,
  type FinancialReportLoadingStage
} from '~/constants/financialReportLoading'

const props = withDefaults(
  defineProps<{
    active: boolean
    stages?: FinancialReportLoadingStage[]
    title?: string
    subtitle?: string
    compact?: boolean
  }>(),
  {
    stages: () => FINANCIAL_REPORT_LOADING_STAGES,
    title: FINANCIAL_REPORT_LOADING_COPY.title,
    subtitle: FINANCIAL_REPORT_LOADING_COPY.subtitle
  }
)

const activeRef = toRef(props, 'active')
const { currentIndex, currentStage, progress, isLastStage } = useStagedLoading(
  props.stages,
  activeRef
)

function stageState(index: number): 'done' | 'active' | 'pending' {
  if (index < currentIndex.value) return 'done'
  if (index === currentIndex.value) return 'active'
  return 'pending'
}
</script>

<template>
  <div
    class="mm-financial-report-loading"
    :class="{ 'mm-financial-report-loading--compact': compact }"
    role="status"
    aria-live="polite"
    :aria-busy="active"
  >
    <div class="mm-financial-report-loading__head">
      <p class="mm-financial-report-loading__title">{{ title }}</p>
      <p class="mm-financial-report-loading__subtitle">
        {{ isLastStage && active ? FINANCIAL_REPORT_LOADING_COPY.waitMore : subtitle }}
      </p>
    </div>

    <div class="mm-financial-report-loading__progress" aria-hidden="true">
      <div class="mm-financial-report-loading__progress-track">
        <div
          class="mm-financial-report-loading__progress-fill"
          :style="{ width: `${progress}%` }"
        />
      </div>
      <span class="mm-financial-report-loading__progress-label">{{ progress }}%</span>
    </div>

    <p class="mm-financial-report-loading__current">
      {{ currentStage.label }}
      <span v-if="currentStage.hint" class="mm-financial-report-loading__current-hint">
        — {{ currentStage.hint }}
      </span>
    </p>

    <ol class="mm-financial-report-loading__steps">
      <li
        v-for="(stage, index) in stages"
        :key="stage.id"
        class="mm-financial-report-loading__step"
        :class="`mm-financial-report-loading__step--${stageState(index)}`"
      >
        <span class="mm-financial-report-loading__step-icon" aria-hidden="true">
          <Check v-if="stageState(index) === 'done'" class="size-3.5" />
          <Loader2 v-else-if="stageState(index) === 'active'" class="size-3.5 animate-spin" />
          <Circle v-else class="size-2.5 fill-current stroke-none opacity-40" />
        </span>
        <span class="mm-financial-report-loading__step-copy">
          <span class="mm-financial-report-loading__step-label">{{ stage.label }}</span>
          <span v-if="stage.hint && !compact" class="mm-financial-report-loading__step-hint">
            {{ stage.hint }}
          </span>
        </span>
      </li>
    </ol>
  </div>
</template>
