import { translate, type Locale } from '~/i18n/messages';

export function useLocale() {
  const locale = useState<Locale>('ui-locale', () => 'ru');
  const t = (text: string) => translate(locale.value, text);
  return { locale, t };
}
