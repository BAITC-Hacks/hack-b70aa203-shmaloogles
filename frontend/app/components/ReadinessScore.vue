<script setup lang="ts">
import type { ReadinessScore, ReadinessBreakdown } from '~/types/task';
const props = defineProps<{ score: ReadinessScore; breakdown?: ReadinessBreakdown[]; compact?: boolean; detailed?: boolean }>();
const categories = [
  ['context_and_need', 'Контекст и потребность', 20], ['data', 'Данные', 20],
  ['expected_result', 'Ожидаемый результат', 15], ['success_criteria', 'Критерии успеха', 15],
  ['constraints', 'Ограничения', 10], ['users', 'Пользователи', 10],
  ['business_communication', 'Связь с бизнесом', 10],
] as const;
const rows = computed(() => categories.map(([name, label, weight]) => ({ name, label, weight, item: props.breakdown?.find(item => item.name === name) })));
const validBreakdown = computed(() => rows.value.every(({ item, weight }) => item && item.maxPoints === weight && Number.isFinite(item.points) && item.points >= 0 && item.points <= weight) && rows.value.reduce((sum, row) => sum + (row.item?.points ?? 0), 0) === props.score.total);
const level = computed(() => ({ draft: 'Требует уточнения', workable: 'Можно работать', ready: 'Готова к работе', priority: 'Приоритетная' }[props.score.level.toLowerCase()] || props.score.level));
</script>

<template>
  <div class="readiness" :class="{ 'readiness-compact': compact }">
    <div class="readiness-top"><div><span class="readiness-label">Готовность задачи</span><strong class="readiness-level">{{ level }}</strong></div><span class="score-value">{{ score.total }}<small>/100</small></span></div>
    <div class="weighted-track" role="meter" :aria-valuenow="score.total" aria-valuemin="0" aria-valuemax="100" :aria-valuetext="`${score.total} из 100. ${level}`" aria-label="Готовность задачи">
      <template v-if="validBreakdown"><span v-for="row in rows" :key="row.name" class="weighted-segment" :style="{ flex: row.weight }" aria-hidden="true"><i :style="{ width: `${(row.item!.points / row.weight) * 100}%` }" /></span></template>
      <span v-else class="weighted-segment" style="flex: 1" aria-hidden="true"><i :style="{ width: `${score.total}%` }" /></span>
    </div>
    <dl v-if="detailed && validBreakdown" class="readiness-categories"><div v-for="row in rows" :key="row.name"><dt>{{ row.label }}</dt><dd>{{ row.item!.points }}<span> / {{ row.weight }}</span></dd></div></dl>
    <p v-if="detailed && !validBreakdown" class="readiness-note">Подробная разбивка пока недоступна. Она появится после подтверждения карточки и расчёта оценки.</p>
    <p v-if="detailed" class="readiness-note">Оценка полноты описания, а не качества решения. Низкий балл не мешает публикации.</p>
  </div>
</template>
