<script setup lang="ts">
import { Copy, Users } from 'lucide-vue-next'
import type { ChallengeItem, ChallengeType } from '~/types/api'

const props = defineProps<{
  challenges: ChallengeItem[]
  selectedChallengeId: string | null
  inviteUrl: string
  inviteCopied: boolean
  loading?: boolean
  newTitle: string
  newType: ChallengeType
  newDuration: number
}>()

const emit = defineEmits<{
  'update:newTitle': [value: string]
  'update:newType': [value: ChallengeType]
  'update:newDuration': [value: number]
  create: []
  select: [id: string]
  copyInvite: []
}>()

const selectedChallenge = computed(() =>
  props.challenges.find((c) => c.id === props.selectedChallengeId) ?? null
)
</script>

<template>
  <Card>
    <CardHeader>
      <CardTitle class="text-base">Задания</CardTitle>
      <CardDescription>Кто меньше потратил, кто больше отложил, streak без срывов</CardDescription>
    </CardHeader>
    <CardContent class="space-y-6">
      <form class="flex flex-col gap-3" @submit.prevent="emit('create')">
        <Input
          :model-value="newTitle"
          type="text"
          placeholder="Название задания"
          @update:model-value="emit('update:newTitle', String($event))"
        />
        <div class="flex flex-col gap-2 sm:flex-row">
          <Select
            :model-value="newType"
            @update:model-value="emit('update:newType', $event as ChallengeType)"
          >
            <SelectTrigger class="w-full sm:flex-1">
              <SelectValue placeholder="Тип" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="least_spend">Меньше трат</SelectItem>
              <SelectItem value="most_saved">Больше накоплений</SelectItem>
              <SelectItem value="streak">Серия без срывов</SelectItem>
            </SelectContent>
          </Select>
          <Input
            :model-value="newDuration"
            type="number"
            min="1"
            max="30"
            class="sm:w-28"
            placeholder="Дней"
            @update:model-value="emit('update:newDuration', Number($event))"
          />
        </div>
        <Button type="submit" class="w-full sm:w-auto">Создать</Button>
      </form>

      <div
        v-if="inviteUrl"
        class="rounded-lg border border-dashed border-primary/30 bg-primary/5 px-4 py-3"
      >
        <p class="text-xs font-medium text-primary">Пригласить друзей</p>
        <p class="mt-1 break-all font-mono text-xs text-muted-foreground">{{ inviteUrl }}</p>
        <div class="mt-3 flex flex-wrap gap-2">
          <Button type="button" size="sm" variant="secondary" @click="emit('copyInvite')">
            <Copy class="mr-1.5 size-3.5" aria-hidden="true" />
            {{ inviteCopied ? 'Скопировано' : 'Копировать ссылку' }}
          </Button>
        </div>
        <p class="mt-2 text-xs text-muted-foreground">
          Рейтинг анонимный — видны только относительные баллы, не суммы трат.
        </p>
      </div>

      <div v-if="loading && !challenges.length" class="space-y-2">
        <Skeleton class="h-14 w-full" />
        <Skeleton class="h-14 w-full" />
      </div>

      <div
        v-else-if="!challenges.length"
        class="rounded-lg border border-dashed px-4 py-8 text-center"
      >
        <Users class="mx-auto size-8 text-muted-foreground/60" aria-hidden="true" />
        <p class="mt-3 text-sm font-medium">Заданий пока нет</p>
        <p class="mt-1 text-sm text-muted-foreground">
          Создайте первый — отправьте ссылку друзьям, когда появится приглашение.
        </p>
      </div>

      <ul v-else class="space-y-2">
        <li
          v-for="challenge in challenges"
          :key="challenge.id"
          class="flex cursor-pointer flex-col gap-1 rounded-lg border px-4 py-3 text-sm transition-colors sm:flex-row sm:items-center sm:justify-between"
          :class="
            selectedChallengeId === challenge.id
              ? 'border-primary bg-primary/10'
              : 'hover:bg-muted/40'
          "
          @click="emit('select', challenge.id)"
        >
          <span class="font-medium">{{ challenge.title }}</span>
          <span class="text-muted-foreground">{{ challenge.participants }} участников</span>
        </li>
      </ul>

      <p v-if="selectedChallenge && !inviteUrl" class="text-xs text-muted-foreground">
        У «{{ selectedChallenge.title }}» нет ссылки приглашения — создайте своё задание, чтобы
        пригласить друзей.
      </p>
    </CardContent>
  </Card>
</template>
