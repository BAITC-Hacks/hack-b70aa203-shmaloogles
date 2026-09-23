<script setup lang="ts">
const { listTasks } = usePlatformApi();
const { data: publishedTasks } = await useAsyncData('home-tasks', () => listTasks());
const { setRole } = useRole();
const featured = computed(() => (publishedTasks.value || []).slice(0, 3));
useHead({ title: 'Мост — задачи бизнеса и студенческие команды' });
</script>

<template>
  <div class="page home-page">
    <section class="hero hero-bridge-edition">
      <div class="hero-copy">
        <span class="eyebrow">БИЗНЕС × СТУДЕНЧЕСКИЕ КОМАНДЫ</span>
        <h1>Задачи бизнеса.<br><span>Возможности<br class="hero-title-break"> для команд.</span></h1>
        <p>Опишите проблему и ожидаемый результат. Команды предложат решение, план и сроки.</p>
        <div class="hero-actions"><NuxtLink to="/tasks/new" class="button button-primary" @click="setRole('business')">Создать задачу <AppIcon name="plus" :size="18" /></NuxtLink><NuxtLink to="/catalog" class="button button-white" @click="setRole('team')">Выбрать проект <AppIcon name="arrow" :size="18" /></NuxtLink></div>
      </div>
      <HeroBridge :task="featured[0]" />
    </section>

    <section class="home-projects">
      <div class="section-heading"><div><span class="eyebrow muted">ОТКРЫТЫ К СОТРУДНИЧЕСТВУ</span><h2>Открытые задачи</h2><p>Сравните результат, тему и полноту описания.</p></div><NuxtLink to="/catalog" class="text-link">Весь каталог <AppIcon name="arrow" :size="18" /></NuxtLink></div>

      <div class="task-grid"><TaskCard v-for="task in featured" :key="task.id" :task="task" /></div>
    </section>

    <section id="how-it-works" class="how-section"><BridgeRule />
      <div class="section-heading"><div><span class="eyebrow muted">ПОРЯДОК РАБОТЫ</span><h2>От описания к сотрудничеству</h2></div></div>
      <div class="steps-grid">
        <article class="how-step"><div class="step-top"><span>01</span><AppIcon name="message" :size="22" /></div><h3>Расскажите о задаче</h3><p>Опишите, что хотите изменить. AI-ассистент задаст вопросы и поможет составить понятную карточку.</p><NuxtLink to="/tasks/new" @click="setRole('business')">Описать задачу <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>02</span><AppIcon name="compass" :size="22" /></div><h3>Найдите друг друга</h3><p>Бизнес публикует задачу, а команды предлагают свои идеи, план работы и сроки.</p><NuxtLink to="/catalog">Посмотреть задачи <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>03</span><AppIcon name="bolt" :size="22" /></div><h3>Выберите предложение</h3><p>Выберите подходящее предложение и начните сотрудничество. Решение всегда остаётся за бизнесом.</p><NuxtLink to="/proposals">К предложениям <AppIcon name="arrow" :size="16" /></NuxtLink></article>
      </div>
    </section>

  </div>
</template>
