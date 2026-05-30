<script setup lang="ts">
const year = new Date().getFullYear()
</script>

<template>
  <LandingReveal>
    <footer class="mm-landing-footer mm-safe-bottom">
      <p class="mm-landing-footer__text">
        © {{ year }} Поток. Все права защищены.
      </p>
    </footer>
  </LandingReveal>
</template>
