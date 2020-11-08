import settings from "../lib/settings.js";

export default function ({ i18n }, value, currencyCode) {
    if (!value) {
        return value;
    }

    value = value / 100;
    const currencyDisplayMode = settings.getCurrencyDisplayMode();

    if (currencyDisplayMode === 'code') {
        return `${value} ${currencyCode}`;
    } else if (currencyDisplayMode === 'name') {
        const name = i18n.t(`currency.${currencyCode}`);
        return `${value} ${name}`;
    } else {
        return value;
    }
}
