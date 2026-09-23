<script setup lang="ts">
import type { QuestionDto } from '~/types/api';
import { mapCard } from '~/types/api';
const platform = usePlatformApi();
const step = ref(1); const description = ref(''); const questions = ref<QuestionDto[]>([]); const answers = reactive<Record<string, string>>({});
const loading = ref(false); const error = ref('');
async function clarify() {
  if (description.value.trim().length < 20) return;
  loading.value = true; error.value = '';
  try { questions.value = (await platform.clarify(description.value.trim())).questions; step.value = 2; } catch (cause) { error.value = cause instanceof Error ? cause.message : 'Не удалось получить вопросы'; } finally { loading.value = false; }
}
async function buildCard() {
  loading.value = true; error.value = '';
  try {
    const descriptionText = description.value.trim();
    const generated = await platform.generate(descriptionText, questions.value.map(question => ({ field: question.field, answer: answers[question.field]?.trim() || '' })));
    const task = await platform.createTask(descriptionText);
    await platform.updateTask(task.id, mapCard(generated.card));
    await navigateTo(`/tasks/${task.id}/edit`);
  } catch (cause) { error.value = cause instanceof Error ? cause.message : 'Не удалось создать карточку'; } finally { loading.value = false; }
}
</script>
<template><div><ol class="creation-stepper"><li :class="{ current: step === 1, done: step > 1 }"><span>1</span><div>Ваша идея<small>Расскажите о задаче</small></div></li><li :class="{ current: step === 2 }"><span>2</span><div>Уточнения AI<small>Ответьте на вопросы</small></div></li><li><span>3</span><div>Карточка<small>Проверьте и опубликуйте</small></div></li></ol><div class="form-layout"><div class="panel creation-panel"><div class="assistant-heading"><span class="assistant-avatar"><AppIcon name="sparkles" /></span><div><strong>Поможем разложить всё по полочкам</strong><span>Ассистент Моста</span></div></div><form v-if="step === 1" class="creation-form" @submit.prevent="clarify"><h2>Что вы хотите изменить?</h2><p class="form-intro">Опишите проблему, идею или процесс своими словами.</p><label class="field"><span>Описание задачи <span class="required">*</span></span><textarea v-model="description" required minlength="20" maxlength="20000" rows="9" /></label><p v-if="error" class="error-text" role="alert">{{ error }}</p><div class="form-actions"><button class="button button-primary" :disabled="loading || description.trim().length < 20">{{ loading ? 'Анализируем…' : 'Получить вопросы' }}</button></div></form><form v-else class="creation-form" @submit.prevent="buildCard"><h2>Уточните важные детали</h2><label v-for="(question, index) in questions" :key="question.field" class="field question-field"><span><span class="question-number">{{ String(index + 1).padStart(2, '0') }}</span>{{ question.question }}</span><textarea v-model="answers[question.field]" rows="3" maxlength="20000" /></label><p v-if="error" class="error-text" role="alert">{{ error }}</p><div class="form-actions"><button type="button" class="button button-ghost" @click="step = 1">Назад</button><button class="button button-primary" :disabled="loading">{{ loading ? 'Собираем…' : 'Собрать карточку' }}</button></div></form></div><aside class="form-sidebar"><div class="aside-note"><h3>Не ищите идеальную формулировку.</h3><p>AI задаст вопросы и сохранит отсутствующие сведения пустыми. Всё можно изменить перед публикацией.</p></div></aside></div></div></template>
