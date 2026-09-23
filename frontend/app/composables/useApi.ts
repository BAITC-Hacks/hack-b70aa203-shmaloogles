export function useApi() {
  const config = useRuntimeConfig();
  const base = import.meta.server ? config.apiBase : config.public.apiBase;

  const request = async <T>(path: string, init?: RequestInit): Promise<T> => {
    const res = await fetch(`${base}${path}`, init);
    if (!res.ok) {
      let message = `Ошибка API (${res.status})`;
      try { message = (await res.json()).error?.message || message; } catch { /* Keep the status message. */ }
      throw new Error(message);
    }
    return res.json() as Promise<T>;
  };

  return {
    get: <T>(path: string) => request<T>(path),
    post: <T>(path: string, body: unknown) =>
      request<T>(path, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      }),
    patch: <T>(path: string, body: unknown) =>
      request<T>(path, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body),
      }),
    put: <T>(path: string, body: unknown) => request<T>(path, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) }),
  };
}
