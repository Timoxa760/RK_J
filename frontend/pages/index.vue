<script setup lang="ts">
import {
  ArrowRight,
  CalendarDays,
  Mic,
  Moon,
  ShoppingBag,
  Sparkles,
  TrendingDown
} from 'lucide-vue-next'

definePageMeta({ layout: 'welcome' })

useHead({
  title: 'Поток — голосовой помощник по тратам',
  meta: [
    {
      name: 'description',
      content:
        'Деньги идут потоком — Поток помогает их разобрать. Голосовой анализ расходов после покупки, вечером и раз в неделю. Без таблиц и подключения банка.'
    }
  ]
})

const moments = [
  {
    icon: ShoppingBag,
    title: 'После траты',
    description:
      'Новая покупка — новая волна. Одной фразой фиксируете расход, Поток показывает, вписывается ли она в ваш обычный ритм.'
  },
  {
    icon: Moon,
    title: 'Вечером',
    description:
      'Краткий итог дня: куда ушёл поток трат и где он вышел за привычные берега.'
  },
  {
    icon: CalendarDays,
    title: 'Раз в неделю',
    description:
      'Смотрите на недельный поток целиком — одна рекомендация, куда его направить в следующие семь дней.'
  }
]

const benefits = [
  {
    icon: Mic,
    title: 'Голос вместо таблиц',
    description:
      'Описываете траты своими словами — Поток сам направит их в нужное русло, без категорий и чеков.'
  },
  {
    icon: Sparkles,
    title: 'Разбор потока',
    description:
      'Не сухой отчёт, а понятный разбор: что в норме, что разлилось и почему это важно сейчас.'
  },
  {
    icon: TrendingDown,
    title: 'Одна правка курса',
    description:
      'Не десять советов — одно действие, которое сужает поток там, где он слишком широкий.'
  }
]

const exampleLines = [
  {
    role: 'user',
    text: 'Вышел из «Пятёрочки» — потратил 1 332 ₽'
  },
  {
    role: 'assistant',
    text: 'На 11% выше вашего обычного чека в продуктовых. В рамках недели, но ближе к верхней границе — имеет смысл отметить, что именно добавило сумму.'
  }
]
</script>

