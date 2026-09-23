<template>
  <div class="role-controls">
    <div class="role-switcher" role="group" aria-label="Выбор демонстрационной роли">
      <button aria-label="Бизнес" :class="{ selected: role === 'business' }" :aria-pressed="role === 'business'" @click="setRole('business')"><AppIcon name="briefcase" :size="16" /><span>Бизнес</span></button>
      <button aria-label="Команда" :class="{ selected: role === 'team' }" :aria-pressed="role === 'team'" @click="setRole('team')"><AppIcon name="users" :size="16" /><span>Команда</span></button>
    </div>
    <label v-if="role === 'team'" class="team-picker" :class="{ 'has-error': error }">
      <AppIcon name="users" :size="15" /><span class="sr-only">Выбранная команда</span>
      <select v-model="selectedTeamId" :disabled="pending || !teams?.length" aria-label="Выбранная команда">
        <option value="" disabled>{{ pending ? 'Загрузка…' : error ? 'Ошибка загрузки' : teams?.length ? 'Выберите команду' : 'Нет команд' }}</option>
        <option v-for="team in teams || []" :key="team.id" :value="team.id">{{ team.name }}</option>
      </select><AppIcon name="down" :size="13" />
    </label>
  </div>
</template>

<script setup lang="ts">
const { role, setRole, selectedTeamId } = useRole();
const { listTeams } = usePlatformApi();
const { data: teams, pending, error } = await useAsyncData('role-teams', listTeams);
watchEffect(() => { if (!selectedTeamId.value && teams.value?.[0]) selectedTeamId.value = teams.value[0].id; });
</script>
