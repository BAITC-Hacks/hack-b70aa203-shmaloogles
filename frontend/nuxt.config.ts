// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  typescript: { strict: true },
  css: ['~/assets/css/main.css', '~/assets/css/bureau.css'],
  app: {
    head: {
      htmlAttrs: { lang: 'ru' },
      title: 'Мост — бизнес, идеи, команды',
      meta: [
        { name: 'description', content: 'Превратите задачу бизнеса в реальный проект. Мост соединяет компании и студенческие команды, готовые создавать новое.' },
        { name: 'theme-color', content: '#f5f8fc' },
      ],
      link: [{ rel: 'icon', type: 'image/svg+xml', href: '/favicon.svg' }],
    },
  },
  runtimeConfig: {
    apiBase: 'http://localhost:8080',
    public: {
      apiBase: 'http://localhost:8080',
    },
  },
})
