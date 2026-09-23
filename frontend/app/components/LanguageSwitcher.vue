<script setup lang="ts">
const { locale, t } = useLocale();
</script>

<template>
  <select v-model="locale" class="language-switcher" :aria-label="t('Язык интерфейса')">
    <option value="ru" lang="ru">Русский</option>
    <option value="kk" lang="kk">Қазақша</option>
    <option value="en" lang="en">English</option>
  </select>
</template>

<style scoped>
.language-switcher { max-width: 110px; padding: 6px; border: 1px solid #dce3ec; border-radius: 8px; background: white; color: #25334a; font: inherit; font-size: 12px; }
@media (max-width: 480px) { .language-switcher { max-width: 88px; } }
</style>
