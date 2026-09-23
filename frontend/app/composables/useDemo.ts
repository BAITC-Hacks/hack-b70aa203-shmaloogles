import { demoProposals, demoTasks } from '~/data/demo';
import type { Proposal, Task, TaskCard } from '~/types/task';
import { calculateReadiness } from '~/utils/readiness';

export function useDemo() {
  const tasks = useState<Task[]>('demo-tasks', () => structuredClone(demoTasks));
  const proposals = useState<Proposal[]>('demo-proposals', () => structuredClone(demoProposals));
  const publishedTasks = computed(() => tasks.value.filter(task => task.status === 'published'));

  function createTask(description: string, card: TaskCard) {
    const task: Task = {
      id: crypto.randomUUID(), description, card, score: calculateReadiness(card), status: 'draft',
      organization: 'Ваша компания', category: 'Новая задача', tags: ['Новая задача'], icon: 'briefcase', color: 'blue',
    };
    tasks.value.unshift(task);
    return task;
  }

  function updateTask(id: string, card: TaskCard, publish = false) {
    const task = tasks.value.find(item => item.id === id);
    if (!task) return;
    task.card = structuredClone(card);
    task.score = calculateReadiness(card);
    if (publish) task.status = 'published';
  }

  function submitProposal(proposal: Omit<Proposal, 'id' | 'status'>) {
    proposals.value.unshift({ ...proposal, id: crypto.randomUUID(), status: 'pending' });
  }

  function decideProposal(id: string, status: 'accepted' | 'rejected') {
    const proposal = proposals.value.find(item => item.id === id);
    if (proposal?.status === 'pending') proposal.status = status;
  }

  return { tasks, proposals, publishedTasks, createTask, updateTask, submitProposal, decideProposal };
}