<template>
  <section class="relative mx-auto max-w-[1400px] px-4 pb-20 pt-6 sm:px-8 sm:pb-24 sm:pt-10 lg:px-12 lg:pt-14">
    <div class="relative z-10 grid gap-10 lg:grid-cols-[1fr_minmax(280px,380px)] lg:items-center lg:gap-12">
      <div class="max-w-xl">
        <p class="text-[10px] font-medium uppercase tracking-[0.18em] text-[color:var(--mm-primary)] sm:text-xs sm:tracking-[0.25em]">
          Поток денег под контролем
        </p>

        <h1 class="mm-hero-title mt-4 text-[color:var(--mm-text)] sm:mt-5">
          <span class="block">Расходы текут —</span>
          <span class="block mm-hero-title__flow-line">
            <SharedHeroFlowWord />
          </span>
          <span class="block text-[color:var(--mm-primary)]">их разберёт</span>
        </h1>

        <p class="mt-6 text-[15px] leading-relaxed text-[color:var(--mm-text-muted)] sm:mt-7 sm:text-lg">
          Деньги не стоят на месте — они текут. Вы говорите, что потратили,
          а Поток показывает, куда уходит бюджет и что стоит скорректировать.
          Без таблиц и подключения банка.
        </p>

        <div class="mt-8 flex flex-col gap-3 sm:flex-row sm:flex-wrap">
          <NuxtLink to="/login" class="mm-btn-primary group w-full sm:w-auto">
            Войти в Поток
            <ArrowRight class="h-4 w-4 transition group-hover:translate-x-0.5" stroke-width="2" />
          </NuxtLink>
          <NuxtLink to="/login" class="mm-btn-secondary w-full sm:w-auto">
            <Mic class="h-4 w-4" stroke-width="1.75" />
            Как это работает
          </NuxtLink>
        </div>

        <p class="mt-5 text-sm text-[color:var(--mm-text-soft)]">
          Первый разбор — за пару минут.
        </p>
      </div>

      <div class="mm-voice-demo mm-card relative overflow-hidden p-5 sm:p-6">
        <div class="flex items-center gap-3 border-b border-[color:var(--mm-border-subtle)] pb-4">
          <div
            class="flex h-11 w-11 shrink-0 items-center justify-center rounded-full bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-primary)]"
          >
            <Mic class="h-5 w-5" stroke-width="1.75" />
          </div>
          <div>
            <p class="text-sm font-medium text-[color:var(--mm-text)]">Поток слушает</p>
            <p class="text-xs text-[color:var(--mm-text-soft)]">После похода в магазин</p>
          </div>
          <div class="mm-voice-demo__bars ml-auto flex items-end gap-0.5" aria-hidden="true">
            <span /><span /><span /><span /><span />
          </div>
        </div>

        <div class="mt-4 space-y-3">
          <div
            v-for="(line, index) in exampleLines"
            :key="index"
            class="rounded-2xl px-3.5 py-2.5 text-sm leading-relaxed"
            :class="
              line.role === 'user'
                ? 'ml-4 bg-[color:var(--mm-bg-muted)] text-[color:var(--mm-text-muted)]'
                : 'mr-2 border border-[color:var(--mm-primary-muted)] bg-[color:var(--mm-primary-soft)] text-[color:var(--mm-text)]'
            "
          >
            {{ line.text }}
          </div>
        </div>
      </div>
    </div>

    <div id="when" class="relative z-10 mt-16 sm:mt-20">
      <h2 class="text-sm font-medium text-[color:var(--mm-text)]">Когда включать Поток</h2>
      <p class="mt-2 max-w-xl text-sm text-[color:var(--mm-text-muted)]">
        Встраивается в ваш ритм — там, где поток трат уже пошёл, а разобраться некогда.
      </p>

      <div class="mt-6 grid gap-4 sm:grid-cols-3">
        <article
          v-for="item in moments"
          :key="item.title"
          class="mm-card p-5 sm:p-6"
        >
          <component
            :is="item.icon"
            class="h-5 w-5 text-[color:var(--mm-primary)]"
            stroke-width="1.75"
          />
          <h3 class="mt-4 font-medium text-[color:var(--mm-text)]">{{ item.title }}</h3>
          <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
            {{ item.description }}
          </p>
        </article>
      </div>
    </div>

    <div id="how" class="relative z-10 mt-16 sm:mt-20">
      <h2 class="text-sm font-medium text-[color:var(--mm-text)]">Что даёт Поток</h2>
      <div class="mt-6 grid gap-4 sm:grid-cols-3">
        <article
          v-for="item in benefits"
          :key="item.title"
          class="mm-card flex flex-col p-5 sm:p-6"
        >
          <component
            :is="item.icon"
            class="h-5 w-5 text-[color:var(--mm-primary)]"
            stroke-width="1.75"
          />
          <h3 class="mt-4 font-medium text-[color:var(--mm-text)]">{{ item.title }}</h3>
          <p class="mt-2 flex-1 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
            {{ item.description }}
          </p>
        </article>
      </div>
    </div>

    <div class="relative z-10 mt-16 sm:mt-20">
      <div class="mm-card flex flex-col items-start gap-5 p-6 sm:flex-row sm:items-center sm:justify-between sm:p-8">
        <div class="max-w-lg">
          <h2 class="text-lg font-semibold text-[color:var(--mm-text)] sm:text-xl">
            Следующий поток трат — уже под контролем
          </h2>
          <p class="mt-2 text-sm leading-relaxed text-[color:var(--mm-text-muted)]">
            Подключите Поток и после следующей покупки получите разбор
            с одной рекомендацией — куда направить бюджет дальше.
          </p>
        </div>
        <NuxtLink to="/login" class="mm-btn-primary w-full shrink-0 sm:w-auto">
          Начать бесплатно
        </NuxtLink>
      </div>
    </div>

    <p class="relative z-10 mt-10 max-w-2xl text-sm leading-relaxed text-[color:var(--mm-text-soft)]">
      Не учёт ради учёта. Поток делает ваш денежный поток прозрачным —
      в момент, когда траты только что ушли, день закончился или подошла неделя.
    </p>
  </section>
</template>
