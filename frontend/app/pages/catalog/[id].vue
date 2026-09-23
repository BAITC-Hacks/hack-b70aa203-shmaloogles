<script setup lang="ts">
const route = useRoute();
const { publishedTasks } = useDemo();
const { role, setRole } = useRole();
const task = computed(() => publishedTasks.value.find(item => item.id === route.params.id));
const sections = computed(() => task.value ? [
  { title: 'Для кого мы это делаем', icon: 'users', value: task.value.card.users },
  { title: 'Данные и материалы', icon: 'layers', value: task.value.card.data },
  { title: 'Что важно учесть', icon: 'shield', value: task.value.card.constraints },
  { title: 'Ожидаемый результат', icon: 'target', value: task.value.card.expectedResult },
  { title: 'Как поймём, что получилось', icon: 'check', value: task.value.card.successCriteria },
  { title: 'Будем на связи', icon: 'message', value: task.value.card.contact },
] : []);
useHead({ title: computed(() => `${task.value?.card.title || 'Задача'} — Мост`) });
</script>

<template>
  <div class="page">
    <NuxtLink to="/catalog" class="back-link"><AppIcon name="back" :size="17" />В каталог задач</NuxtLink>
    <template v-if="task">
      <div class="detail-heading"><div class="detail-company"><span class="company-icon" :class="`tone-${task.color || 'blue'}`"><AppIcon :name="task.icon || 'briefcase'" :size="27" /></span><div><strong>{{ task.organization }}</strong><span>{{ task.category }} <span v-if="task.isExample">· Пример задачи</span></span></div><span class="badge badge-success"><span class="status-dot" />Открыта</span></div><h1>{{ task.card.title }}</h1><div class="tag-list"><span v-for="tag in task.tags" :key="tag" class="tag">{{ tag }}</span></div></div>
      <div class="detail-layout">
        <div class="detail-main"><article class="panel task-full-card"><section class="detail-section"><h2><AppIcon name="briefcase" />О задаче</h2><p>{{ task.card.context || 'Бизнес пока не добавил описание.' }}</p></section><section v-for="section in sections" :key="section.title" class="detail-section"><h2><AppIcon :name="section.icon" />{{ section.title }}</h2><template v-if="section.value.length"><ul v-if="Array.isArray(section.value)" class="detail-list"><li v-for="item in section.value" :key="item">{{ item }}</li></ul><p v-else>{{ section.value }}</p></template><p v-else class="missing-value">Пока не указано. Это можно уточнить у бизнеса.</p></section></article>
          <section id="proposal-form" class="panel proposal-form-panel"><ProposalForm v-if="role === 'team'" :task-id="task.id" /><div v-else class="role-invite"><span class="empty-icon"><AppIcon name="users" :size="26" /></span><h2>Есть идея решения?</h2><p>Переключитесь в роль команды, чтобы рассказать о своём подходе.</p><button class="button button-primary" @click="setRole('team')">Я представляю команду <AppIcon name="arrow" :size="18" /></button></div></section>
        </div>
        <aside class="detail-sidebar"><div class="panel score-panel"><span class="eyebrow muted">ХОРОШАЯ ЗАДАЧА — ЯСНАЯ ЗАДАЧА</span><ReadinessScore :score="task.score" /><p>Рейтинг отражает, насколько подробно бизнес описал задачу и ожидаемый результат.</p><div class="score-explanation"><AppIcon name="info" :size="17" /><span>Откликнуться можно при любом рейтинге.</span></div><a v-if="role === 'team'" href="#proposal-form" class="button button-primary button-full">Предложить решение <AppIcon name="arrow" :size="17" /></a><NuxtLink v-else :to="`/tasks/${task.id}/edit`" class="button button-outline button-full"><AppIcon name="edit" :size="17" />Редактировать задачу</NuxtLink></div><div class="aside-note"><AppIcon name="users" :size="22" /><h3>Начните с диалога</h3><p>Не обязательно знать все ответы. Расскажите о подходе, команде и том, что хотите попробовать.</p></div><div v-if="task.isExample" class="demo-disclaimer"><AppIcon name="info" :size="16" /><span>Это вымышленная задача для знакомства с платформой. Контакты — демонстрационные.</span></div></aside>
      </div>
    </template>
    <EmptyState v-else title="Задача не найдена" description="Возможно, она ещё не опубликована. Другие проекты ждут вас в каталоге." icon="file"><NuxtLink to="/catalog" class="button button-primary">Открыть каталог</NuxtLink></EmptyState>
  </div>
</template>
