<script setup lang="ts">
import type { TaskCard } from '~/types/task';
import { calculateReadiness, isFilled, readinessFields, toLines } from '~/utils/readiness';

const route = useRoute();
const { tasks, updateTask } = useDemo();
const { role, setRole } = useRole();
const task = computed(() => tasks.value.find(item => item.id === route.params.id));
const form = reactive({ title: '', context: '', users: '', data: '', constraints: '', expectedResult: '', successCriteria: '', contact: '' });
const confirmed = ref(false);
const saved = ref(false);
watch(task, (value) => {
  if (!value) return;
  for (const key of Object.keys(form) as (keyof TaskCard)[]) {
    const field = value.card[key];
    form[key] = Array.isArray(field) ? field.join('\n') : field;
  }
}, { immediate: true });
const card = computed<TaskCard>(() => ({
  title: form.title.trim(), context: form.context.trim(), users: toLines(form.users), data: toLines(form.data), constraints: toLines(form.constraints),
  expectedResult: form.expectedResult.trim(), successCriteria: toLines(form.successCriteria), contact: form.contact.trim(),
}));
const score = computed(() => calculateReadiness(card.value));
watch(form, () => { saved.value = false; confirmed.value = false; });
const fields = [
  { key: 'context' as const, label: 'Контекст и потребность', placeholder: 'Что происходит сейчас? Почему важно решить эту задачу?', rows: 4 },
  { key: 'users' as const, label: 'Пользователи', placeholder: 'Кто будет пользоваться решением? Каждый пункт — с новой строки.', rows: 2 },
  { key: 'data' as const, label: 'Данные и материалы', placeholder: 'Что вы готовы предоставить команде? Каждый пункт — с новой строки.', rows: 3 },
  { key: 'constraints' as const, label: 'Ограничения', placeholder: 'Сроки, технологии, бюджет или другие условия. Каждый пункт — с новой строки.', rows: 2 },
  { key: 'expectedResult' as const, label: 'Ожидаемый результат', placeholder: 'Какой конкретный результат должна передать команда?', rows: 3 },
  { key: 'successCriteria' as const, label: 'Критерии успеха', placeholder: 'Как вы оцените результат? Каждый критерий — с новой строки.', rows: 3 },
  { key: 'contact' as const, label: 'Контакт и формат взаимодействия', placeholder: 'Как с вами связаться и как часто вы готовы обсуждать проект?', rows: 2 },
];
function saveDraft() {
  if (!task.value || !card.value.title || role.value !== 'business') return;
  updateTask(task.value.id, card.value);
  saved.value = true;
}
async function publish() {
  if (!task.value || !confirmed.value || role.value !== 'business') return;
  updateTask(task.value.id, card.value, true);
  await navigateTo(`/catalog/${task.value.id}`);
}
useHead({ title: 'Карточка задачи — Мост' });
</script>

<template>
  <div class="page">
    <NuxtLink to="/tasks" class="back-link"><AppIcon name="back" :size="17" />К моим задачам</NuxtLink>
    <template v-if="task && role === 'business'">
      <div class="page-heading"><div><span class="eyebrow muted">ПОЧТИ ГОТОВО К НОВЫМ ВСТРЕЧАМ</span><h1>Ваша идея обрела <span>форму.</span></h1><p>Проверьте карточку, добавьте детали — и дайте командам возможность откликнуться.</p></div><span class="badge" :class="task.status === 'draft' ? 'badge-neutral' : 'badge-success'">{{ task.status === 'draft' ? 'Черновик' : 'Опубликована' }}</span></div>
      <form class="form-layout editor-layout" @submit.prevent="publish">
        <div class="panel editor-panel"><div class="panel-title"><h2><AppIcon name="file" />Карточка задачи</h2><span>Вы управляете каждым словом</span></div><label class="field"><span>Название задачи <span class="required">*</span></span><input v-model="form.title" required minlength="5" maxlength="120" placeholder="Короткое и понятное название"></label><label v-for="field in fields" :key="field.key" class="field"><span>{{ field.label }}<span v-if="!form[field.key].trim()" class="missing-label">Не заполнено</span></span><textarea v-model="form[field.key]" :rows="field.rows" :placeholder="field.placeholder" maxlength="4000" /></label></div>
        <aside class="editor-sidebar"><div class="panel score-panel"><span class="eyebrow muted">ЯСНОСТЬ ПОМОГАЕТ НАЧАТЬ</span><ReadinessScore :score="score" /><div class="score-breakdown"><div v-for="field in readinessFields" :key="field.key" :class="{ filled: isFilled(card[field.key]) }"><span><span class="breakdown-check"><AppIcon v-if="isFilled(card[field.key])" name="check" :size="12" /><span v-else>−</span></span>{{ field.label }}</span><strong>{{ isFilled(card[field.key]) ? field.weight : 0 }}<small>/{{ field.weight }}</small></strong></div></div><p class="score-note">Заполните недостающие поля, чтобы команде было проще предложить решение. Публикация доступна при любом рейтинге.</p></div><div class="panel publish-panel"><h3>Всё выглядит верно?</h3><label class="checkbox-field"><input v-model="confirmed" type="checkbox" required><span>Я проверил(а) информацию и подтверждаю публикацию задачи.</span></label><button class="button button-primary button-full" type="submit" :disabled="!confirmed || form.title.trim().length < 5">{{ task.status === 'published' ? 'Обновить публикацию' : 'Опубликовать задачу' }}<AppIcon name="arrowUp" :size="18" /></button><button class="button button-outline button-full" type="button" :disabled="form.title.trim().length < 5" @click="saveDraft">{{ task.status === 'published' ? 'Сохранить изменения' : 'Сохранить черновик' }}</button><p v-if="saved" class="success-inline" role="status"><AppIcon name="check" :size="16" />Изменения сохранены</p></div><div class="demo-disclaimer"><AppIcon name="info" :size="16" /><span>Деморежим. Изменения сохраняются локально в вашем браузере.</span></div></aside>
      </form>
    </template>
    <EmptyState v-else-if="role !== 'business'" icon="briefcase" title="Редактирование доступно бизнесу" description="Переключитесь в роль бизнеса, чтобы дополнить карточку задачи."><button class="button button-primary" @click="setRole('business')">Перейти в роль бизнеса</button></EmptyState>
    <EmptyState v-else icon="file" title="Задача не найдена" description="Вернитесь к списку или создайте новую задачу."><NuxtLink to="/tasks" class="button button-primary">Мои задачи</NuxtLink></EmptyState>
  </div>
</template>
