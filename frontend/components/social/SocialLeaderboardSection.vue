<script setup lang="ts">
import { Trophy } from 'lucide-vue-next'
import type { LeaderboardEntry } from '~/types/api'

defineProps<{
  leaderboard: LeaderboardEntry[]
  loading?: boolean
  challengeTitle?: string | null
}>()
</script>

<template>
  <Card data-demo="leaderboard">
    <CardHeader>
      <CardTitle class="text-base">Рейтинг друзей</CardTitle>
      <CardDescription>
        Относительный рейтинг без сумм
        <span v-if="challengeTitle"> — «{{ challengeTitle }}»</span>
      </CardDescription>
    </CardHeader>
    <CardContent>
      <Skeleton v-if="loading" class="h-40 w-full" />

      <div
        v-else-if="!leaderboard.length"
        class="rounded-lg border border-dashed px-4 py-10 text-center"
      >
        <Trophy class="mx-auto size-8 text-muted-foreground/60" aria-hidden="true" />
        <p class="mt-3 text-sm font-medium">Пока пусто</p>
        <p class="mt-1 text-sm text-muted-foreground">
          Выберите задание или создайте своё — таблица заполнится, когда подключатся друзья.
        </p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full min-w-[280px] text-sm">
          <thead>
            <tr class="border-b text-left text-muted-foreground">
              <th class="pb-2 pr-4">#</th>
              <th class="pb-2 pr-4">Участник</th>
              <th class="pb-2 text-right">Балл</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="row in leaderboard"
              :key="row.position"
              class="border-b"
              :class="row.display_name === 'Вы' ? 'bg-primary/10' : ''"
            >
              <td class="py-3 pr-4">{{ row.rank ?? row.position }}</td>
              <td class="py-3 pr-4 font-medium">{{ row.display_name ?? row.username }}</td>
              <td class="py-3 text-right text-muted-foreground">
                {{ row.relative_score.toLocaleString('ru-RU', { maximumFractionDigits: 2 }) }}
              </td>
            </tr>
          </tbody>
        </table>
        <p class="mt-3 text-xs text-muted-foreground">
          Чем ниже балл, тем лучше результат — без раскрытия сумм трат.
        </p>
      </div>
    </CardContent>
  </Card>
</template>
