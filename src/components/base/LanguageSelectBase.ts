import { ref, computed } from 'vue';

import type { LanguageOption } from '@/locales/index.ts';
import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

export interface LanguageSelectBaseProps {
    disabled?: boolean;
    loading?: boolean;
    includeSystemDefault?: boolean;
    useModelValue?: boolean;
    modelValue?: string;
}

export interface LanguageSelectBaseEmits {
    (e: 'update:modelValue', value: string): void;
    (e: 'error', message: string): void;
}

export function useLanguageSelectButtonBase(props: LanguageSelectBaseProps, emit: LanguageSelectBaseEmits) {
    const { getCurrentLanguageTag, getCurrentLanguageDisplayName, getAllLanguageOptions, getLanguageInfo, setLanguage } = useI18n();

    const settingsStore = useSettingsStore();

    const settingLanguage = ref<boolean>(false);

    const allLanguages = computed<LanguageOption[]>(() => getAllLanguageOptions(!!props.includeSystemDefault));

    const currentLocale = computed<string>(() => getCurrentLanguageTag());

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

    function updateLanguage(languageTag: string): Promise<void> {
        if (props.disabled || props.loading || settingLanguage.value) {
            return Promise.resolve();
        }

        if (props.useModelValue) {
            emit('update:modelValue', languageTag);
            return Promise.resolve();
        }

        settingLanguage.value = true;

        return setLanguage(languageTag).then(localeDefaultSettings => {
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
        }).finally(() => {
            settingLanguage.value = false;
        });
    }

    function isLanguageSelected(languageTag: string): boolean {
        if (props.useModelValue) {
            return props.modelValue === languageTag;
        } else {
            return currentLocale.value === languageTag;
        }
    }

    return {
        // states,
        settingLanguage,
        // computed states
        allLanguages,
        currentLocale,
        currentLanguageName,
        // functions
        updateLanguage,
        isLanguageSelected
    }
}
