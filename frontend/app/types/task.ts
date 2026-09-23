export type TaskStatus = 'draft' | 'confirmed' | 'published';
export type ProposalStatus = 'pending' | 'accepted' | 'rejected';

export interface TaskCard {
  title: string;
  topic: string;
  context: string;
  need: string;
  users: string;
  data: string;
  constraints: string;
  expectedResult: string;
  successCriteria: string;
  contact: string;
  interactionFormat: string;
}

export interface ReadinessBreakdown { name: string; points: number; maxPoints: number }
export interface ReadinessScore { total: number; level: string }

export interface Task {
  id: string;
  description: string;
  card: TaskCard;
  score: ReadinessScore;
  status: TaskStatus;
  breakdown: ReadinessBreakdown[];
  missingInformation: string[];
  suggestions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Team { id: string; name: string; interests: string[]; skills: string[]; technologies: string[] }
export interface Proposal {
  id: string;
  taskId: string;
  teamId: string;
  team: string;
  idea: string;
  plan: string;
  timeframe: string;
  prototypeUrl?: string;
  status: ProposalStatus;
}
