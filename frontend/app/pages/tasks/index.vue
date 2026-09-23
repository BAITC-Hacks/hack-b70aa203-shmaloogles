<script setup lang="ts">
const { tasks, proposals } = useDemo();
const { role, setRole } = useRole();
const status = ref('all');
const visible = computed(() => tasks.value.filter(task => status.value === 'all' || task.status === status.value));
const draftCount = computed(() => tasks.value.filter(task => task.status === 'draft').length);
const pendingCount = computed(() => proposals.value.filter(proposal => proposal.status === 'pending').length);
useHead({ title: 'Мои задачи — Мост' });
</script>

<template>
  <div class="page">
    <div class="page-heading"><div><span class="eyebrow muted">ВАШИ ИДЕИ В ДВИЖЕНИИ</span><h1>Мои задачи<span class="heading-dot">.</span></h1><p>От первого наброска до команды, с которой всё получится.</p></div><NuxtLink v-if="role === 'business'" to="/tasks/new" class="button button-primary"><AppIcon name="plus" :size="18" />Новая задача</NuxtLink></div>
    <template v-if="role === 'business'">
      <div class="stats-grid"><div class="stat-card"><span class="stat-icon tone-blue"><AppIcon name="layers" /></span><div><span>Всего задач</span><strong>{{ tasks.length }}<small>в вашем пространстве</small></strong></div></div><div class="stat-card"><span class="stat-icon tone-peach"><AppIcon name="edit" /></span><div><span>Черновики</span><strong>{{ draftCount }}<small>идеи на старте</small></strong></div></div><div class="stat-card"><span class="stat-icon tone-mint"><AppIcon name="message" /></span><div><span>Новые предложения</span><strong>{{ pendingCount }}<small>ждут вашего решения</small></strong></div></div></div>
      <div class="catalog-tabs-row"><div class="filter-tabs" role="group" aria-label="Статус задачи"><button :class="{ active: status === 'all' }" :aria-pressed="status === 'all'" @click="status = 'all'">Все задачи <span>{{ tasks.length }}</span></button><button :class="{ active: status === 'published' }" :aria-pressed="status === 'published'" @click="status = 'published'">Опубликованные</button><button :class="{ active: status === 'draft' }" :aria-pressed="status === 'draft'" @click="status = 'draft'">Черновики <span>{{ draftCount }}</span></button></div><span class="demo-caption">Демо-пространство бизнеса</span></div>
      <div v-if="visible.length" class="task-grid tasks-management"><TaskCard v-for="task in visible" :key="task.id" :task="task" manage /></div>
      <EmptyState v-else icon="file" :title="status === 'draft' ? 'Все идеи уже увидели свет' : 'Здесь начнётся что-то хорошее'" description="Создайте задачу — мы поможем превратить вашу идею в понятный проект."><NuxtLink to="/tasks/new" class="button button-primary"><AppIcon name="plus" :size="18" />Создать задачу</NuxtLink></EmptyState>
    </template>
    <EmptyState v-else icon="briefcase" title="Пространство для ваших бизнес-задач" description="Чтобы создавать задачи и рассматривать предложения, переключитесь в роль бизнеса."><button class="button button-primary" @click="setRole('business')">Перейти в роль бизнеса <AppIcon name="arrow" :size="18" /></button><NuxtLink to="/catalog" class="text-link">Я хочу найти проект</NuxtLink></EmptyState>
  </div>
</template>
