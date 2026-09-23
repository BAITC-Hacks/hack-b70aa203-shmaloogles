<script setup lang="ts">
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
    <a class="skip-link" href="#main-content">Перейти к содержимому</a>
    <button v-if="menuOpen" class="sidebar-overlay" aria-label="Закрыть меню" @click="menuOpen = false" />
    <aside id="main-navigation" ref="navigation" class="sidebar" :class="{ 'is-open': menuOpen }" @keydown="trapMenuFocus">
      <NuxtLink to="/" class="brand" aria-label="Мост — на главную"><BrandMark /><span>мост<span class="brand-dot">.</span></span></NuxtLink>
      <button v-if="menuOpen" class="icon-button sidebar-close" aria-label="Закрыть навигацию" @click="menuOpen = false"><AppIcon name="close" :size="18" /></button>
      <div class="sidebar-caption">Бизнес и студенческие команды</div>
      <div class="nav-label">ПЛАТФОРМА</div>
      <nav class="main-nav" aria-label="Основная навигация">
        <NuxtLink to="/" class="nav-link" exact-active-class="is-active"><AppIcon name="grid" /><span>Обзор</span></NuxtLink>
        <NuxtLink to="/catalog" class="nav-link" active-class="is-active"><AppIcon name="compass" /><span>Каталог задач</span><span class="nav-dot" /></NuxtLink>
        <NuxtLink to="/tasks" class="nav-link" :class="{ 'is-active': route.path.startsWith('/tasks') && route.path !== '/tasks/new' }"><AppIcon name="layers" /><span>Мои задачи</span></NuxtLink>
        <NuxtLink to="/proposals" class="nav-link" active-class="is-active"><AppIcon name="message" /><span>{{ role === 'business' ? 'Предложения команд' : 'Мои предложения' }}</span></NuxtLink>
      </nav>
      <NuxtLink :to="role === 'business' ? '/tasks/new' : '/catalog'" class="button button-primary sidebar-create"><AppIcon :name="role === 'business' ? 'plus' : 'search'" />{{ role === 'business' ? 'Создать задачу' : 'Найти задачу' }}</NuxtLink>
      <div class="sidebar-bottom">
        <div class="sidebar-help">

          <h3>Как начать проект</h3>
          <p>Описание → уточнение →<br>предложение команды.</p>
          <NuxtLink to="/#how-it-works">Как работает Мост <AppIcon name="arrow" :size="17" /></NuxtLink>
        </div>
        <div class="sidebar-status"><span class="status-dot" />Демо-пространство <span class="version">v.01</span></div>
      </div>
    </aside>

    <div class="workspace" :inert="menuOpen || undefined">
      <header class="topbar">
        <div class="topbar-left">
          <button ref="menuButton" class="icon-button mobile-menu" :aria-expanded="menuOpen" aria-controls="main-navigation" :aria-label="menuOpen ? 'Закрыть меню' : 'Открыть меню'" @click="menuOpen = !menuOpen"><AppIcon :name="menuOpen ? 'close' : 'menu'" /></button>
          <div class="breadcrumbs"><span class="breadcrumb-root">Платформа</span><span class="breadcrumb-divider">/</span><span>{{ pageName }}</span></div>
        </div>
        <div class="topbar-right"><span class="role-caption">Я здесь как</span><RoleSwitcher /><span class="profile-avatar" :aria-label="role === 'business' ? 'Роль: бизнес' : 'Роль: команда'">{{ role === 'business' ? 'Б' : 'К' }}<span /></span></div>
      </header>
      <main id="main-content" class="main-content" tabindex="-1"><NuxtPage /></main>
      <footer class="footer"><span>Мост<span class="brand-dot">.</span> Задачи бизнеса. Работа команд.</span><span>Проект начинается с понятной задачи</span></footer>
    </div>
  </div>
</template>
