<script setup lang="ts">
const { role, setRole } = useRole();
useHead({ title: 'Новая задача — Мост' });
</script>

<template>
  <div class="page">
    <NuxtLink to="/tasks" class="back-link"><AppIcon name="back" :size="17" />К моим задачам</NuxtLink>
    <div class="page-heading"><div><span class="eyebrow muted">КАЖДЫЙ ПРОЕКТ НАЧИНАЕТСЯ С МЫСЛИ</span><h1>Давайте создадим <span>что-то важное.</span></h1><p>Вам не нужно готовое техническое задание. Начните своими словами.</p></div></div>
    <ClarificationChat v-if="role === 'business'" />
    <EmptyState v-else icon="briefcase" title="Расскажите о задаче вашего бизнеса" description="Переключитесь в роль бизнеса, чтобы подготовить и опубликовать задачу."><button class="button button-primary" @click="setRole('business')">Перейти в роль бизнеса</button></EmptyState>
  </div>
</template>
