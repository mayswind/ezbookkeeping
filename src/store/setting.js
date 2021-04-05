import {
    UPDATE_DEFAULT_SETTING
} from './mutations.js';

export function updateLocalizedDefaultSettings(context, { defaultCurrency, defaultFirstDayOfWeek }) {
    context.commit(UPDATE_DEFAULT_SETTING, {
        defaultCurrency,
        defaultFirstDayOfWeek,
    });
}
