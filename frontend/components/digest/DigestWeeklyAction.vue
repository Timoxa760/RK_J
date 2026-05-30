<script setup lang="ts">
defineProps<{
  action: string
  periodLabel?: string
}>()
</script>

<template>
  <Card class="border-primary/25 bg-primary/5" data-demo="digest-action">
    <CardHeader class="pb-2">
      <CardDescription class="text-primary">Главное действие недели</CardDescription>
      <CardTitle v-if="periodLabel" class="text-xs font-normal text-muted-foreground">
        {{ periodLabel }}
      </CardTitle>
    </CardHeader>
    <CardContent>
      <p class="text-base leading-relaxed">{{ action }}</p>
    </CardContent>
  </Card>
</template>
