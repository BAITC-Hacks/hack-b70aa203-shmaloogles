export type Role = 'business' | 'team';

const role = ref<Role>('business');

export function useRole() {
  const setRole = (next: Role) => {
    role.value = next;
  };

  return { role, setRole };
}
