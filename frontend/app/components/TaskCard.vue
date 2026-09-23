<script setup lang="ts">
import type { Task } from '~/types/task';
defineProps<{ task: Task; manage?: boolean; registry?: boolean }>();
</script>
<template>
  <article class="task-card" :class="{ 'task-registry-row': registry }">
    <div class="task-card-top"><span class="task-reference">№ {{ task.id }}</span><span class="category-label">{{ topicLabel(task.card.topic) }}</span></div>
    <h3><NuxtLink :to="manage ? `/tasks/${task.id}/edit` : `/catalog/${task.id}`">{{ task.card.title || 'Без названия' }}</NuxtLink></h3>
    <p class="task-description">{{ task.card.need || task.description }}</p>
    <div class="task-result"><span>ОЖИДАЕМЫЙ РЕЗУЛЬТАТ</span><p>{{ task.card.expectedResult || 'Пока не указан' }}</p></div>
    <ReadinessScore :score="task.score" :breakdown="task.breakdown" compact />
    <div class="task-card-footer"><span>{{ task.status === 'published' ? 'Опубликована' : task.status === 'confirmed' ? 'Подтверждена' : 'Не опубликована' }}</span><NuxtLink :to="manage ? `/tasks/${task.id}/proposals` : `/catalog/${task.id}`">{{ manage ? 'Отклики' : 'Открыть задачу' }} <span aria-hidden="true">→</span></NuxtLink></div>
  </article>
</template>
