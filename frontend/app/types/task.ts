export interface ClarificationQuestion {
  id: string;
  question: string;
  field: string;
}

export interface TaskCard {
  title: string;
  context: string;
  users: string[];
  data: string[];
  constraints: string[];
  expectedResult: string;
  successCriteria: string[];
  contact: string;
}

export interface ReadinessScore {
  total: number;
  level: 'Черновик' | 'Рабочая' | 'Готовая' | 'Приоритетная';
}

export type TaskStatus = 'draft' | 'published';

export interface Task {
  id: string;
  description: string;
  card: TaskCard;
  score: ReadinessScore;
  status: TaskStatus;
}

export type ProposalStatus = 'pending' | 'accepted' | 'rejected';

export interface Proposal {
  id: string;
  taskId: string;
  team: string;
  idea: string;
  plan: string;
  timeframe: string;
  prototypeUrl?: string;
  status: ProposalStatus;
}
