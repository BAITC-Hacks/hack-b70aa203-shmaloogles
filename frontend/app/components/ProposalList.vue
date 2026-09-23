<script setup lang="ts">
const { t } = useLocale();

import type { ProposalStatus } from '~/types/task';
const props = defineProps<{ taskId?: string }>();
const platform = usePlatformApi();
const { role, selectedTeamId } = useRole();
const status = ref('all');
const changing = ref('');
const actionError = ref('');
const { data: proposals, pending, error, refresh: reload } = await useAsyncData(() => `proposals-${props.taskId || 'all'}`, () => platform.listProposals(props.taskId));
const refresh = () => reload();
const visible = computed(() => (proposals.value || []).filter(item => (role.value !== 'team' || item.teamId === selectedTeamId.value) && (status.value === 'all' || item.status === status.value)));
const statuses = [{ key: 'all', label: 'Все' }, { key: 'pending', label: 'На рассмотрении' }, { key: 'accepted', label: 'Принятые' }, { key: 'rejected', label: 'Отклонённые' }];
async function decide(id: string, next: Exclude<ProposalStatus, 'pending'>) {
  changing.value = id; actionError.value = '';
  try { await platform.updateProposal(id, next); await refresh(); }
  catch (cause) { actionError.value = cause instanceof Error ? cause.message : 'Не удалось изменить статус'; }
  finally { changing.value = ''; }
}
</script>
<template><div><div class="catalog-tabs-row"><div class="filter-tabs" role="group" :aria-label="t('Фильтр предложений')"><button v-for="item in statuses" :key="item.key" :class="{ active: status === item.key }" :aria-pressed="status === item.key" @click="status = item.key">{{ t(item.label) }}</button></div><span class="demo-caption"><AppIcon name="message" :size="14" />{{ visible.length }} {{ t("предложений") }}</span></div><div v-if="pending" class="state-panel" role="status"><span class="state-icon"><AppIcon name="message" /></span><div><strong>{{ t("Загружаем предложения") }}</strong><p>{{ t("Собираем идеи и планы команд.") }}</p></div></div><div v-else-if="error" class="state-panel state-error" role="alert"><span class="state-icon"><AppIcon name="info" /></span><div><strong>{{ t("Не удалось загрузить предложения") }}</strong><button class="text-link" @click="refresh">{{ t("Попробовать снова") }}</button></div></div><p v-if="actionError" class="error-text inline-alert" role="alert">{{ t(actionError) }}</p><div v-if="!pending && !error && visible.length" class="proposal-list"><article v-for="proposal in visible" :key="proposal.id" class="panel proposal-card"><div class="proposal-card-heading"><span class="team-avatar">{{ proposal.team.slice(0, 2).toUpperCase() }}</span><div><h2>{{ proposal.team }}</h2><NuxtLink :to="`/catalog/${proposal.taskId}`">{{ t("Открыть задачу") }} <AppIcon name="arrow" :size="14" /></NuxtLink></div><span class="badge" :class="proposalBadgeClass(proposal.status)"><AppIcon :name="proposal.status === 'accepted' ? 'check' : proposal.status === 'rejected' ? 'close' : 'clock'" :size="13" />{{ proposalStatusLabel(proposal.status) }}</span></div><div class="proposal-content"><div><h3><AppIcon name="sparkles" :size="15" />{{ t("Идея решения") }}</h3><p class="preserve-lines">{{ proposal.idea }}</p></div><div><h3><AppIcon name="layers" :size="15" />{{ t("План команды") }}</h3><p class="preserve-lines">{{ proposal.plan }}</p></div></div><div class="proposal-card-bottom"><div class="proposal-meta"><span><AppIcon name="clock" :size="15" />{{ proposal.timeframe }}</span><a v-if="proposal.prototypeUrl" class="text-link" :href="proposal.prototypeUrl" target="_blank" rel="noopener noreferrer"><AppIcon name="link" :size="15" />{{ t("Прототип") }}</a></div><div v-if="role === 'business' && proposal.status === 'pending'" class="proposal-actions"><button class="button button-outline button-small" :disabled="changing === proposal.id" @click="decide(proposal.id, 'rejected')"><AppIcon name="close" :size="14" />{{ t("Отклонить") }}</button><button class="button button-primary button-small" :disabled="changing === proposal.id" @click="decide(proposal.id, 'accepted')"><AppIcon name="check" :size="14" />{{ changing === proposal.id ? t("Сохраняем…") : t("Принять") }}</button></div><span v-else class="proposal-status-note">{{ t("Статус:") }} {{ proposalStatusLabel(proposal.status) }}</span></div></article></div><EmptyState v-else-if="!pending && !error" :title="t('Предложений пока нет')" :description="t('Здесь появятся предложения выбранной команды или отклики на задачу.')" /></div></template>
