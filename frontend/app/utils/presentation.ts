const labels: Record<string, string> = {
  draft: 'Черновик', confirmed: 'Подтверждена', published: 'Опубликована',
  pending: 'На рассмотрении', accepted: 'Принято', rejected: 'Отклонено',
  workable: 'Можно начинать', ready: 'Готова к работе', priority: 'Высокая готовность',
  title: 'Название', topic: 'Тема', context: 'Контекст', need: 'Потребность', users: 'Пользователи',
  data: 'Данные и материалы', constraints: 'Ограничения', expected_result: 'Ожидаемый результат',
  expectedResult: 'Ожидаемый результат', success_criteria: 'Критерии успеха', successCriteria: 'Критерии успеха',
  contact: 'Контакт', interaction_format: 'Формат взаимодействия', interactionFormat: 'Формат взаимодействия',
  context_and_need: 'Контекст и потребность', business_communication: 'Связь с бизнесом',
  description: 'Описание', initial_description: 'Исходное описание', clarity: 'Ясность', completeness: 'Полнота',
  specificity: 'Конкретность', feasibility: 'Реализуемость', business_value: 'Ценность для бизнеса',
  education: 'Образование', healthcare: 'Здравоохранение', health: 'Здоровье', finance: 'Финансы',
  fintech: 'Финтех', retail: 'Ритейл', ecommerce: 'Электронная коммерция', 'e-commerce': 'Электронная коммерция',
  logistics: 'Логистика', marketing: 'Маркетинг', analytics: 'Аналитика', ai: 'Искусственный интеллект',
  'artificial intelligence': 'Искусственный интеллект', machine_learning: 'Машинное обучение',
  'machine learning': 'Машинное обучение', automation: 'Автоматизация', ecology: 'Экология',
  environment: 'Экология', agriculture: 'Сельское хозяйство', tourism: 'Туризм', hr: 'HR', support: 'Поддержка',
  social: 'Социальные проекты', cybersecurity: 'Кибербезопасность', it: 'IT', web: 'Веб-разработка', mobile: 'Мобильная разработка',
};

const fallback = (value: string) => value
  .replace(/[_-]+/g, ' ')
  .replace(/([a-zа-я])([A-ZА-Я])/g, '$1 $2')
  .trim()
  .replace(/^./, letter => letter.toLocaleUpperCase('ru-RU'));

export const presentationLabel = (value?: string | null, empty = 'Не указано') => {
  if (!value?.trim()) return empty;
  const key = value.trim();
  return labels[key] ?? labels[key.toLowerCase()] ?? fallback(key);
};

export const taskStatusLabel = (value: string) => presentationLabel(value);
export const proposalStatusLabel = (value: string) => presentationLabel(value);
export const readinessLabel = (value: string) => presentationLabel(value);
export const breakdownLabel = (value: string) => presentationLabel(value);
export const missingFieldLabel = (value: string) => presentationLabel(value);
export const topicLabel = (value?: string | null) => presentationLabel(value, 'Задача бизнеса');

export const proposalBadgeClass = (status: string) => ({
  accepted: 'badge-success', rejected: 'badge-danger', pending: 'badge-blue',
}[status] ?? 'badge-neutral');
