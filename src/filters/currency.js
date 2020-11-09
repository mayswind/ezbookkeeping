import settings from "../lib/settings.js";
import utils from "../lib/utils.js";

export default function ({ i18n }, value, currencyCode) {
    if (!utils.isNumber(value) && !utils.isString(value)) {
        return value;
    }

    if (utils.isNumber(value)) {
        value = value.toString();
    }

    if (value.length === 0) {
        value = '0.00';
    } else if (value.length === 1) {
        value = '0.0' + value;
    } else if (value.length === 2) {
        value = '0.' + value;
    } else {
        const integer = value.substr(0, value.length - 2);
        const decimals = value.substr(value.length - 2, 2);
        value = `${integer}.${decimals}`;
    }

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
