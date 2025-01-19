import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

// @ts-expect-error the above file is migrating to ts
import { useRootStore } from '@/stores/index.js';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { CategoryType } from '@/core/category.ts';
import type { RegisterResponse } from '@/models/auth_response.ts';
import type { User } from '@/models/user.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export function useSignupPageBase() {
    const { tt, getCurrentLanguageTag, getLanguageInfo, setLanguage } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const user = ref<User>(userStore.generateNewUserModel(getCurrentLanguageTag()));
    const submitting = ref<boolean>(false);

    const currentLocale = computed<string>({
        get: () => getCurrentLanguageTag(),
        set: (value: string) => {
            const isCurrencyDefault = user.value.defaultCurrency === settingsStore.localeDefaultSettings.currency;
            const isFirstWeekDayDefault = user.value.firstDayOfWeek === settingsStore.localeDefaultSettings.firstDayOfWeek;

            user.value.language = value;

            const localeDefaultSettings = setLanguage(value);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            if (isCurrencyDefault) {
                user.value.defaultCurrency = settingsStore.localeDefaultSettings.currency;
            }

            if (isFirstWeekDayDefault) {
                user.value.firstDayOfWeek = settingsStore.localeDefaultSettings.firstDayOfWeek;
            }
        },
    });

    const currentLanguageName = computed<string>(() => {
        const languageInfo = getLanguageInfo(currentLocale.value);

        if (!languageInfo) {
            return '';
        }

        return languageInfo.displayName;
    });

    const inputEmptyProblemMessage = computed<string>(() => {
        if (!user.value.username) {
            return 'Username cannot be blank';
        } else if (!user.value.password) {
            return 'Password cannot be blank';
        } else if (!user.value.confirmPassword) {
            return 'Password confirmation cannot be blank';
        } else if (!user.value.email) {
            return 'Email address cannot be blank';
        } else if (!user.value.nickname) {
            return 'Nickname cannot be blank';
        } else if (!user.value.defaultCurrency) {
            return 'Default currency cannot be blank';
        } else {
            return '';
        }
    });

    const inputInvalidProblemMessage = computed<string>(() => {
        if (user.value.password && user.value.confirmPassword && user.value.password !== user.value.confirmPassword) {
            return 'Password and password confirmation do not match';
        } else {
            return '';
        }
    });

    const inputIsEmpty = computed<boolean>(() => !!inputEmptyProblemMessage.value);
    const inputIsInvalid = computed<boolean>(() => !!inputInvalidProblemMessage.value);

    function getCategoryTypeName(categoryType: CategoryType): string {
        switch (categoryType) {
            case CategoryType.Income:
                return tt('Income Categories');
            case CategoryType.Expense:
                return tt('Expense Categories');
            case CategoryType.Transfer:
                return tt('Transfer Categories');
            default:
                return tt('Transaction Categories');
        }
    }

    function doAfterSignupSuccess(response: RegisterResponse): void {
        if (response.user) {
            const localeDefaultSettings = setLanguage(response.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
        }

        if (settingsStore.appSettings.autoUpdateExchangeRatesData) {
            exchangeRatesStore.getLatestExchangeRates({ silent: true, force: false });
        }

        if (response.notificationContent) {
            rootStore.setNotificationContent(response.notificationContent);
        }
    }

    return {
        // states
        user,
        submitting,
        // computed states
        currentLocale,
        currentLanguageName,
        inputEmptyProblemMessage,
        inputInvalidProblemMessage,
        inputIsEmpty,
        inputIsInvalid,
        // functions
        getCategoryTypeName,
        doAfterSignupSuccess
    };
}
