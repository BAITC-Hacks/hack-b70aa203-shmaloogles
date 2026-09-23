<script setup lang="ts">
const props = defineProps<{ taskId: string }>();
const { submitProposal } = useDemo();
const { role } = useRole();
const sent = ref(false);
const form = reactive({ team: '', idea: '', plan: '', timeframe: '', prototypeUrl: '' });
const urlError = ref('');
function submit() {
  if (role.value !== 'team') return;
  urlError.value = '';
  if (form.prototypeUrl.trim()) {
    try {
      if (!['http:', 'https:'].includes(new URL(form.prototypeUrl.trim()).protocol)) throw new Error('Invalid protocol');
    } catch {
      urlError.value = 'Укажите ссылку, начинающуюся с https:// или http://.';
      return;
    }
  }
  submitProposal({ taskId: props.taskId, team: form.team.trim(), idea: form.idea.trim(), plan: form.plan.trim(), timeframe: form.timeframe.trim(), prototypeUrl: form.prototypeUrl.trim() || undefined });
  sent.value = true;
}
</script>

<template>
  <div v-if="sent" class="proposal-success" role="status"><span class="success-symbol"><AppIcon name="check" :size="30" /></span><span class="eyebrow muted">ПЕРВЫЙ ШАГ СДЕЛАН</span><h2>Предложение отправлено!</h2><p>Теперь слово за бизнесом. Статус решения появится в разделе ваших предложений.</p><NuxtLink to="/proposals" class="button button-primary">Мои предложения <AppIcon name="arrow" :size="17" /></NuxtLink><span class="demo-caption">Предложение сохранено в демо-пространстве</span></div>
  <form v-else class="proposal-form" @submit.prevent="submit">
    <span class="eyebrow muted">ВАШ ПОДХОД МОЖЕТ ВСЁ ИЗМЕНИТЬ</span><h2>Давайте знакомиться.</h2><p class="form-intro">Расскажите, как ваша команда видит решение этой задачи.</p>
    <label class="field"><span>Название команды <span class="required">*</span></span><input v-model="form.team" required pattern=".*\S.*" maxlength="100" placeholder="Как к вам обращаться?"></label>
    <label class="field"><span>Ваша идея <span class="required">*</span></span><textarea v-model="form.idea" required minlength="10" maxlength="4000" rows="4" placeholder="Как вы предлагаете решить задачу? В чём особенность вашего подхода?" /></label>
    <label class="field"><span>Краткий план <span class="required">*</span></span><textarea v-model="form.plan" required minlength="10" maxlength="4000" rows="3" placeholder="С чего начнёте и какие шаги планируете пройти?" /></label>
    <div class="form-two-columns"><label class="field"><span>Предполагаемый срок <span class="required">*</span></span><input v-model="form.timeframe" required pattern=".*\S.*" maxlength="100" placeholder="Например, 3–4 недели"></label><label class="field"><span>Ссылка на прототип <small>необязательно</small></span><input v-model="form.prototypeUrl" type="url" maxlength="2000" placeholder="https://" :aria-invalid="!!urlError" :aria-describedby="urlError ? 'prototype-error' : undefined"><small v-if="urlError" id="prototype-error" class="error-text" role="alert">{{ urlError }}</small></label></div>
    <div class="form-actions"><span class="form-assurance"><AppIcon name="info" :size="16" />Решение принимает бизнес</span><button class="button button-primary" type="submit">Отправить предложение <AppIcon name="send" :size="17" /></button></div>
  </form>
</template>
