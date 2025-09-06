import fa from './fa.json';

type PersianCalendarLocaleDataKey = 'monthNames' | 'monthShortNames';
type PersianCalendarLocaleData = {
    [K in PersianCalendarLocaleDataKey]: string[];
};

export const DEFAULT_CONTENT: PersianCalendarLocaleData = fa;

export const ALL_LANGUAGES: Record<string, PersianCalendarLocaleData> = {
    'fa': fa
}
