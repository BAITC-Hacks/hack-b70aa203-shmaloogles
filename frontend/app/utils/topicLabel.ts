/** Present known catalog codes without changing stored values or custom topics. */
export function topicLabel(topic: string): string {
  const labels: Record<string, string> = {
    analytics: 'Аналитика', automation: 'Автоматизация', ecology: 'Экология',
    education: 'Образование', marketing: 'Маркетинг', logistics: 'Логистика',
  };
  return labels[topic.toLowerCase()] || topic || 'Тема не указана';
}
