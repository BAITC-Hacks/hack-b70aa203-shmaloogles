<script setup lang="ts">
withDefaults(defineProps<{ title: string; description: string; icon?: string }>(), { icon: 'search' });
</script>

<template>
  <div class="empty-state"><span class="empty-icon"><AppIcon :name="icon" :size="30" /></span><h2>{{ title }}</h2><p>{{ description }}</p><slot /></div>
</template>
