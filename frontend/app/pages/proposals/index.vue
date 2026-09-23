<script setup lang="ts">
const { t } = useLocale();

const { role } = useRole();
useHead(() => ({ title: t("Предложения — Мост") }));
</script>

<template>
  <div class="page"><div class="page-heading"><div><span class="eyebrow"><AppIcon name="message" :size="14" />{{ t("РАБОТА С ОТКЛИКАМИ") }}</span><h1>{{ role === 'business' ? t("Предложения команд") : t("Мои предложения") }}<span class="heading-dot">.</span></h1><p>{{ role === 'business' ? t("Сравните идеи, планы и сроки, затем примите решение.") : t("Следите за статусом предложений выбранной команды.") }}</p></div><NuxtLink v-if="role === 'team'" to="/catalog" class="button button-primary"><AppIcon name="search" :size="16" />{{ t("Найти новую задачу") }}</NuxtLink></div><ProposalList /></div>
</template>
