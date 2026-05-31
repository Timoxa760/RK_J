import { createRequire } from 'node:module'
import { brandPalette } from './constants/brandPalette'
import { vueSfcRegisterTs } from './plugins/vue-sfc-register-ts'
import { patchVueRouterDevtools } from './plugins/vite/patch-vue-router-devtools'

// Vue SFC compiler needs TypeScript to resolve defineProps<ImportedType>() from reka-ui
createRequire(import.meta.url)('vue/compiler-sfc')

export default defineNuxtConfig({
  srcDir: '.',
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  experimental: {
    appManifest: false
  },
  modules: ['@nuxtjs/tailwindcss', '@vite-pwa/nuxt', '@pinia/nuxt', 'shadcn-nuxt'],
  build: {
    transpile: ['reka-ui', 'vue-sonner']
  },
  typescript: {
    strict: true,
    typeCheck: false
  },
  vite: {
    server: {
      proxy: {
        '/api': {
          target: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8000',
          changeOrigin: true
        }
      }
    },
    plugins: [patchVueRouterDevtools(), vueSfcRegisterTs()],
    optimizeDeps: {
      include: ['typescript'],
    },
    vue: {
      features: {
        prodDevtools: false
      }
    }
  },
  tailwindcss: {
    configPath: './tailwind.config.ts',
    cssPath: '~/assets/css/main.css'
  },
  shadcn: {
    prefix: '',
    componentDir: './components/ui'
  },
  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:8000',
      demoMode: process.env.NUXT_PUBLIC_DEMO_MODE !== 'false'
    }
  },
  app: {
    pageTransition: false,
    layoutTransition: false,
    head: {
      title: 'Поток — голосовой помощник по тратам',
      link: [],
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
        { name: 'theme-color', content: brandPalette.primary },
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
      theme_color: brandPalette.primary,
      background_color: brandPalette.bg,
      display: 'standalone',
      start_url: '/'
    }
  },
  routeRules: {
    '/': { ssr: true },
    '/login': { ssr: true },
    '/dashboard': { ssr: false },
    '/profile': { ssr: false },
    '/receipts': { ssr: false },
    '/credits': { ssr: false },
    '/onboarding': { ssr: false },
    '/advisor': { ssr: false }
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
