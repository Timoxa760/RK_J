export default defineNuxtConfig({
  srcDir: '.',
  compatibilityDate: '2025-07-15',
  modules: ['@nuxtjs/tailwindcss', '@vite-pwa/nuxt', '@pinia/nuxt'],
  css: ['~/assets/css/main.css'],
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8000',
      demoMode: process.env.NUXT_PUBLIC_DEMO_MODE !== 'false'
    }
  },
  app: {
    head: {
      title: 'Поток — голосовой помощник по тратам',
      meta: [
        {
          name: 'viewport',
          content: 'width=device-width, initial-scale=1, viewport-fit=cover'
        },
        {
          name: 'description',
          content:
            'Голосовой помощник анализирует траты после магазина, в конце дня и недели — и предлагает конкретную оптимизацию.'
        },
        { name: 'theme-color', content: '#fffcf9' },
        { property: 'og:title', content: 'Поток — голосовой помощник по тратам' },
        {
          property: 'og:description',
          content:
            'Скажите, что купили — Поток разберёт траты и подскажет, что можно улучшить уже на этой неделе.'
        },
        { property: 'og:type', content: 'website' }
      ]
    }
  },
  pwa: {
    registerType: 'autoUpdate',
    manifest: {
      name: 'Поток',
      short_name: 'Поток',
      description: 'Понятный контроль личных финансов',
      theme_color: '#fffcf9',
      background_color: '#fffcf9',
      display: 'standalone',
      start_url: '/'
    }
  },
  nitro: {
    devProxy: {
      '/api': {
        target: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8000',
        changeOrigin: true
      }
    }
  }
})
