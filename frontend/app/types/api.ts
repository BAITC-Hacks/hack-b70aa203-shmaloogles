import type { Proposal, ProposalStatus, Task, TaskCard, TaskStatus, Team } from './task';
import { topicLabel } from '~/utils/presentation';

export interface TaskCardDto {
  title: string | null; topic: string | null; context: string | null; need: string | null;
  users: string | null; data: string | null; constraints: string | null;
  expected_result: string | null; success_criteria: string | null; contact: string | null;
  interaction_format: string | null;
}
export interface TaskDto extends TaskCardDto {
  id: number; initial_description: string | null; clarification: unknown; status: TaskStatus;
  readiness_score: number; readiness_level: string;
  readiness_breakdown: { name: string; points: number; max_points: number }[] | Record<string, never> | null;
  missing_information: string[] | null; suggestions: string[] | null;
  confirmed_at: string | null; published_at: string | null; created_at: string; updated_at: string;
}
export interface TeamDto { id: number; name: string; interests: string[]; skills: string[]; technologies: string[]; created_at: string; updated_at: string }
export interface ProposalDto { id: number; task_id: number; team_id: number; team_name: string; solution_idea: string; plan: string; timeline: string; prototype_url: string | null; status: ProposalStatus; created_at: string; updated_at: string }
export interface QuestionDto { field: string; question: string }
export interface ClarificationDto { questions: QuestionDto[]; mode: string; fallback_reason?: string }
export interface GenerationDto { card: TaskCardDto; mode: string; fallback_reason?: string }

const text = (value: string | null) => value ?? '';
export const mapCard = (dto: TaskCardDto): TaskCard => ({
  title: text(dto.title), topic: dto.topic ? topicLabel(dto.topic) : '', context: text(dto.context), need: text(dto.need),
  users: text(dto.users), data: text(dto.data), constraints: text(dto.constraints),
  expectedResult: text(dto.expected_result), successCriteria: text(dto.success_criteria),
  contact: text(dto.contact), interactionFormat: text(dto.interaction_format),
});
export const toCardDto = (card: TaskCard): TaskCardDto => ({
  title: card.title || null, topic: card.topic || null, context: card.context || null, need: card.need || null,
  users: card.users || null, data: card.data || null, constraints: card.constraints || null,
  expected_result: card.expectedResult || null, success_criteria: card.successCriteria || null,
  contact: card.contact || null, interaction_format: card.interactionFormat || null,
});
export const mapTask = (dto: TaskDto): Task => ({
  id: String(dto.id), description: text(dto.initial_description), card: mapCard(dto), status: dto.status,
  score: { total: dto.readiness_score, level: dto.readiness_level },
  breakdown: (Array.isArray(dto.readiness_breakdown) ? dto.readiness_breakdown : []).map(item => ({ name: item.name, points: item.points, maxPoints: item.max_points })),
  missingInformation: dto.missing_information ?? [], suggestions: dto.suggestions ?? [],
  createdAt: dto.created_at, updatedAt: dto.updated_at,
});
export const mapTeam = (dto: TeamDto): Team => ({ ...dto, id: String(dto.id) });
export const mapProposal = (dto: ProposalDto): Proposal => ({ id: String(dto.id), taskId: String(dto.task_id), teamId: String(dto.team_id), team: dto.team_name, idea: dto.solution_idea, plan: dto.plan, timeframe: dto.timeline, prototypeUrl: dto.prototype_url ?? undefined, status: dto.status });
