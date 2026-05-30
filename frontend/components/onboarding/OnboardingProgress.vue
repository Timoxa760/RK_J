<script setup lang="ts">
const props = defineProps<{
  current: number
  total: number
  label?: string
  stepLabels?: string[]
}>()

const percent = computed(() =>
  props.total > 0 ? Math.min(100, Math.round((props.current / props.total) * 100)) : 0
)

const showStepper = computed(
  () => props.stepLabels && props.stepLabels.length > 0 && props.stepLabels.length <= 8
)
</script>

<template>
  <div class="space-y-2" aria-label="Прогресс опроса">
    <div class="flex items-center justify-between text-xs text-[color:var(--mm-text-muted)]">
      <span>{{ label ?? `Шаг ${current} из ${total}` }}</span>
      <span class="tabular-nums font-medium">{{ percent }}%</span>
    </div>

    <div v-if="showStepper" class="mm-onb-stepper" role="presentation">
      <div
        v-for="(_, i) in stepLabels"
        :key="i"
        class="mm-onb-stepper__segment"
        :class="{
          'mm-onb-stepper__segment--done': i + 1 < current,
          'mm-onb-stepper__segment--current': i + 1 === current
        }"
      />
    </div>
    <div
      v-else
      class="h-1.5 overflow-hidden rounded-full"
      style="background-color: rgb(245 220 200 / 0.55)"
    >
      <div
        class="h-full rounded-full transition-all duration-500 ease-out"
        :style="{
          width: `${percent}%`,
          background: 'linear-gradient(90deg, #e8955f, #f0b07a)'
        }"
      />
    </div>

    <div v-if="showStepper && stepLabels" class="mm-onb-stepper__labels mm-onb-stepper__labels--dense">
      <span
        v-for="(name, i) in stepLabels"
        :key="name"
        class="truncate px-0.5"
        :class="{
          'mm-onb-stepper__label--active': i + 1 === current,
          'mm-onb-stepper__label--done': i + 1 < current
        }"
      >
        {{ name }}
      </span>
    </div>
  </div>
</template>
