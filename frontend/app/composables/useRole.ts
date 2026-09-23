export type Role = 'business' | 'team';

export function useRole() {
  const role = useState<Role>('role', () => 'business');
  const selectedTeamId = useState<string>('selected-team-id', () => '');
  const setRole = (next: Role) => {
    role.value = next;
  };

  return { role, setRole, selectedTeamId };
}
