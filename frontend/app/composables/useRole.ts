export type Role = 'business' | 'team';

export function useRole() {
  const role = useState<Role>('role', () => 'business');
  const setRole = (next: Role) => {
    role.value = next;
  };

  return { role, setRole };
}
