<script setup lang="ts">
import type { Task } from '~/types/task';
defineProps<{ task?: Task }>();
</script>
<template>
  <div class="hero-project">
    <div class="project-caption"><span>ИЗ КАТАЛОГА</span><svg width="48" height="22" viewBox="0 0 48 22" fill="none" aria-hidden="true"><path d="M1 20H47M8 20V15C8-3 40-3 40 15V20M16 20V15C16 7 32 7 32 15V20" stroke="currentColor" stroke-width="2" /></svg></div>
    <template v-if="task"><NuxtLink :to="`/catalog/${task.id}`" class="hero-project-title">{{ task.card.title || 'Задача без названия' }} <span aria-hidden="true">↗</span></NuxtLink><p>{{ task.card.expectedResult || 'Ожидаемый результат пока не указан' }}</p><ReadinessScore :score="task.score" :breakdown="task.breakdown" compact /></template>
    <template v-else><h2>Сначала — понятная задача</h2><p>Уточните проблему, данные и результат. Готовность покажет, какие сведения ещё нужны команде.</p><NuxtLink to="/catalog" class="text-link">Открыть каталог →</NuxtLink></template>
  </div>
</template>
