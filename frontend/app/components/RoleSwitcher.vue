<template>
  <div class="role-switcher" role="group" aria-label="Выбор демонстрационной роли">
    <button :class="{ selected: role === 'business' }" :aria-pressed="role === 'business'" @click="setRole('business')"><AppIcon name="briefcase" :size="16" />Бизнес</button>
    <button :class="{ selected: role === 'team' }" :aria-pressed="role === 'team'" @click="setRole('team')"><AppIcon name="users" :size="16" />Команда</button>
    <select v-if="role === 'team'" v-model="selectedTeamId" aria-label="Выбранная команда"><option value="" disabled>Команда</option><option v-for="team in teams || []" :key="team.id" :value="team.id">{{ team.name }}</option></select>
  </div>
</template>

<script setup lang="ts">
const { role, setRole, selectedTeamId } = useRole();
const { listTeams } = usePlatformApi();
const { data: teams } = await useAsyncData('role-teams', listTeams);
watchEffect(() => { if (!selectedTeamId.value && teams.value?.[0]) selectedTeamId.value = teams.value[0].id; });
</script>
