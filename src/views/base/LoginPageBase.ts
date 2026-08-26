import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useRootStore } from '@/stores/index.ts';
import { useSettingsStore } from '@/stores/setting.ts';
import { useExchangeRatesStore } from '@/stores/exchangeRates.ts';

import type { AuthResponse } from '@/models/auth_response.ts';

import { updateMapCacheExpiration } from '@/lib/cache.ts';
import { getOAuth2Provider, getOIDCCustomDisplayNames, getLoginPageTips } from '@/lib/server_settings.ts';
import { getClientDisplayVersion } from '@/lib/version.ts';
import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import {
    PASSWORD_MAX_LENGTH,
    PASSWORD_MIN_LENGTH,
    type ValidationProblem,
    isValidPassword
} from '@/lib/validation.ts';

export function useLoginPageBase(platform: 'mobile' | 'desktop') {
    const { tt, getServerMultiLanguageConfigContent, getLocalizedOAuth2LoginText, setLanguage } = useI18n();

    const rootStore = useRootStore();
    const settingsStore = useSettingsStore();
    const exchangeRatesStore = useExchangeRatesStore();

    const version = `${getClientDisplayVersion()}`;

    const username = ref<string>('');
    const password = ref<string>('');
    const passcode = ref<string>('');
    const backupCode = ref<string>('');
    const tempToken = ref<string>('');
    const twoFAVerifyType = ref<string>('passcode');
    const oauth2ClientSessionId = ref<string>('');

    const loggingInByPassword = ref<boolean>(false);
    const loggingInByOAuth2 = ref<boolean>(false);
    const verifying = ref<boolean>(false);

    function getPasswordProblem(passwordValue: string): ValidationProblem | null {
        if (passwordValue.length < PASSWORD_MIN_LENGTH) {
            return {
                message: 'parameterizedError.parameter too short',
                options: {
                    parameter: tt('parameter.password'),
                    length: PASSWORD_MIN_LENGTH
                }
            };
        } else if (passwordValue.length > PASSWORD_MAX_LENGTH) {
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

    const inputInvalidProblem = computed<ValidationProblem | null>(() => password.value && !isValidPassword(password.value) ? getPasswordProblem(password.value) : null);
    const inputInvalidProblemMessage = computed<string>(() => inputInvalidProblem.value ? getProblemMessage(inputInvalidProblem.value) : '');
    const passwordRules = computed<((value: string) => true | string)[]>(() => [
        (value: string): true | string => {
            const problem = getPasswordProblem(value);
            return !value || !problem || getProblemMessage(problem);
        }
    ]);
    const inputIsEmpty = computed<boolean>(() => !username.value || !password.value);
    const inputIsInvalid = computed<boolean>(() => !!inputInvalidProblemMessage.value);
    const twoFAInputIsEmpty = computed<boolean>(() => {
        if (twoFAVerifyType.value === 'backupcode') {
            return !backupCode.value;
        } else {
            return !passcode.value;
        }
    });

    const oauth2LoginUrl = computed<string>(() => rootStore.generateOAuth2LoginUrl(platform, oauth2ClientSessionId.value));
    const oauth2LoginDisplayName = computed<string>(() => getLocalizedOAuth2LoginText(getOAuth2Provider(), getOIDCCustomDisplayNames()));
    const tips = computed<string>(() => getServerMultiLanguageConfigContent(getLoginPageTips()));

    function doAfterLogin(authResponse: AuthResponse): void {
        if (authResponse.user) {
            const localeDefaultSettings = setLanguage(authResponse.user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(authResponse.user.expenseAmountColor, authResponse.user.incomeAmountColor);
        }

        updateMapCacheExpiration(settingsStore.appSettings.mapCacheExpiration);
        exchangeRatesStore.removeExpiredExchangeRates(true);
        exchangeRatesStore.autoUpdateExchangeRatesData();

        if (authResponse.notificationContent) {
            rootStore.setNotificationContent(authResponse.notificationContent);
        }
    }

    return {
        // constants
        version,
        // states
        username,
        password,
        passcode,
        backupCode,
        tempToken,
        twoFAVerifyType,
        oauth2ClientSessionId,
        loggingInByPassword,
        loggingInByOAuth2,
        verifying,
        // computed states
        inputIsEmpty,
        inputIsInvalid,
        inputInvalidProblem,
        inputInvalidProblemMessage,
        passwordRules,
        twoFAInputIsEmpty,
        oauth2LoginUrl,
        oauth2LoginDisplayName,
        tips,
        // functions
        doAfterLogin
    }
}
