<script setup lang="ts">
const { createTask } = useDemo();
const step = ref(1);
const title = ref('');
const description = ref('');
const answers = reactive({ users: '', data: '', result: '' });
const questionHeading = useTemplateRef('questionHeading');
const questions = [
  { key: 'users' as const, number: '01', icon: 'users', title: 'Кому должно стать лучше?', hint: 'Кто будет пользоваться решением? Расскажите об этих людях.', placeholder: 'Например: управляющие трёх кофеен, которые планируют закупки' },
  { key: 'data' as const, number: '02', icon: 'layers', title: 'Что уже есть для начала?', hint: 'Данные, материалы, исследования. Каждый пункт — с новой строки.', placeholder: 'Например: таблица продаж за год\nСписок товаров и закупочные цены' },
  { key: 'result' as const, number: '03', icon: 'target', title: 'Как выглядит хороший результат?', hint: 'Что команда должна передать вам в конце работы?', placeholder: 'Например: дашборд с прогнозом спроса и рекомендациями по закупкам' },
];
async function nextStep() {
  if (title.value.trim().length < 5 || description.value.trim().length < 20) return;
  step.value = 2;
  await nextTick();
  questionHeading.value?.focus();
}
async function buildCard() {
  const task = createTask(description.value.trim(), {
    title: title.value.trim(), context: description.value.trim(), users: toLines(answers.users), data: toLines(answers.data),
    constraints: [], expectedResult: answers.result.trim(), successCriteria: [], contact: '',
  });
  await navigateTo(`/tasks/${task.id}/edit`);
}
</script>

<template>
  <div>
    <ol class="creation-stepper" aria-label="Этапы создания задачи"><li :class="{ current: step === 1, done: step > 1 }" :aria-current="step === 1 ? 'step' : undefined"><span><AppIcon v-if="step > 1" name="check" :size="16" /><template v-else>1</template></span><div>Ваша идея<small>Расскажите о задаче</small></div></li><li :class="{ current: step === 2 }" :aria-current="step === 2 ? 'step' : undefined"><span>2</span><div>Немного деталей<small>Ответьте на вопросы</small></div></li><li><span>3</span><div>Готово к встрече<small>Проверьте и опубликуйте</small></div></li></ol>
    <div class="form-layout">
      <div class="panel creation-panel">
        <div class="assistant-heading"><span class="assistant-avatar"><AppIcon name="sparkles" :size="23" /></span><div><strong>Поможем разложить всё по полочкам</strong><span>Ассистент Моста <span class="badge badge-blue">Демо</span></span></div></div>
        <form v-if="step === 1" class="creation-form" @submit.prevent="nextStep">
          <h2>Что вы хотите изменить?</h2><p class="form-intro">Опишите проблему, идею или процесс, который можно сделать лучше. Простого описания достаточно.</p>
          <label class="field"><span>Рабочее название <span class="required">*</span></span><input v-model="title" required minlength="5" maxlength="120" placeholder="Например: научиться точнее планировать закупки"><small>Коротко о главном. Название можно будет изменить.</small></label>
          <label class="field"><span>Расскажите о вашей задаче <span class="required">*</span></span><textarea v-model="description" required minlength="20" maxlength="4000" rows="7" placeholder="Мы — небольшая сеть кофеен. Сейчас планируем закупки на глаз, и часть продуктов приходится списывать. Хотим понять, сколько закупать, опираясь на данные о продажах…" /><small class="field-bottom"><span>Что происходит сейчас и что хотелось бы улучшить?</span><span>{{ description.length }} / 4000</span></small></label>
          <div class="form-actions"><span class="form-assurance"><AppIcon name="shield" :size="16" />Публикация только после проверки</span><button class="button button-primary" type="submit" :disabled="title.trim().length < 5 || description.trim().length < 20">Продолжить <AppIcon name="arrow" :size="18" /></button></div>
        </form>
        <form v-else class="creation-form" @submit.prevent="buildCard">
          <h2 ref="questionHeading" tabindex="-1">Чуть больше контекста — и мы готовы</h2><p class="form-intro">Три вопроса помогут командам понять вашу задачу. Если ответа пока нет, оставьте поле пустым.</p>
          <div class="description-preview"><AppIcon name="message" :size="17" /><p>{{ description }}</p></div>
          <label v-for="question in questions" :key="question.key" class="field question-field"><span><span class="question-number">{{ question.number }}</span>{{ question.title }}</span><small>{{ question.hint }}</small><textarea v-model="answers[question.key]" rows="3" maxlength="3000" :placeholder="question.placeholder" /></label>
          <div class="form-actions"><button class="button button-ghost" type="button" @click="step = 1"><AppIcon name="back" :size="17" />Назад</button><button class="button button-primary" type="submit">Собрать карточку <AppIcon name="sparkles" :size="18" /></button></div>
        </form>
      </div>
      <aside class="form-sidebar"><div class="aside-note"><span class="note-label">МАЛЕНЬКАЯ ПОДСКАЗКА</span><h3>Не ищите идеальную<br>формулировку.</h3><p>Хорошая задача начинается не с терминов, а с понимания проблемы.</p><ul class="tip-list"><li><AppIcon name="check" :size="16" />Расскажите о своём бизнесе</li><li><AppIcon name="check" :size="16" />Опишите, что не получается</li><li><AppIcon name="check" :size="16" />Поделитесь желаемым результатом</li></ul><div class="note-illustration" aria-hidden="true"><AppIcon name="file" :size="55" /><span><AppIcon name="sparkles" :size="26" /></span></div></div><div class="demo-disclaimer"><AppIcon name="info" :size="17" /><span>Демонстрация интерфейса: вопросы подготовлены заранее. Карточка собирается только из ваших ответов, без обращения к AI.</span></div></aside>
    </div>
  </div>
</template>
