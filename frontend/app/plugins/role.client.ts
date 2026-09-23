export default defineNuxtPlugin((nuxtApp) => {
  const { role, selectedTeamId } = useRole();
  nuxtApp.hook('app:mounted', () => {
    const savedRole = localStorage.getItem('most-role');
    if (savedRole === 'business' || savedRole === 'team') role.value = savedRole;
    selectedTeamId.value = localStorage.getItem('most-team-id') || '';
    watch(role, value => localStorage.setItem('most-role', value));
    watch(selectedTeamId, value => localStorage.setItem('most-team-id', value));
  });
});
