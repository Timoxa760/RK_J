<script setup lang="ts">
const { challenges, leaderboard, loading, newTitle, loadAll, createChallenge } = useSocial()

onMounted(() => {
  loadAll()
})
</script>

<template>
  <div class="mm-page-shell">
    <SharedSkeletonLoader v-if="loading" height="400px" />

    <template v-else>
      <section class="mm-card p-4 sm:p-6">
        <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Создать челлендж</h2>
        <form class="mt-4 flex flex-col gap-2 sm:flex-row sm:flex-wrap" @submit.prevent="createChallenge(newTitle)">
          <input
            v-model="newTitle"
            type="text"
            placeholder="Название челленджа"
            class="w-full min-w-0 flex-1 rounded-lg border border-[color:var(--mm-border)] px-3 py-2.5 text-sm sm:min-w-[200px]"
          />
          <button type="submit" class="mm-btn-primary w-full !py-2.5 sm:w-auto sm:!px-4">
            Создать
          </button>
        </form>

        <ul class="mt-6 space-y-2">
          <li
            v-for="c in challenges"
            :key="c.id"
            class="flex flex-col gap-1 rounded-lg border border-[color:var(--mm-border-subtle)] px-4 py-3 text-sm sm:flex-row sm:justify-between sm:gap-2"
          >
            <span class="font-medium text-[color:var(--mm-text)]">{{ c.title }}</span>
            <span class="text-[color:var(--mm-text-soft)]">{{ c.participants }} участников</span>
          </li>
        </ul>
      </section>

      <section class="mm-card p-4 sm:p-6" data-demo="leaderboard">
        <h2 class="text-sm font-semibold text-[color:var(--mm-text)]">Лидерборд</h2>
        <p class="mt-1 text-xs text-[color:var(--mm-text-soft)]">Только относительный рейтинг, без сумм</p>
        <div class="mm-table-scroll mt-4">
          <table class="w-full min-w-[280px] text-sm">
            <thead>
              <tr class="border-b border-[color:var(--mm-border-subtle)] text-left text-[color:var(--mm-text-soft)]">
                <th class="pb-2 pr-4">#</th>
                <th class="pb-2 pr-4">Участник</th>
                <th class="pb-2 text-right">Балл</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="row in leaderboard"
                :key="row.id"
                class="border-b border-[color:var(--mm-border-subtle)]"
                :class="row.display_name === 'Вы' ? 'bg-[color:var(--mm-primary-soft)]' : ''"
              >
                <td class="py-3 pr-4">{{ row.rank }}</td>
                <td class="py-3 pr-4 font-medium text-[color:var(--mm-text)]">{{ row.display_name }}</td>
                <td class="py-3 text-right text-[color:var(--mm-text-muted)]">{{ row.relative_score }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </template>
  </div>
</template>
