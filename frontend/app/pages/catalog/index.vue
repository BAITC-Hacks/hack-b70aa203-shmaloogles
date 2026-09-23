<script setup lang="ts">
const { publishedTasks } = useDemo();
const { role } = useRole();
const search = ref('');
const category = ref('Все задачи');
const readiness = ref('all');
const sort = ref('rating');
const view = ref('grid');
const categories = ['Все задачи', 'Разработка', 'Дизайн', 'Аналитика'];
const filtered = computed(() => {
  const query = search.value.trim().toLocaleLowerCase('ru');
  const result = publishedTasks.value.filter(task =>
    (category.value === 'Все задачи' || task.category === category.value)
    && (readiness.value === 'all' || (readiness.value === 'priority' ? task.score.total >= 90 : readiness.value === 'ready' ? task.score.total >= 70 && task.score.total < 90 : task.score.total < 70))
    && (!query || [task.card.title, task.description, task.organization, ...(task.tags || [])].join(' ').toLocaleLowerCase('ru').includes(query)),
  );
  return sort.value === 'rating' ? result.sort((a, b) => b.score.total - a.score.total) : sort.value === 'rating-asc' ? result.sort((a, b) => a.score.total - b.score.total) : result;
});
function resetFilters() { search.value = ''; category.value = 'Все задачи'; readiness.value = 'all'; }
useHead({ title: 'Каталог задач — Мост' });
</script>

<template>
  <div class="page">
    <div class="page-heading"><div><span class="eyebrow muted">НАЙДИТЕ СВОЮ ТОЧКУ ПРИЛОЖЕНИЯ</span><h1>Задачи, которые ждут <span>вас.</span></h1><p>Настоящие вызовы бизнеса. Пространство для ваших идей.</p></div><NuxtLink v-if="role === 'business'" to="/tasks/new" class="button button-primary"><AppIcon name="plus" :size="18" />Создать задачу</NuxtLink></div>
    <div class="catalog-callout"><span class="callout-icon"><AppIcon name="compass" :size="25" /></span><div><strong>Не просто строчка в портфолио.</strong><p>Выбирайте задачу, предлагайте решение и создавайте то, чем будут пользоваться.</p></div><span class="callout-decoration" aria-hidden="true">↗</span></div>
    <div class="catalog-toolbar"><label class="search-field"><AppIcon name="search" :size="20" /><input v-model="search" type="search" placeholder="Название, навык или компания" aria-label="Поиск задач"></label><label class="select-field"><AppIcon name="filter" :size="18" /><select v-model="readiness" aria-label="Уровень готовности"><option value="all">Любая готовность</option><option value="priority">Приоритетная · 90–100</option><option value="ready">Готовая · 70–89</option><option value="working">В работе · 0–69</option></select></label></div>
    <div class="catalog-tabs-row"><div class="filter-tabs" role="group" aria-label="Направление задачи"><button v-for="item in categories" :key="item" :class="{ active: category === item }" :aria-pressed="category === item" @click="category = item">{{ item }}<span v-if="item === 'Все задачи'">{{ publishedTasks.length }}</span></button></div><span class="demo-caption">Демокаталог</span></div>
    <div class="results-toolbar"><span aria-live="polite">Найдено задач: <strong>{{ filtered.length }}</strong></span><div class="results-controls"><label class="sort-label"><span>Сначала</span><select v-model="sort" aria-label="Сортировка задач"><option value="rating">с высоким рейтингом</option><option value="rating-asc">с низким рейтингом</option><option value="new">новые</option></select></label><div class="view-switcher" role="group" aria-label="Вид каталога"><button :class="{ active: view === 'grid' }" :aria-pressed="view === 'grid'" aria-label="Сетка" @click="view = 'grid'"><AppIcon name="grid" :size="17" /></button><button :class="{ active: view === 'list' }" :aria-pressed="view === 'list'" aria-label="Список" @click="view = 'list'"><AppIcon name="menu" :size="18" /></button></div></div></div>
    <div v-if="filtered.length" class="task-grid" :class="{ 'list-view': view === 'list' }"><TaskCard v-for="task in filtered" :key="task.id" :task="task" /></div>
    <EmptyState v-else title="Пока ничего не нашлось" description="Попробуйте другое название или уберите несколько фильтров."><button class="button button-primary" @click="resetFilters">Сбросить фильтры</button></EmptyState>
    <div class="catalog-bottom-note"><AppIcon name="info" :size="17" /><p>Рейтинг показывает полноту описания, а не сложность задачи. Предложить решение можно при любом рейтинге.</p></div>
  </div>
</template>
