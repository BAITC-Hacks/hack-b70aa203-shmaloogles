<template>
  <div class="role-controls">
    <div class="role-switcher" role="group" :aria-label="t('Выбор демонстрационной роли')">
      <button :aria-label="t('Бизнес')" :class="{ selected: role === 'business' }" :aria-pressed="role === 'business'" @click="setRole('business')"><AppIcon name="briefcase" :size="16" /><span>{{ t("Бизнес") }}</span></button>
      <button :aria-label="t('Команда')" :class="{ selected: role === 'team' }" :aria-pressed="role === 'team'" @click="setRole('team')"><AppIcon name="users" :size="16" /><span>{{ t("Команда") }}</span></button>
    </div>
    <label v-if="role === 'team'" class="team-picker" :class="{ 'has-error': error }">
      <AppIcon name="users" :size="15" /><span class="sr-only">{{ t("Выбранная команда") }}</span>
      <select v-model="selectedTeamId" :disabled="pending || !teams?.length" :aria-label="t('Выбранная команда')">
        <option value="" disabled>{{ pending ? t("Загрузка…") : error ? t("Ошибка загрузки") : teams?.length ? t("Выберите команду") : t("Нет команд") }}</option>
        <option v-for="team in teams || []" :key="team.id" :value="team.id">{{ team.name }}</option>
      </select><AppIcon name="down" :size="13" />
    </label>
  </div>
</template>

<script setup lang="ts">
const { t } = useLocale();

const { role, setRole, selectedTeamId } = useRole();
const { listTeams } = usePlatformApi();
const { data: teams, pending, error } = await useAsyncData('role-teams', listTeams);
watchEffect(() => { if (!selectedTeamId.value && teams.value?.[0]) selectedTeamId.value = teams.value[0].id; });
</script>
