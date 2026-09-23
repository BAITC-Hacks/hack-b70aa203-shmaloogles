import type { Task, Proposal } from '~/types/task';

const storageKey = 'most-demo-v1';
const strings = (value: unknown): value is string[] => Array.isArray(value) && value.every(item => typeof item === 'string');

function isTask(value: unknown): value is Task {
  if (!value || typeof value !== 'object') return false;
  const task = value as Task;
  return typeof task.id === 'string' && typeof task.description === 'string'
    && ['draft', 'published'].includes(task.status) && !!task.card && !!task.score
    && typeof task.card.title === 'string' && typeof task.card.context === 'string'
    && typeof task.card.expectedResult === 'string' && typeof task.card.contact === 'string'
    && ['users', 'data', 'constraints', 'successCriteria'].every(key => strings(task.card[key as keyof typeof task.card]))
    && typeof task.score.total === 'number' && typeof task.score.level === 'string'
    && (task.tags === undefined || strings(task.tags));
}

function isProposal(value: unknown): value is Proposal {
  if (!value || typeof value !== 'object') return false;
  const proposal = value as Proposal;
  return ['id', 'taskId', 'team', 'idea', 'plan', 'timeframe'].every(key => typeof proposal[key as keyof Proposal] === 'string')
    && ['pending', 'accepted', 'rejected'].includes(proposal.status)
    && (proposal.prototypeUrl === undefined || typeof proposal.prototypeUrl === 'string');
}

export default defineNuxtPlugin((nuxtApp) => {
  const { tasks, proposals } = useDemo();
  const { role } = useRole();

  nuxtApp.hook('app:mounted', () => {
    try {
      const saved = JSON.parse(localStorage.getItem(storageKey) || 'null');
      if (saved && Array.isArray(saved.tasks) && saved.tasks.every(isTask) && Array.isArray(saved.proposals) && saved.proposals.every(isProposal)) {
        tasks.value = saved.tasks;
        proposals.value = saved.proposals;
        if (saved.role === 'business' || saved.role === 'team') role.value = saved.role;
      }
    } catch {
      // The preview also works when storage is disabled or contains outdated data.
    }

    watch([tasks, proposals, role], () => {
      try {
        localStorage.setItem(storageKey, JSON.stringify({ tasks: tasks.value, proposals: proposals.value, role: role.value }));
      } catch {
        // In restricted browsers, retain the current session in memory.
      }
    }, { deep: true });
  });
});
