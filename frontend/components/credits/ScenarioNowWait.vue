<script setup lang="ts">
import type { MortgageBreakdownResponse } from '~/types/api'

defineProps<{
  breakdown: MortgageBreakdownResponse
}>()
</script>

<template>
  <div class="grid gap-4 lg:grid-cols-2">
    <Card class="border-primary/20">
      <CardHeader class="pb-2">
        <CardDescription class="text-primary">Если взять сейчас</CardDescription>
        <CardTitle class="text-base font-medium">Сценарий «сейчас»</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-sm leading-relaxed text-muted-foreground">{{ breakdown.scenario_now }}</p>
      </CardContent>
    </Card>

    <Card>
      <CardHeader class="pb-2">
        <CardDescription>Если подождать {{ breakdown.wait_months }} мес.</CardDescription>
        <CardTitle class="text-base font-medium">Сценарий «подождать»</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-sm leading-relaxed text-muted-foreground">{{ breakdown.scenario_wait }}</p>
      </CardContent>
    </Card>
  </div>
</template>
