<script setup lang="ts">
const { role } = useRole();
useHead({ title: 'Предложения — Мост' });
</script>

<template>
  <div class="page"><div class="page-heading"><div><h1>{{ role === 'business' ? 'Предложения команд' : 'Мои предложения' }}<span class="heading-dot">.</span></h1><p>{{ role === 'business' ? 'Рассмотрите идеи и примите решение.' : 'Предложения выбранной команды.' }}</p></div><NuxtLink v-if="role === 'team'" to="/catalog" class="button button-primary">Найти новую задачу</NuxtLink></div><ProposalList /></div>
</template>
