<script setup lang="ts">
import { chartHeightClass } from '~/utils/chartTheme'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    title: string
    description?: string
    size?: 'sm' | 'md' | 'lg' | 'full'
    loading?: boolean
    colSpan?: '1' | '2'
  }>(),
  {
    size: 'md',
    colSpan: '1'
  }
)

const sizeClass = computed(() => {
  const map = {
    sm: 'mm-chart-wrap--sm',
    md: 'mm-chart-wrap--md',
    lg: 'mm-chart-wrap--lg',
    full: 'mm-chart-wrap--full'
  }
  return map[props.size]
})

const colClass = computed(() => (props.colSpan === '2' ? 'lg:col-span-2' : ''))
</script>

<template>
  <Card :class="[colClass, 'flex h-full flex-col']" v-bind="$attrs">
    <CardHeader class="shrink-0">
      <CardTitle class="text-base">{{ title }}</CardTitle>
      <CardDescription v-if="description">{{ description }}</CardDescription>
    </CardHeader>
    <CardContent class="flex min-h-0 flex-1 flex-col overflow-hidden pt-0">
      <Skeleton v-if="loading" :class="['w-full shrink-0', chartHeightClass(size)]" />
      <div v-else :class="['mm-chart-wrap shrink-0', sizeClass]">
        <slot />
      </div>
    </CardContent>
  </Card>
</template>
