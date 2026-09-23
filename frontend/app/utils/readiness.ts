import type { ReadinessScore, TaskCard } from '~/types/task';

export const readinessFields = [
  { key: 'context', label: 'Контекст и потребность', weight: 20 },
  { key: 'data', label: 'Данные и материалы', weight: 20 },
  { key: 'expectedResult', label: 'Ожидаемый результат', weight: 15 },
  { key: 'successCriteria', label: 'Критерии успеха', weight: 15 },
  { key: 'constraints', label: 'Ограничения', weight: 10 },
  { key: 'users', label: 'Пользователи', weight: 10 },
  { key: 'contact', label: 'Связь с бизнесом', weight: 10 },
] as const;

export function isFilled(value: string | string[]): boolean {
  return Array.isArray(value) ? value.some(item => item.trim().length > 0) : value.trim().length > 0;
}

export function calculateReadiness(card: TaskCard): ReadinessScore {
  const total = readinessFields.reduce((sum, field) => sum + (isFilled(card[field.key]) ? field.weight : 0), 0);
  const level = total >= 90 ? 'Приоритетная' : total >= 70 ? 'Готовая' : total >= 40 ? 'Рабочая' : 'Черновик';
  return { total, level };
}

export const toLines = (value: string): string[] => value.split('\n').map(line => line.trim()).filter(Boolean);
