<script setup lang="ts">
const { t } = useLocale();

const { listTasks } = usePlatformApi();
const { data: publishedTasks, pending, error, refresh: reload } = await useAsyncData('home-tasks', () => listTasks());
const refresh = () => reload();
const { setRole } = useRole();
const featured = computed(() => (publishedTasks.value || []).slice(0, 3));
useHead(() => ({ title: t("Мост — у больших идей есть начало") }));
</script>

<template>
  <div class="page home-page">
    <div class="welcome-line"><span><span class="tiny-star">✳</span> {{ t("Место встречи бизнеса и нового поколения") }}</span><span class="welcome-detail">{{ t("Меньше расстояний. Больше возможностей.") }}</span></div>
    <section class="hero">
      <div class="hero-copy">
        <span class="eyebrow"><span class="status-dot" />{{ t("ОТ ИДЕИ К ОБЩЕМУ ДЕЛУ") }}</span>
        <h1>{{ t("Ваши задачи.") }}<br>{{ t("Их свежий взгляд.") }}<br><span>{{ t("Общий результат.") }}</span></h1>
        <p>{{ t("Соединяем бизнес и студенческие команды,") }}<br class="desktop-break"> {{ t("чтобы хорошие идеи становились реальными проектами.") }}</p>
        <div class="hero-actions"><NuxtLink to="/tasks/new" class="button button-primary" @click="setRole('business')">{{ t("Разместить задачу") }} <AppIcon name="arrowUp" :size="18" /></NuxtLink><NuxtLink to="/catalog" class="button button-white" @click="setRole('team')">{{ t("Найти проект") }} <AppIcon name="arrow" :size="18" /></NuxtLink></div>
        <div class="hero-footnote"><span class="mini-avatars"><span>{{ t("Б") }}</span><span>{{ t("К") }}</span><span><AppIcon name="plus" :size="12" /></span></span><span>{{ t("Разные сильные стороны. Одна цель.") }}</span></div>
      </div>
      <HeroBridge />
    </section>

    <section class="value-strip" :aria-label="t('Возможности платформы')">
      <div class="value-item"><span class="value-icon"><AppIcon name="sparkles" :size="23" /></span><div><strong>{{ t("Из мысли — в задачу") }}</strong><span>{{ t("AI поможет найти нужные слова") }}</span></div></div>
      <div class="value-item"><span class="value-icon"><AppIcon name="target" :size="23" /></span><div><strong>{{ t("Ясность с первого шага") }}</strong><span>{{ t("Понятные цели и рейтинг готовности") }}</span></div></div>
      <div class="value-item"><span class="value-icon"><AppIcon name="users" :size="23" /></span><div><strong>{{ t("Новый взгляд на привычное") }}</strong><span>{{ t("Команды, которым интересно создавать") }}</span></div></div>
    </section>

    <section class="home-projects">
      <div class="section-heading"><div><span class="eyebrow muted">{{ t("РЕАЛЬНЫЙ ОПЫТ НАЧИНАЕТСЯ ЗДЕСЬ") }}</span><h2>{{ t("Задачи со смыслом") }}<span class="heading-dot">.</span></h2><p>{{ t("Найдите то, во что захочется вложить свои знания.") }}</p></div><NuxtLink to="/catalog" class="text-link">{{ t("Весь каталог") }} <AppIcon name="arrow" :size="18" /></NuxtLink></div>
      <div class="demo-caption"><span class="status-dot" />{{ t("Опубликованные задачи") }}</div>
      <div v-if="pending" class="state-panel" role="status"><span class="state-icon"><AppIcon name="compass" /></span><div><strong>{{ t("Загружаем свежие задачи") }}</strong><p>{{ t("Подбираем проекты для знакомства с платформой.") }}</p></div></div>
      <div v-else-if="error" class="state-panel state-error" role="alert"><span class="state-icon"><AppIcon name="info" /></span><div><strong>{{ t("Не удалось загрузить задачи") }}</strong><button class="text-link" @click="refresh">{{ t("Попробовать снова") }}</button></div></div>
      <div v-else-if="featured.length" class="task-grid"><TaskCard v-for="task in featured" :key="task.id" :task="task" /></div>
      <EmptyState v-else :title="t('Опубликованных задач пока нет')" :description="t('Скоро здесь появятся новые проекты от бизнеса.')" />
    </section>

    <section id="how-it-works" class="how-section">
      <div class="section-heading"><div><span class="eyebrow muted">{{ t("ВСЁ ПРОЩЕ, ЧЕМ КАЖЕТСЯ") }}</span><h2>{{ t("Один мост. Три шага.") }}</h2></div><span class="section-aside">{{ t("От «а что, если» до «мы сделали»") }}</span></div>
      <div class="steps-grid">
        <article class="how-step"><div class="step-top"><span>01</span><AppIcon name="message" :size="22" /></div><h3>{{ t("Расскажите о задаче") }}</h3><p>{{ t("Опишите, что хотите изменить. AI-ассистент задаст вопросы и поможет составить понятную карточку.") }}</p><NuxtLink to="/tasks/new" @click="setRole('business')">{{ t("Сформулировать идею") }} <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>02</span><AppIcon name="compass" :size="22" /></div><h3>{{ t("Найдите друг друга") }}</h3><p>{{ t("Бизнес публикует задачу, а команды предлагают свои идеи, план работы и сроки.") }}</p><NuxtLink to="/catalog">{{ t("Посмотреть задачи") }} <AppIcon name="arrow" :size="16" /></NuxtLink></article>
        <article class="how-step"><div class="step-top"><span>03</span><AppIcon name="bolt" :size="22" /></div><h3>{{ t("Создайте что-то важное") }}</h3><p>{{ t("Выберите подходящее предложение и начните сотрудничество. Решение всегда остаётся за бизнесом.") }}</p><NuxtLink to="/proposals">{{ t("К предложениям") }} <AppIcon name="arrow" :size="16" /></NuxtLink></article>
      </div>
    </section>
    <section class="bottom-banner"><div class="banner-symbol"><BrandMark /></div><div><h2>{{ t("Следующая большая идея может быть вашей.") }}</h2><p>{{ t("Начните с небольшой задачи. Дальше — вместе.") }}</p></div><NuxtLink to="/tasks/new" class="button button-primary" @click="setRole('business')">{{ t("Давайте начнём") }} <AppIcon name="arrowUp" :size="18" /></NuxtLink></section>
  </div>
</template>
