<script setup lang="ts">
const { role } = useRole();
useHead({ title: 'Предложения — Мост' });
</script>

<template>
  <div class="page"><div class="page-heading"><div><span class="eyebrow muted">ТАК НАЧИНАЕТСЯ СОТРУДНИЧЕСТВО</span><h1>{{ role === 'business' ? 'Свежий взгляд на ваши задачи' : 'Мои предложения' }}<span class="heading-dot">.</span></h1><p>{{ role === 'business' ? 'Знакомьтесь с идеями команд и выбирайте подход, который вам близок.' : 'Ваши идеи уже в пути. Здесь можно следить за решениями бизнеса.' }}</p></div><NuxtLink v-if="role === 'team'" to="/catalog" class="button button-primary">Найти новую задачу <AppIcon name="arrow" :size="17" /></NuxtLink></div><div class="inline-notice"><AppIcon name="info" :size="19" /><p>Демо-пространство: здесь показаны примеры и отправленные вами предложения. Можно переключать роли и проверять весь сценарий.</p></div><ProposalList /></div>
</template>
