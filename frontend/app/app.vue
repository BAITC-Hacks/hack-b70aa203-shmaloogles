<script setup lang="ts">
const { locale, t } = useLocale();
useHead(() => ({ htmlAttrs: { lang: locale.value } }));

const route = useRoute();
const { role } = useRole();
const menuOpen = ref(false);
const navigation = useTemplateRef<HTMLElement>('navigation');
const menuButton = useTemplateRef<HTMLButtonElement>('menuButton');
const pageName = computed(() => {
  if (route.path === '/') return 'Обзор';
  if (route.path === '/tasks/new') return 'Новая задача';
  if (route.path.includes('/edit')) return 'Редактирование задачи';
  if (route.path.startsWith('/proposals') || route.path.endsWith('/proposals')) return 'Предложения';
  if (route.path.startsWith('/tasks')) return 'Мои задачи';
  return 'Каталог задач';
});
watch(() => route.fullPath, () => { menuOpen.value = false; });
watch(menuOpen, async (open) => {
  if (import.meta.client) document.body.classList.toggle('menu-is-open', open);
  await nextTick();
  if (open) navigation.value?.querySelector<HTMLAnchorElement>('a')?.focus();
  else menuButton.value?.focus();
});
function trapMenuFocus(event: KeyboardEvent) {
  if (!menuOpen.value || event.key !== 'Tab') return;
  const items = navigation.value?.querySelectorAll<HTMLElement>('a[href], button');
  if (!items?.length) return;
  const first = items[0];
  const last = items[items.length - 1];
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault();
    last?.focus();
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault();
    first?.focus();
  }
}
onMounted(() => {
  const desktop = window.matchMedia('(min-width: 768px)');
  const closeOnDesktop = () => { if (desktop.matches) menuOpen.value = false; };
  desktop.addEventListener('change', closeOnDesktop);
  onBeforeUnmount(() => desktop.removeEventListener('change', closeOnDesktop));
});
</script>

<template>
  <div class="app-shell" @keydown.esc="menuOpen = false">
    <NuxtRouteAnnouncer />
    <a class="skip-link" href="#main-content">{{ t("Перейти к содержимому") }}</a>
    <button v-if="menuOpen" class="sidebar-overlay" :aria-label="t('Закрыть меню')" @click="menuOpen = false" />
    <aside id="main-navigation" ref="navigation" class="sidebar" :class="{ 'is-open': menuOpen }" @keydown="trapMenuFocus">
      <NuxtLink to="/" class="brand" :aria-label="t('Мост — на главную')"><BrandMark /><span>{{ t("мост") }}<span class="brand-dot">.</span></span></NuxtLink>
      <button v-if="menuOpen" class="icon-button sidebar-close" :aria-label="t('Закрыть навигацию')" @click="menuOpen = false"><AppIcon name="close" :size="18" /></button>
      <div class="sidebar-caption">{{ t("Идеи встречают возможности") }}</div>
      <div class="nav-label">{{ t("ПЛАТФОРМА") }}</div>
      <nav class="main-nav" :aria-label="t('Основная навигация')">
        <NuxtLink to="/" class="nav-link" exact-active-class="is-active"><AppIcon name="grid" /><span>{{ t("Обзор") }}</span></NuxtLink>
        <NuxtLink to="/catalog" class="nav-link" active-class="is-active"><AppIcon name="compass" /><span>{{ t("Каталог задач") }}</span><span class="nav-dot" /></NuxtLink>
        <NuxtLink to="/tasks" class="nav-link" :class="{ 'is-active': route.path.startsWith('/tasks') && route.path !== '/tasks/new' }"><AppIcon name="layers" /><span>{{ t("Мои задачи") }}</span></NuxtLink>
        <NuxtLink to="/proposals" class="nav-link" active-class="is-active"><AppIcon name="message" /><span>{{ role === 'business' ? t("Предложения команд") : t("Мои предложения") }}</span></NuxtLink>
      </nav>
      <NuxtLink :to="role === 'business' ? '/tasks/new' : '/catalog'" class="button button-primary sidebar-create"><AppIcon :name="role === 'business' ? 'plus' : 'search'" />{{ role === 'business' ? t("Создать задачу") : t("Найти задачу") }}</NuxtLink>
      <div class="sidebar-bottom">
        <div class="sidebar-help">
          <span class="help-spark"><AppIcon name="sparkles" :size="23" /></span>
          <h3>{{ t("У больших идей") }}<br>{{ t("есть начало.") }}</h3>
          <p>{{ t("Расскажите о задаче.") }}<br>{{ t("Всё остальное — по шагам.") }}</p>
          <NuxtLink to="/#how-it-works">{{ t("Как работает Мост") }} <AppIcon name="arrow" :size="17" /></NuxtLink>
        </div>
        <div class="sidebar-status"><span class="status-dot" />{{ t("Демо-пространство") }} <span class="version">v.01</span></div>
      </div>
    </aside>

    <div class="workspace" :inert="menuOpen || undefined">
      <header class="topbar">
        <div class="topbar-left">
          <button ref="menuButton" class="icon-button mobile-menu" :aria-expanded="menuOpen" aria-controls="main-navigation" :aria-label="menuOpen ? t('Закрыть меню') : t('Открыть меню')" @click="menuOpen = !menuOpen"><AppIcon :name="menuOpen ? 'close' : 'menu'" /></button>
          <div class="breadcrumbs"><span class="breadcrumb-root">{{ t("Платформа") }}</span><span class="breadcrumb-divider">/</span><span>{{ t(pageName) }}</span></div>
        </div>
        <div class="topbar-right"><LanguageSwitcher /><span class="role-caption">{{ t("Я здесь как") }}</span><RoleSwitcher /><span class="profile-avatar" :aria-label="role === 'business' ? t('Роль: бизнес') : t('Роль: команда')">{{ role === 'business' ? t("Б") : t("К") }}<span /></span></div>
      </header>
      <main id="main-content" class="main-content" tabindex="-1"><NuxtPage /></main>
      <footer class="footer"><span>{{ t("Мост") }}<span class="brand-dot">.</span> {{ t("Соединяем, чтобы создавать.") }}</span><span>{{ t("Сделано для идей с будущим") }} <AppIcon name="sparkles" :size="14" /></span></footer>
    </div>
  </div>
</template>
