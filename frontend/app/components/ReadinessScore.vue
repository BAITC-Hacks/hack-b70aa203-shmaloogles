<template>
  <div>
    <!-- Рейтинг готовности 0-100 и уровень: Черновик / Рабочая / Готовая / Приоритетная -->
  </div>
</template>

<script setup lang="ts">
import type { ReadinessScore } from '~/types/task';

defineProps<{ score: ReadinessScore }>();
</script>
