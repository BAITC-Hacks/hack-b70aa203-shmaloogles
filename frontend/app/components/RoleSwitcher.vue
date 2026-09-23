<template>
  <div class="role-switcher" role="group" aria-label="Выбор демонстрационной роли">
    <button :class="{ selected: role === 'business' }" :aria-pressed="role === 'business'" @click="setRole('business')"><AppIcon name="briefcase" :size="16" />Бизнес</button>
    <button :class="{ selected: role === 'team' }" :aria-pressed="role === 'team'" @click="setRole('team')"><AppIcon name="users" :size="16" />Команда</button>
  </div>
</template>

<script setup lang="ts">
const { role, setRole } = useRole();
</script>
