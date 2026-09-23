<script setup lang="ts">
const { t } = useLocale();

const { role, setRole } = useRole();
useHead(() => ({ title: t("Новая задача — Мост") }));
</script>

<template>
  <div class="page">
    <NuxtLink to="/tasks" class="back-link"><AppIcon name="back" :size="17" />{{ t("К моим задачам") }}</NuxtLink>
    <div class="page-heading"><div><span class="eyebrow muted">{{ t("КАЖДЫЙ ПРОЕКТ НАЧИНАЕТСЯ С МЫСЛИ") }}</span><h1>{{ t("Давайте создадим") }} <span>{{ t("что-то важное.") }}</span></h1><p>{{ t("Вам не нужно готовое техническое задание. Начните своими словами.") }}</p></div></div>
    <ClarificationChat v-if="role === 'business'" />
    <EmptyState v-else icon="briefcase" :title="t('Расскажите о задаче вашего бизнеса')" :description="t('Переключитесь в роль бизнеса, чтобы подготовить и опубликовать задачу.')"><button class="button button-primary" @click="setRole('business')">{{ t("Перейти в роль бизнеса") }}</button></EmptyState>
  </div>
</template>
