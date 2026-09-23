<template>
  <div class="readiness" :class="[{ 'readiness-compact': compact }, score.total >= 90 ? 'score-priority' : score.total >= 70 ? 'score-ready' : 'score-working']">
    <div class="readiness-top"><span class="readiness-label"><AppIcon :name="score.total >= 90 ? 'bolt' : 'chart'" :size="14" />{{ compact ? readinessLabel(score.level) : t("Готовность задачи") }}</span><span class="score-value">{{ score.total }}<small>/100</small></span></div>
    <div class="score-track" role="meter" :aria-valuenow="score.total" aria-valuemin="0" aria-valuemax="100" :aria-label="t('Рейтинг готовности задачи')"><span :style="{ width: `${score.total}%` }" /></div>
    <div v-if="!compact" class="readiness-bottom"><span>{{ readinessLabel(score.level) }}</span><span>{{ t("Чем яснее задача, тем легче начать") }}</span></div>
  </div>
</template>

<script setup lang="ts">
const { t } = useLocale();

import type { ReadinessScore } from '~/types/task';

defineProps<{ score: ReadinessScore; compact?: boolean }>();
</script>
