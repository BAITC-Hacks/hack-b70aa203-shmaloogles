<template>
  <article class="task-card">
    <div class="task-card-top"><span class="company-icon" :class="`tone-${task.color || 'blue'}`"><AppIcon :name="task.icon || 'briefcase'" :size="24" /></span><span class="category-label">{{ task.category || 'Задача бизнеса' }}</span><span v-if="task.status === 'draft'" class="badge badge-neutral">Черновик</span><AppIcon v-else name="arrowUp" :size="18" class="card-arrow" /></div>
    <span class="company-name">{{ task.organization || 'Компания' }}</span>
    <h3><NuxtLink :to="manage ? `/tasks/${task.id}/edit` : `/catalog/${task.id}`">{{ task.card.title }}</NuxtLink></h3>
    <p class="task-description">{{ task.description }}</p>
    <div class="tag-list"><span v-for="tag in task.tags" :key="tag" class="tag">{{ tag }}</span></div>
    <ReadinessScore :score="task.score" compact />
    <div class="task-card-footer"><span><span class="status-dot" />{{ task.status === 'published' ? 'Открыта для предложений' : 'Ждёт публикации' }}</span><NuxtLink :to="manage ? `/tasks/${task.id}/proposals` : `/catalog/${task.id}`" :aria-label="`${manage ? 'Предложения' : 'Подробнее'}: ${task.card.title}`">{{ manage ? 'Отклики' : 'Подробнее' }}<AppIcon name="arrow" :size="15" /></NuxtLink></div>
  </article>
</template>

<script setup lang="ts">
import type { Task } from '~/types/task';

defineProps<{ task: Task; manage?: boolean }>();
</script>
