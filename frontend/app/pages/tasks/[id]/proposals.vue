<script setup lang="ts">
const route = useRoute();
const { tasks } = useDemo();
const task = computed(() => tasks.value.find(item => item.id === route.params.id));
useHead({ title: 'Предложения по задаче — Мост' });
</script>

<template>
  <div class="page"><NuxtLink to="/tasks" class="back-link"><AppIcon name="back" :size="17" />К моим задачам</NuxtLink><template v-if="task"><div class="page-heading"><div><span class="eyebrow muted">У ВАШЕЙ ИДЕИ ЕСТЬ ОТКЛИК</span><h1>Предложения команд<span class="heading-dot">.</span></h1><p>{{ task.card.title }}</p></div><NuxtLink :to="task.status === 'published' ? `/catalog/${task.id}` : `/tasks/${task.id}/edit`" class="button button-outline">Карточка задачи <AppIcon name="arrowUp" :size="17" /></NuxtLink></div><ProposalList :task-id="task.id" /></template><EmptyState v-else icon="file" title="Задача не найдена" description="Вернитесь к вашим задачам, чтобы посмотреть предложения команд."><NuxtLink to="/tasks" class="button button-primary">Мои задачи</NuxtLink></EmptyState></div>
</template>
