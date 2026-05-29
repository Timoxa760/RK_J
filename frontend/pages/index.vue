<script setup lang="ts">
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
</script>

<template>
  <div>
    <LandingHero />
    <LandingHowItWorks />
    <LandingWhen />
    <LandingBenefits />
    <LandingProof />
    <LandingCtaBand />
    <LandingFooter />
  </div>
</template>
