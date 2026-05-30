<script setup lang="ts">
import { chartHeightClass } from '~/utils/chartTheme'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    title: string
    /** Одна строка пояснения под заголовком */
    subtitle?: string
    size?: 'sm' | 'md' | 'lg' | 'full' | 'auto'
    loading?: boolean
    colSpan?: '1' | '2'
  }>(),
  {
    size: 'md',
    colSpan: '1'
  }
)

const sizeClass = computed(() => {
  if (props.size === 'auto') return ''
  const map = {
    sm: 'mm-chart-wrap--sm',
    md: 'mm-chart-wrap--md',
    lg: 'mm-chart-wrap--lg',
    full: 'mm-chart-wrap--full'
  }
  return map[props.size]
})

const colClass = computed(() => (props.colSpan === '2' ? 'lg:col-span-2' : ''))

const skeletonClass = computed(() => {
  if (props.size === 'auto') return 'h-32'
  return chartHeightClass(props.size)
})
</script>

<template>
  <Card :class="[colClass, 'flex h-full flex-col']" v-bind="$attrs">
    <CardHeader class="gap-1.5 p-4 pb-2 sm:p-5 sm:pb-3">
      <CardTitle class="text-lg font-semibold leading-snug">{{ title }}</CardTitle>
      <p v-if="subtitle" class="text-base leading-relaxed text-muted-foreground">
        {{ subtitle }}
      </p>
    </CardHeader>
    <CardContent class="flex min-h-0 flex-1 flex-col p-4 pt-0 sm:p-5 sm:pt-0">
      <Skeleton v-if="loading" :class="['w-full shrink-0 rounded-xl', skeletonClass]" />
      <div v-else :class="size === 'auto' ? 'w-full' : ['mm-chart-wrap shrink-0', sizeClass]">
        <slot />
      </div>
    </CardContent>
  </Card>
</template>
