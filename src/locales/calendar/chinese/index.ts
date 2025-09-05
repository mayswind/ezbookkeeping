import zhHans from './zh_Hans.json';
import zhHant from './zh_Hant.json';

type ChineseCalendarLocaleDataKey = 'numerals' | 'monthNames' | 'dayNames' | 'leapMonthPrefix' | 'solarTermNames';
type ChineseCalendarLocaleData = {
    [K in ChineseCalendarLocaleDataKey]: K extends 'leapMonthPrefix' ? string : string[];
};

export const DEFAULT_CONTENT: ChineseCalendarLocaleData = zhHans;

export const ALL_LANGUAGES: Record<string, ChineseCalendarLocaleData> = {
    'zh-Hans': zhHans,
    'zh-Hant': zhHant
}
