<script setup lang="ts">
import type { Task } from '~/types/task';
defineProps<{ task: Task; manage?: boolean }>();
</script>
<template>
  <article class="task-card">
    <div class="task-card-top"><span class="category-label">{{ task.card.topic || 'Тема не указана' }}</span></div>
    <h3><NuxtLink :to="manage ? `/tasks/${task.id}/edit` : `/catalog/${task.id}`">{{ task.card.title || 'Без названия' }}</NuxtLink></h3>
    <p class="task-description">{{ task.card.need || task.description }}</p>
    <div class="task-result"><span>ОЖИДАЕМЫЙ РЕЗУЛЬТАТ</span><p>{{ task.card.expectedResult || 'Пока не указан' }}</p></div>
    <ReadinessScore :score="task.score" :breakdown="task.breakdown" compact />
    <div class="task-card-footer"><span>{{ task.status === 'published' ? 'Опубликована' : task.status === 'confirmed' ? 'Подтверждена' : 'Не опубликована' }}</span><NuxtLink :to="manage ? `/tasks/${task.id}/proposals` : `/catalog/${task.id}`">{{ manage ? 'Отклики' : 'Подробнее' }} <span aria-hidden="true">↗</span></NuxtLink></div>
  </article>
</template>
