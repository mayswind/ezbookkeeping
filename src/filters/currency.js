import settings from "../lib/settings.js";
import utils from "../lib/utils.js";

export default function ({ i18n }, value, currencyCode) {
    if (!utils.isNumber(value) && !utils.isString(value)) {
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
