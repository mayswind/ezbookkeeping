import { computed } from 'vue';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

export interface LanguageSelectBaseProps {
    disabled?: boolean;
    includeSystemDefault?: boolean;
    useModelValue?: boolean;
    modelValue?: string;
}

export interface LanguageSelectBaseEmits {
    (e: 'update:modelValue', value: string): void;
}

export function useLanguageSelectButtonBase(props: LanguageSelectBaseProps, emit: LanguageSelectBaseEmits) {
    const { getCurrentLanguageTag, getCurrentLanguageDisplayName, getAllLanguageOptions, getLanguageInfo, setLanguage } = useI18n();

    const settingsStore = useSettingsStore();

    const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(!!props.includeSystemDefault));

    const currentLocale = computed<string>({
        get: () => getCurrentLanguageTag(),
        set: (value: string) => {
            const localeDefaultSettings = setLanguage(value);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }
    });

    const currentLanguageName = computed<string>(() => {
        if (props.useModelValue && props.modelValue) {
            const languageInfo = getLanguageInfo(props.modelValue);

            if (!languageInfo) {
                return '';
            }

            return languageInfo.displayName;
        } else {
            return getCurrentLanguageDisplayName()
        }
    });

    function updateLanguage(languageTag: string): void {
        if (props.useModelValue) {
            emit('update:modelValue', languageTag);
        } else {
            currentLocale.value = languageTag;
        }
    }

    function isLanguageSelected(languageTag: string): boolean {
        if (props.useModelValue) {
            return props.modelValue === languageTag;
        } else {
            return currentLocale.value === languageTag;
        }
    }

    return {
        // computed states
        allLanguages,
        currentLocale,
        currentLanguageName,
        // functions
        updateLanguage,
        isLanguageSelected
    }
}
