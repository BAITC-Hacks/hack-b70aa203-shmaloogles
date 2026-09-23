export default defineNuxtPlugin(() => {
  const saved = useCookie<string>('most-locale', { maxAge: 60 * 60 * 24 * 365, sameSite: 'lax', path: '/' });
  const { locale } = useLocale();
  locale.value = saved.value === 'kk' || saved.value === 'en' ? saved.value : 'ru';
  watch(locale, value => { saved.value = value; });
});
