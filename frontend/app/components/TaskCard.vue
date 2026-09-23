<template>
  <div>
    <!-- Карточка задачи: название, контекст, пользователи, данные,
         ограничения, ожидаемый результат, критерии, контакт -->
  </div>
</template>

<script setup lang="ts">
import type { Task } from '~/types/task';

defineProps<{ task: Task }>();
</script>
