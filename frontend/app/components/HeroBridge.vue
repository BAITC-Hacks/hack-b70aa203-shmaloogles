<script setup lang="ts">
import type { Task } from '~/types/task';
defineProps<{ task?: Task }>();
</script>

<template>
  <div class="bridge-study">
    <div class="bridge-study-heading"><span>{{ task ? `Задача № ${task.id} / Из каталога` : 'От задачи к проекту' }}</span><span class="bridge-study-index" aria-hidden="true">МОСТ / 01</span></div>
    <NuxtLink v-if="task" :to="`/catalog/${task.id}`" class="bridge-study-title">{{ task.card.title || 'Задача без названия' }} <span aria-hidden="true">→</span></NuxtLink>
    <h2 v-else class="bridge-study-title">Начните с описания проблемы</h2>

    <svg class="bridge-drawing" viewBox="0 0 480 130" fill="none" aria-hidden="true">
      <path d="M24 116H456M68 22V123M240 8V123M412 22V123" stroke="currentColor" stroke-opacity=".18" />
      <path d="M113 110V67C113 5 237 5 237 67V110M237 110V67C237 5 361 5 361 67V110" stroke="currentColor" stroke-width="19" />
      <path d="M113 110V67C113 5 237 5 237 67V110M237 110V67C237 5 361 5 361 67V110" stroke="white" stroke-opacity=".5" stroke-width="1" />
      <path d="M42 110H438M42 106V114M438 106V114" stroke="currentColor" stroke-width="2" />
      <path d="M68 128H412M68 124V130M412 124V130" stroke="currentColor" stroke-opacity=".35" />
    </svg>

    <div class="bridge-endpoints">
      <div><span>01 / Потребность</span><p>{{ task ? task.card.need || 'Пока не указана' : 'Что бизнесу нужно изменить?' }}</p></div>
      <div><span>02 / Ожидаемый результат</span><p>{{ task ? task.card.expectedResult || 'Пока не указан' : 'Что команда должна подготовить?' }}</p></div>
    </div>
    <ReadinessScore v-if="task" :score="task.score" :breakdown="task.breakdown" compact />
    <NuxtLink v-else to="/tasks/new" class="text-link">Описать задачу →</NuxtLink>
  </div>
</template>
