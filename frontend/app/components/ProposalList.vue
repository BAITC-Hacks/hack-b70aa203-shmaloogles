<script setup lang="ts">
const props = defineProps<{ taskId?: string }>();
const { proposals, tasks, decideProposal } = useDemo();
const { role } = useRole();
const status = ref('all');
const allProposals = computed(() => proposals.value.filter(item => !props.taskId || item.taskId === props.taskId));
const filtered = computed(() => allProposals.value.filter(item => status.value === 'all' || item.status === status.value));
const statuses = [
  { key: 'all', label: 'Все предложения' }, { key: 'pending', label: 'На рассмотрении' },
  { key: 'accepted', label: 'Принятые' }, { key: 'rejected', label: 'Отклонённые' },
];
const statusLabel = { pending: 'На рассмотрении', accepted: 'Принято', rejected: 'Отклонено' };
function safeUrl(url?: string) {
  if (!url) return undefined;
  try { return ['http:', 'https:'].includes(new URL(url).protocol) ? url : undefined; } catch { return undefined; }
}
</script>

<template>
  <div>
    <div class="catalog-tabs-row"><div class="filter-tabs" role="group" aria-label="Статус предложения"><button v-for="item in statuses" :key="item.key" :class="{ active: status === item.key }" :aria-pressed="status === item.key" @click="status = item.key">{{ item.label }}<span v-if="item.key === 'all'">{{ allProposals.length }}</span></button></div></div>
    <div v-if="filtered.length" class="proposal-list">
      <article v-for="proposal in filtered" :key="proposal.id" class="panel proposal-card">
        <div class="proposal-card-heading"><span class="team-avatar">{{ proposal.team.substring(0, 2).toUpperCase() }}</span><div><h2>{{ proposal.team }}</h2><NuxtLink :to="`/catalog/${proposal.taskId}`">{{ tasks.find(task => task.id === proposal.taskId)?.card.title || 'Задача' }}<AppIcon name="arrowUp" :size="13" /></NuxtLink></div><span class="badge" :class="proposal.status === 'accepted' ? 'badge-success' : proposal.status === 'rejected' ? 'badge-neutral' : 'badge-blue'"><AppIcon :name="proposal.status === 'accepted' ? 'check' : proposal.status === 'rejected' ? 'close' : 'clock'" :size="13" />{{ statusLabel[proposal.status] }}</span></div>
        <div class="proposal-content"><div><h3>Идея решения</h3><p>{{ proposal.idea }}</p></div><div><h3>План команды</h3><p class="preserve-lines">{{ proposal.plan }}</p></div></div>
        <div class="proposal-card-bottom"><div class="proposal-meta"><span><AppIcon name="clock" :size="16" />{{ proposal.timeframe }}</span><a v-if="safeUrl(proposal.prototypeUrl)" :href="safeUrl(proposal.prototypeUrl)" target="_blank" rel="noopener noreferrer" class="text-link"><AppIcon name="link" :size="16" />Прототип <AppIcon name="arrowUp" :size="13" /></a></div><div v-if="role === 'business' && proposal.status === 'pending'" class="proposal-actions"><button class="button button-outline button-small" @click="decideProposal(proposal.id, 'rejected')">Отклонить</button><button class="button button-primary button-small" @click="decideProposal(proposal.id, 'accepted')"><AppIcon name="check" :size="16" />Принять предложение</button></div><span v-else class="proposal-status-note" role="status">{{ proposal.status === 'accepted' ? 'Начало хорошего сотрудничества' : proposal.status === 'rejected' ? 'Бизнес выбрал другой подход' : 'Ожидаем решения бизнеса' }}</span></div>
      </article>
    </div>
    <EmptyState v-else icon="message" title="Новые возможности ещё впереди" :description="status === 'all' ? 'Здесь появятся предложения команд и решения по ним. Начните со знакомства с задачами.' : 'Предложений с таким статусом пока нет. Загляните в соседнюю вкладку.'"><NuxtLink v-if="status === 'all'" to="/catalog" class="button button-primary">Открыть каталог <AppIcon name="arrow" :size="17" /></NuxtLink><button v-else class="button button-outline" @click="status = 'all'">Все предложения</button></EmptyState>
  </div>
</template>
