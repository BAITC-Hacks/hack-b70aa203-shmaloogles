<script setup lang="ts">
withDefaults(defineProps<{ name: string; size?: number }>(), { size: 20 });

const paths: Record<string, string> = {
  grid: 'M3 3h7v7H3z M14 3h7v7h-7z M3 14h7v7H3z M14 14h7v7h-7z',
  compass: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M16 8l-2.5 5.5L8 16l2.5-5.5z',
  briefcase: 'M8 7V5a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2 M3 7h18v13H3z M3 11l9 4 9-4 M10 14h4',
  layers: 'M12 3 2 8l10 5 10-5-10-5z M2 12l10 5 10-5 M2 16l10 5 10-5',
  message: 'M21 11.5a8.5 8.5 0 0 1-8.5 8.5H4l-3 2 2-6a8.5 8.5 0 1 1 18-4.5z M7 10h10 M7 14h6',
  arrow: 'M4 12h16 M14 6l6 6-6 6',
  arrowUp: 'M6 18 18 6 M6 6h12v12',
  chevron: 'm9 5 7 7-7 7',
  down: 'm6 9 6 6 6-6',
  back: 'M20 12H4 M10 6l-6 6 6 6',
  plus: 'M12 5v14 M5 12h14',
  check: 'm5 12 4 4L19 6',
  close: 'm6 6 12 12 M6 18 18 6',
  search: 'M21 21l-5-5 M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16',
  filter: 'M4 7h16 M7 12h10 M10 17h4',
  sparkles: 'm12 3 2.5 6.5L21 12l-6.5 2.5L12 21l-2.5-6.5L3 12l6.5-2.5L12 3z M20 2v4 M18 4h4',
  users: 'M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2 M9 11a4 4 0 1 0 0-8 4 4 0 0 0 0 8 M22 21v-2a4 4 0 0 0-3-3.87 M16 3.13a4 4 0 0 1 0 7.75',
  building: 'M4 21V3h12v18 M16 9h4v12 M2 21h20 M8 7h4 M8 11h4 M8 15h4 M9 21v-3h2v3',
  clock: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M12 6v6l4 2',
  info: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M12 11v6 M12 7h.01',
  help: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M9.1 9a3 3 0 0 1 5.8 1c0 2-3 2-3 4 M12 17h.01',
  code: 'm8 6-6 6 6 6 M16 6l6 6-6 6 M14 3l-4 18',
  chart: 'M4 3v18h17 M8 16v-5 M13 16V7 M18 16v-8',
  leaf: 'M20 3C9 2 3 7 4 14c1 6 9 7 13 2 3-4 2-8 3-13z M4 21l10-11',
  coffee: 'M4 8h12v9a4 4 0 0 1-4 4H8a4 4 0 0 1-4-4V8z M16 8h2a3 3 0 0 1 0 6h-2 M7 2v3 M12 2v3',
  heart: 'M20.8 4.6a5.5 5.5 0 0 0-7.8 0L12 5.7l-1.1-1.1a5.5 5.5 0 0 0-7.8 7.8L12 21l8.8-8.6a5.5 5.5 0 0 0 0-7.8z',
  book: 'M12 5c-3-2-7-2-10-1v16c3-1 7-1 10 1 3-2 7-2 10-1V4c-3-1-7-1-10 1z M12 5v16',
  globe: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M2 12h20 M12 2c6 6 6 14 0 20-6-6-6-14 0-20',
  target: 'M12 22a10 10 0 1 0 0-20 10 10 0 0 0 0 20 M12 18a6 6 0 1 0 0-12 6 6 0 0 0 0 12 M12 14a2 2 0 1 0 0-4 2 2 0 0 0 0 4',
  file: 'M14 2H4v20h16V8l-6-6z M14 2v6h6 M8 13h8 M8 17h5',
  link: 'M10 13a5 5 0 0 0 7 0l3-3a5 5 0 0 0-7-7l-2 2 M14 11a5 5 0 0 0-7 0l-3 3a5 5 0 0 0 7 7l2-2',
  send: 'm22 2-7 20-4-9-9-4 20-7z M22 2 11 13',
  edit: 'm16 3 5 5 M3 21l5-1L21 7a3.5 3.5 0 0 0-5-5L3 15v6z',
  menu: 'M4 6h16 M4 12h16 M4 18h16',
  bolt: 'm13 2-9 12h7l-1 8 10-12h-7l1-8z',
  shield: 'm12 2 9 4v6c0 6-9 10-9 10S3 18 3 12V6l9-4z M8 12l3 3 5-6',
};
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.65" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true" class="icon">
    <path :d="paths[name] || paths.layers" />
  </svg>
</template>
