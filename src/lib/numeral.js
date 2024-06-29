export function appendThousandsSeparator(value, enable) {
    if (!enable || value.length <= 3) {
        return value;
    }

    const negative = value.charAt(0) === '-';

    if (negative) {
        value = value.substring(1);
    }

    const dotPos = value.indexOf('.');
    const integer = dotPos < 0 ? value : value.substring(0, dotPos);
    const decimals = dotPos < 0 ? '' : value.substring(dotPos + 1, value.length);

    const finalChars = [];

    for (let i = 0; i < integer.length; i++) {
        if (i % 3 === 0 && i > 0) {
            finalChars.push(',');
        }

        finalChars.push(integer.charAt(integer.length - 1 - i));
    }

    finalChars.reverse();

    let newInteger = finalChars.join('');

    if (negative) {
        newInteger = `-${newInteger}`;
    }

    if (dotPos < 0) {
        return newInteger;
    } else {
        return `${newInteger}.${decimals}`;
    }
}

export function formatPercent(value, precision, lowPrecisionValue) {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        return lowPrecisionValue + '%';
    }

    const result = normalizedValue / ratio;
    return result + '%';
}
