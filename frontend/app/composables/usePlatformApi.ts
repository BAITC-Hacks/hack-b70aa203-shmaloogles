import type { ClarificationDto, GenerationDto, ProposalDto, TaskDto, TeamDto } from '~/types/api';
import { mapProposal, mapTask, mapTeam, toCardDto } from '~/types/api';
import type { ProposalStatus, TaskCard } from '~/types/task';

export function usePlatformApi() {
  const api = useApi();
  return {
    listTasks: async (all = false) => (await api.get<TaskDto[]>(`/api/tasks${all ? '?scope=all' : ''}`)).map(mapTask),
    getTask: async (id: string) => mapTask(await api.get<TaskDto>(`/api/tasks/${id}`)),
    createTask: async (description: string) => mapTask(await api.post<TaskDto>('/api/tasks', { initial_description: description })),
    updateTask: async (id: string, card: TaskCard) => mapTask(await api.put<TaskDto>(`/api/tasks/${id}`, toCardDto(card))),
    confirmTask: async (id: string) => mapTask(await api.post<TaskDto>(`/api/tasks/${id}/confirm`, {})),
    publishTask: async (id: string) => mapTask(await api.post<TaskDto>(`/api/tasks/${id}/publish`, {})),
    clarify: (description: string) => api.post<ClarificationDto>('/api/tasks/clarify', { description }),
    generate: (description: string, answers: { field: string; answer: string }[]) => api.post<GenerationDto>('/api/tasks/generate', { description, answers }),
    listTeams: async () => (await api.get<TeamDto[]>('/api/teams')).map(mapTeam),
    listProposals: async (taskId?: string) => (await api.get<ProposalDto[]>(taskId ? `/api/tasks/${taskId}/proposals` : '/api/proposals')).map(mapProposal),
    createProposal: async (taskId: string, input: { teamId: string; idea: string; plan: string; timeframe: string; prototypeUrl?: string }) => mapProposal(await api.post<ProposalDto>(`/api/tasks/${taskId}/proposals`, { team_id: Number(input.teamId), solution_idea: input.idea, plan: input.plan, timeline: input.timeframe, prototype_url: input.prototypeUrl || null })),
    updateProposal: async (id: string, status: Exclude<ProposalStatus, 'pending'>) => mapProposal(await api.patch<ProposalDto>(`/api/proposals/${id}`, { status })),
  };
}
