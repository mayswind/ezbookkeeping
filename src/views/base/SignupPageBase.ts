import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useUserStore } from '@/stores/user.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import { CategoryType } from '@/core/category.ts';
import type { RegisterResponse } from '@/models/auth_response.ts';
import type { User } from '@/models/user.ts';

import { updateMapCacheExpiration } from '@/lib/cache.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import {
    PASSWORD_MAX_LENGTH,
    PASSWORD_MIN_LENGTH,
    type ValidationProblem,
    isValidEmail,
    isValidPassword
} from '@/lib/validation.ts';

export function useSignupPageBase() {
    const { tt, getCurrentLanguageTag, getLanguageInfo, setLanguage } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const userStore = useUserStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const user = ref<User>(userStore.generateNewUserModel(getCurrentLanguageTag()));
    const submitting = ref<boolean>(false);

    const languageTitle = computed<string>(() => {
        const languageInCurrentLanguage = tt('Language');

        if (languageInCurrentLanguage !== 'Language') {
            return `${languageInCurrentLanguage} / Language`;
        }

        return languageInCurrentLanguage;
    });

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

    function getPasswordProblem(password: string): ValidationProblem | null {
        if (password.length < PASSWORD_MIN_LENGTH) {
            return {
                message: 'parameterizedError.parameter too short',
                options: {
                    parameter: tt('parameter.password'),
                    length: PASSWORD_MIN_LENGTH
                }
            };
        } else if (password.length > PASSWORD_MAX_LENGTH) {
            return {
                message: 'parameterizedError.parameter too long',
                options: {
                    parameter: tt('parameter.password'),
                    length: PASSWORD_MAX_LENGTH
                }
            };
        }

        return null;
    }

    function getProblemMessage(problem: ValidationProblem): string {
        return problem.options ? tt(problem.message, problem.options) : tt(problem.message);
    }

    const inputInvalidProblem = computed<ValidationProblem | null>(() => {
        if (user.value.password && !isValidPassword(user.value.password)) {
            return getPasswordProblem(user.value.password);
        } else if (user.value.password && user.value.confirmPassword && user.value.password !== user.value.confirmPassword) {
            return {
                message: 'Password and password confirmation do not match'
            };
        } else if (user.value.email && !isValidEmail(user.value.email)) {
            return {
                message: 'error.email is invalid'
            };
        } else {
            return null;
        }
    });

    const inputInvalidProblemMessage = computed<string>(() => inputInvalidProblem.value ? getProblemMessage(inputInvalidProblem.value) : '');

    const passwordRules = computed<((value: string) => true | string)[]>(() => [
        (value: string): true | string => {
            const problem = getPasswordProblem(value);
            return !value || !problem || getProblemMessage(problem);
        }
    ]);

    const confirmPasswordRules = computed<((value: string) => true | string)[]>(() => [
        (value: string): true | string => !value || !user.value.password || value === user.value.password || tt('Password and password confirmation do not match')
    ]);

    const emailRules = computed<((value: string) => true | string)[]>(() => [
        (value: string): true | string => !value || isValidEmail(value) || tt('error.email is invalid')
    ]);

    const inputIsEmpty = computed<boolean>(() => !!inputEmptyProblemMessage.value);
    const inputIsInvalid = computed<boolean>(() => !!inputInvalidProblemMessage.value);

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

    function doAfterSignupSuccess(response: RegisterResponse): void {
        if (response.user) {
            const localeDefaultSettings = setLanguage(response.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(response.user.expenseAmountColor, response.user.incomeAmountColor);
        }

        updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
        exchangeRatesStore.removeExpiredExchangeRates(true);
        exchangeRatesStore.autoUpdateExchangeRatesData();

        if (response.notificationContent) {
            rootStore.setNotificationContent(response.notificationContent);
        }
    }

    return {
        // states
        user,
        submitting,
        // computed states
        languageTitle,
        currentLocale,
        currentLanguageName,
        inputEmptyProblemMessage,
        inputInvalidProblem,
        inputInvalidProblemMessage,
        inputIsEmpty,
        inputIsInvalid,
        emailRules,
        passwordRules,
        confirmPasswordRules,
        // functions
        getCategoryTypeName,
        doAfterSignupSuccess
    };
}
