import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useOverviewStore } from '@/stores/overview.ts';

import type { TypeAndDisplayName } from '@/core/base.ts';
import { WeekDay } from '@/core/datetime.ts';
import type { LocalizedDigitGroupingType } from '@/core/numeral.ts';

import { type UserBasicInfo, User } from '@/models/user.ts';
import { type CategorizedAccount, Account} from '@/models/account.ts';

import { setExpenseAndIncomeAmountColor } from '@/lib/ui/common.ts';
import { getCategorizedAccounts } from '@/lib/account.ts';

export function useUserProfilePageBase() {
    const {
        tt,
        getDefaultCurrency,
        getDefaultFirstDayOfWeek,
        getAllWeekDays,
        getAllLongDateFormats,
        getAllShortDateFormats,
        getAllLongTimeFormats,
        getAllShortTimeFormats,
        getAllDecimalSeparators,
        getAllDigitGroupingSymbols,
        getAllDigitGroupingTypes,
        getAllCurrencyDisplayTypes,
        getAllExpenseAmountColors,
        getAllIncomeAmountColors,
        getAllTransactionEditScopeTypes,
        setLanguage
    } = useI18n();

    const settingsStore = useSettingsStore();
    const accountsStore = useAccountsStore();
    const overviewStore = useOverviewStore();

    const defaultFirstDayOfWeekName = getDefaultFirstDayOfWeek();
    const defaultFirstDayOfWeek = WeekDay.parse(defaultFirstDayOfWeekName) ? (WeekDay.parse(defaultFirstDayOfWeekName) as WeekDay).type : WeekDay.DefaultFirstDay.type;

    const newProfile = ref<User>(User.createNewUser('', getDefaultCurrency(), defaultFirstDayOfWeek));
    const oldProfile = ref<User>(User.createNewUser('', getDefaultCurrency(), defaultFirstDayOfWeek));

    const emailVerified = ref<boolean>(false);
    const loading = ref<boolean>(false);
    const resending = ref<boolean>(false);
    const saving = ref<boolean>(false);

    const allAccounts = computed<Account[]>(() => accountsStore.allPlainAccounts);
    const allVisibleAccounts = computed<Account[]>(() => accountsStore.allVisiblePlainAccounts);
    const allVisibleCategorizedAccounts = computed<CategorizedAccount[]>(() => getCategorizedAccounts(allVisibleAccounts.value));
    const allWeekDays = computed<TypeAndDisplayName[]>(() => getAllWeekDays());
    const allLongDateFormats = computed<TypeAndDisplayName[]>(() => getAllLongDateFormats());
    const allShortDateFormats = computed<TypeAndDisplayName[]>(() => getAllShortDateFormats());
    const allLongTimeFormats = computed<TypeAndDisplayName[]>(() => getAllLongTimeFormats());
    const allShortTimeFormats = computed<TypeAndDisplayName[]>(() => getAllShortTimeFormats());
    const allDecimalSeparators = computed<TypeAndDisplayName[]>(() => getAllDecimalSeparators());
    const allDigitGroupingSymbols = computed<TypeAndDisplayName[]>(() => getAllDigitGroupingSymbols());
    const allDigitGroupingTypes = computed<LocalizedDigitGroupingType[]>(() => getAllDigitGroupingTypes());
    const allCurrencyDisplayTypes = computed<TypeAndDisplayName[]>(() => getAllCurrencyDisplayTypes());
    const allExpenseAmountColorTypes = computed<TypeAndDisplayName[]>(() => getAllExpenseAmountColors());
    const allIncomeAmountColorTypes = computed<TypeAndDisplayName[]>(() => getAllIncomeAmountColors());
    const allTransactionEditScopeTypes = computed<TypeAndDisplayName[]>(() => getAllTransactionEditScopeTypes());

    const languageTitle = computed<string>(() => {
        const languageInCurrentLanguage = tt('Language');

        if (languageInCurrentLanguage !== 'Language') {
            return `${languageInCurrentLanguage} / Language`;
        }

        return languageInCurrentLanguage;
    });

    const supportDigitGroupingSymbol = computed<boolean>(() => {
        for (const digitGroupingType of allDigitGroupingTypes.value) {
            if (digitGroupingType.type === newProfile.value.digitGrouping) {
                return digitGroupingType.enabled;
            }
        }

        return false;
    });

    const inputIsNotChangedProblemMessage = computed<string | null>(() => {
        if (!newProfile.value.password && !newProfile.value.confirmPassword && !newProfile.value.email && !newProfile.value.nickname) {
            return 'Nothing has been modified';
        } else if (!newProfile.value.password && !newProfile.value.confirmPassword &&
            newProfile.value.email === oldProfile.value.email &&
            newProfile.value.nickname === oldProfile.value.nickname &&
            newProfile.value.defaultAccountId === oldProfile.value.defaultAccountId &&
            newProfile.value.transactionEditScope === oldProfile.value.transactionEditScope &&
            newProfile.value.language === oldProfile.value.language &&
            newProfile.value.defaultCurrency === oldProfile.value.defaultCurrency &&
            newProfile.value.fiscalYearStart === oldProfile.value.fiscalYearStart &&
            newProfile.value.firstDayOfWeek === oldProfile.value.firstDayOfWeek &&
            newProfile.value.longDateFormat === oldProfile.value.longDateFormat &&
            newProfile.value.shortDateFormat === oldProfile.value.shortDateFormat &&
            newProfile.value.longTimeFormat === oldProfile.value.longTimeFormat &&
            newProfile.value.shortTimeFormat === oldProfile.value.shortTimeFormat &&
            newProfile.value.decimalSeparator === oldProfile.value.decimalSeparator &&
            newProfile.value.digitGroupingSymbol === oldProfile.value.digitGroupingSymbol &&
            newProfile.value.digitGrouping === oldProfile.value.digitGrouping &&
            newProfile.value.currencyDisplayType === oldProfile.value.currencyDisplayType &&
            newProfile.value.expenseAmountColor === oldProfile.value.expenseAmountColor &&
            newProfile.value.incomeAmountColor === oldProfile.value.incomeAmountColor) {
            return 'Nothing has been modified';
        } else if (!newProfile.value.password && newProfile.value.confirmPassword) {
            return 'Password cannot be blank';
        } else if (newProfile.value.password && !newProfile.value.confirmPassword) {
            return 'Password confirmation cannot be blank';
        } else {
            return null;
        }
    });

    const inputInvalidProblemMessage = computed<string | null>(() => {
        if (newProfile.value.password && newProfile.value.confirmPassword && newProfile.value.password !== newProfile.value.confirmPassword) {
            return 'Password and password confirmation do not match';
        } else if (!newProfile.value.email) {
            return 'Email address cannot be blank';
        } else if (!newProfile.value.nickname) {
            return 'Nickname cannot be blank';
        } else if (!newProfile.value.defaultCurrency) {
            return 'Default currency cannot be blank';
        } else {
            return null;
        }
    });

    const langAndRegionInputInvalidProblemMessage = computed<string | null>(() => {
        if (!newProfile.value.defaultCurrency) {
            return 'Default currency cannot be blank';
        } else {
            return null;
        }
    });

    const extendInputInvalidProblemMessage = computed<string | null>(() => {
        return null;
    });

    const inputIsNotChanged = computed<boolean>(() => !!inputIsNotChangedProblemMessage.value);
    const inputIsInvalid = computed<boolean>(() => !!inputInvalidProblemMessage.value);
    const langAndRegionInputIsInvalid = computed<boolean>(() => !!langAndRegionInputInvalidProblemMessage.value);
    const extendInputIsInvalid = computed<boolean>(() => !!extendInputInvalidProblemMessage.value);

    function setCurrentUserProfile(profile: UserBasicInfo): void {
        emailVerified.value = profile.emailVerified;
        oldProfile.value.fillFrom(profile);
        newProfile.value.fillFrom(oldProfile.value);
    }

    function reset(): void {
        newProfile.value.fillFrom(oldProfile.value);
    }

    function doAfterProfileUpdate(user: UserBasicInfo): void {
        if (user) {
            if (user.firstDayOfWeek !== oldProfile.value.firstDayOfWeek) {
                overviewStore.resetTransactionOverview();
            }

            setCurrentUserProfile(user);

            const localeDefaultSettings = setLanguage(user.language);
            settingsStore.updateLocalizedDefaultSettings(localeDefaultSettings);

            setExpenseAndIncomeAmountColor(user.expenseAmountColor, user.incomeAmountColor);
        }
    }

    return {
        // states
        newProfile,
        oldProfile,
        emailVerified,
        loading,
        resending,
        saving,
        // computed states
        allAccounts,
        allVisibleAccounts,
        allVisibleCategorizedAccounts,
        allWeekDays,
        allLongDateFormats,
        allShortDateFormats,
        allLongTimeFormats,
        allShortTimeFormats,
        allDecimalSeparators,
        allDigitGroupingSymbols,
        allDigitGroupingTypes,
        allCurrencyDisplayTypes,
        allExpenseAmountColorTypes,
        allIncomeAmountColorTypes,
        allTransactionEditScopeTypes,
        languageTitle,
        supportDigitGroupingSymbol,
        inputIsNotChangedProblemMessage,
        inputInvalidProblemMessage,
        langAndRegionInputInvalidProblemMessage,
        extendInputInvalidProblemMessage,
        inputIsNotChanged,
        inputIsInvalid,
        langAndRegionInputIsInvalid,
        extendInputIsInvalid,
        // functions
        setCurrentUserProfile,
        reset,
        doAfterProfileUpdate
    };
}
