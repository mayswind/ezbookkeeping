export default function (value, precision, lowPrecisionValue) {
    const ratio = Math.pow(10, precision);
    const normalizedValue = Math.floor(value * ratio);

    if (value > 0 && normalizedValue < 1 && lowPrecisionValue) {
        return lowPrecisionValue + '%';
    }

    const result = normalizedValue / ratio;
    return result + '%';
}
