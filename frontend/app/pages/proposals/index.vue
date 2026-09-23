<script setup lang="ts">
const { role } = useRole();
useHead({ title: 'Предложения — Мост' });
</script>

<template>
  <div class="page"><div class="page-heading"><div><span class="eyebrow"><AppIcon name="message" :size="14" />РАБОТА С ОТКЛИКАМИ</span><h1>{{ role === 'business' ? 'Предложения команд' : 'Мои предложения' }}<span class="heading-dot">.</span></h1><p>{{ role === 'business' ? 'Сравните идеи, планы и сроки, затем примите решение.' : 'Следите за статусом предложений выбранной команды.' }}</p></div><NuxtLink v-if="role === 'team'" to="/catalog" class="button button-primary"><AppIcon name="search" :size="16" />Найти новую задачу</NuxtLink></div><ProposalList /></div>
</template>
