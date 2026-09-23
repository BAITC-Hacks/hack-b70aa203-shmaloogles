<script setup lang="ts">
import type { Task } from '~/types/task';
defineProps<{ task?: Task }>();
</script>
<template>
  <div class="hero-project">
    <div class="project-caption"><span>ИЗ КАТАЛОГА</span><BridgeRule /></div>
    <template v-if="task"><NuxtLink :to="`/catalog/${task.id}`" class="hero-project-title">{{ task.card.title || 'Задача без названия' }} <span aria-hidden="true">→</span></NuxtLink><p>{{ task.card.expectedResult || 'Ожидаемый результат пока не указан' }}</p><ReadinessScore :score="task.score" :breakdown="task.breakdown" compact /></template>
    <template v-else><h2>Сначала — понятная задача</h2><p>Уточните проблему, данные и результат. Готовность покажет, какие сведения ещё нужны команде.</p><NuxtLink to="/catalog" class="text-link">Открыть каталог →</NuxtLink></template>
  </div>
</template>
