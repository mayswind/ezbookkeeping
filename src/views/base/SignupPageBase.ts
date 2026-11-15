import { ref, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import type { LocaleDefaultSettings } from '@/core/setting.ts';
import type { RegisterResponse } from '@/models/auth_response.ts';
import type { User } from '@/models/user.ts';

import { isObject, isString } from '@/lib/common.ts';
import { updateMapCacheExpiration, updateTransactionPictureCacheExpiration } from '@/lib/cache.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';

export function useSignupPageBase() {
    const {
        tt,
        getCurrentLanguageTag,
        getLanguageInfo,
        getAllExpenseAmountColors,
        getAllIncomeAmountColors,
        setLanguage
    } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const defaultLocaleSettings = ref<LocaleDefaultSettings>(settingsStore.getLocaleDefaultSettingsSnapshot());
    const user = ref<User>(userStore.generateNewUserModel(getCurrentLanguageTag(), defaultLocaleSettings.value));
    const settingLanguage = ref<boolean>(false);
    const submitting = ref<boolean>(false);

    const allExpenseAmountColorTypes = computed<TypeAndDisplayName[]>(() => getAllExpenseAmountColors());
    const allIncomeAmountColorTypes = computed<TypeAndDisplayName[]>(() => getAllIncomeAmountColors());

    const languageTitle = computed<string>(() => {
        const languageInCurrentLanguage = tt('Language');

        if (languageInCurrentLanguage !== 'Language') {
            return `${languageInCurrentLanguage} / Language`;
        }

        return languageInCurrentLanguage;
    });

    const currentLanguageName = computed<string>(() => {
        const languageInfo = getLanguageInfo(user.value.language);

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

    function updateUserLocaleDefaultSettings(): void {
        const isCurrencyDefault = user.value.defaultCurrency === defaultLocaleSettings.value.currency;
        const isFirstWeekDayDefault = user.value.firstDayOfWeek === defaultLocaleSettings.value.firstDayOfWeek;
        const isExpenseAmountColorDefault = user.value.expenseAmountColor === defaultLocaleSettings.value.preferredExpenseAmountColor?.type;
        const isIncomeAmountColorDefault = user.value.incomeAmountColor === defaultLocaleSettings.value.preferredIncomeAmountColor?.type;

        defaultLocaleSettings.value = settingsStore.getLocaleDefaultSettingsSnapshot();

        if (isCurrencyDefault) {
            user.value.defaultCurrency = settingsStore.localeDefaultSettings.currency;
        }

        if (isFirstWeekDayDefault) {
            user.value.firstDayOfWeek = settingsStore.localeDefaultSettings.firstDayOfWeek;
        }

        if (isExpenseAmountColorDefault && settingsStore.localeDefaultSettings.preferredExpenseAmountColor) {
            user.value.expenseAmountColor = settingsStore.localeDefaultSettings.preferredExpenseAmountColor.type;
        }

        if (isIncomeAmountColorDefault && settingsStore.localeDefaultSettings.preferredIncomeAmountColor) {
            user.value.incomeAmountColor = settingsStore.localeDefaultSettings.preferredIncomeAmountColor.type;
        }
    }

    function updateLanguage(languageTag: string): Promise<void> {
        if (settingLanguage.value) {
            return Promise.resolve();
        }

        settingLanguage.value = true;

        return setLanguage(languageTag).then(localeDefaultSettings => {
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
            user.value.language = languageTag;
            updateUserLocaleDefaultSettings();
        }).finally(() => {
            settingLanguage.value = false;
        });
    }

    function getCategoryTypeName(categoryType: number): string {
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

    function doAfterSignupSuccess(response: RegisterResponse): Promise<string> {
        return new Promise<string>(resolve => {
            const languagePromise = response.user ? setLanguage(response.user.language) : Promise.reject();
            let signupHint: string = '';

            languagePromise.then(localeDefaultSettings => {
                settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);
            }).catch(error => {
                if (isObject(error) && 'message' in error && isString(error.message)) {
                    signupHint = error.message;
                } else if (isString(error)) {
                    signupHint = error;
                }
            }).finally(() => {
                if (response.user) {
                    setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
                }

                updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
                updateTransactionPictureCacheExpiration(settingsStore.appSettings.transactionPictureCacheExpiration);
                exchangeRatesStore.removeExpiredExchangeRates(true);
                exchangeRatesStore.autoUpdateExchangeRatesData();

                if (response.notificationContent) {
                    rootStore.setNotificationContent(response.notificationContent);
                }

                resolve(signupHint);
            });
        });
    }

    watch(() => settingsStore.localeDefaultSettings, () => {
        user.value.language = getCurrentLanguageTag();
        updateUserLocaleDefaultSettings();
    }, {
        deep: true
    });

    return {
        // states
        user,
        settingLanguage,
        submitting,
        // computed states
        allExpenseAmountColorTypes,
        allIncomeAmountColorTypes,
        languageTitle,
        currentLanguageName,
        inputEmptyProblemMessage,
        inputInvalidProblemMessage,
        inputIsEmpty,
        inputIsInvalid,
        // functions
        updateLanguage,
        getCategoryTypeName,
        doAfterSignupSuccess
    };
}
