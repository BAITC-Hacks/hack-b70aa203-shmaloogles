<script setup lang="ts">
const { listTasks } = usePlatformApi();
const { data: publishedTasks } = await useAsyncData('home-tasks', () => listTasks());
const { setRole } = useRole();
const featured = computed(() => (publishedTasks.value || []).slice(0, 3));
useHead({ title: 'Мост — у больших идей есть начало' });
</script>

<template>
  <div class="page home-page">
    <div class="welcome-line"><span><span class="tiny-star">✳</span> Место встречи бизнеса и нового поколения</span><span class="welcome-detail">Меньше расстояний. Больше возможностей.</span></div>
    <section class="hero">
      <div class="hero-copy">
        <span class="eyebrow"><span class="status-dot" />ОТ ИДЕИ К ОБЩЕМУ ДЕЛУ</span>
        <h1>Ваши задачи.<br>Их свежий взгляд.<br><span>Общий результат.</span></h1>
        <p>Соединяем бизнес и студенческие команды,<br class="desktop-break"> чтобы хорошие идеи становились реальными проектами.</p>
        <div class="hero-actions"><NuxtLink to="/tasks/new" class="button button-primary" @click="setRole('business')">Разместить задачу <AppIcon name="arrowUp" :size="18" /></NuxtLink><NuxtLink to="/catalog" class="button button-white" @click="setRole('team')">Найти проект <AppIcon name="arrow" :size="18" /></NuxtLink></div>
        <div class="hero-footnote"><span class="mini-avatars"><span>Б</span><span>К</span><span><AppIcon name="plus" :size="12" /></span></span><span>Разные сильные стороны. Одна цель.</span></div>
      </div>
      <HeroBridge />
    </section>

    <section class="value-strip" aria-label="Возможности платформы">
      <div class="value-item"><span class="value-icon"><AppIcon name="sparkles" :size="23" /></span><div><strong>Из мысли — в задачу</strong><span>AI поможет найти нужные слова</span></div></div>
      <div class="value-item"><span class="value-icon"><AppIcon name="target" :size="23" /></span><div><strong>Ясность с первого шага</strong><span>Понятные цели и рейтинг готовности</span></div></div>
      <div class="value-item"><span class="value-icon"><AppIcon name="users" :size="23" /></span><div><strong>Новый взгляд на привычное</strong><span>Команды, которым интересно создавать</span></div></div>
    </section>

    <section class="home-projects">
      <div class="section-heading"><div><span class="eyebrow muted">РЕАЛЬНЫЙ ОПЫТ НАЧИНАЕТСЯ ЗДЕСЬ</span><h2>Задачи со смыслом<span class="heading-dot">.</span></h2><p>Найдите то, во что захочется вложить свои знания.</p></div><NuxtLink to="/catalog" class="text-link">Весь каталог <AppIcon name="arrow" :size="18" /></NuxtLink></div>
      <div class="demo-caption"><span class="status-dot" />Опубликованные задачи</div>
      <div class="task-grid"><TaskCard v-for="task in featured" :key="task.id" :task="task" /></div>
    </section>

    <section id="how-it-works" class="how-section">
      <div class="section-heading"><div><span class="eyebrow muted">ВСЁ ПРОЩЕ, ЧЕМ КАЖЕТСЯ</span><h2>Один мост. Три шага.</h2></div><span class="section-aside">От «а что, если» до «мы сделали»</span></div>
      <div class="steps-grid">
        <article class="how-step"><div class="step-top"><span>01</span><AppIcon name="message" :size="22" /></div><h3>Расскажите о задаче</h3><p>Опишите, что хотите изменить. AI-ассистент задаст вопросы и поможет составить понятную карточку.</p><NuxtLink to="/tasks/new" @click="setRole('business')">Сформулировать идею <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>02</span><AppIcon name="compass" :size="22" /></div><h3>Найдите друг друга</h3><p>Бизнес публикует задачу, а команды предлагают свои идеи, план работы и сроки.</p><NuxtLink to="/catalog">Посмотреть задачи <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>03</span><AppIcon name="bolt" :size="22" /></div><h3>Создайте что-то важное</h3><p>Выберите подходящее предложение и начните сотрудничество. Решение всегда остаётся за бизнесом.</p><NuxtLink to="/proposals">К предложениям <AppIcon name="arrow" :size="16" /></NuxtLink></article>
      </div>
    </section>
    <section class="bottom-banner"><div class="banner-symbol"><BrandMark /></div><div><h2>Следующая большая идея может быть вашей.</h2><p>Начните с небольшой задачи. Дальше — вместе.</p></div><NuxtLink to="/tasks/new" class="button button-primary" @click="setRole('business')">Давайте начнём <AppIcon name="arrowUp" :size="18" /></NuxtLink></section>
  </div>
</template>
