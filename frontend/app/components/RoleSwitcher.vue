<template>
  <div>
    <label>
      Роль:
      <select v-model="role" @change="setRole(role)">
        <option value="business">Business</option>
        <option value="team">Team</option>
      </select>
    </label>
  </div>
</template>

<script setup lang="ts">
const { role, setRole } = useRole();
</script>
