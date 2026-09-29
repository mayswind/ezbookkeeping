import {
    type Coordinate,
    type CoordinateFormatOptions,
    CoordinateDisplayOrder,
    CoordinateDisplayFormat,
    CoordinateDirectionFormat,
    getNormalizedCoordinate
} from '@/core/coordinate.ts';
import {
    type NumberFormatOptions,
    DecimalSeparator
} from '@/core/numeral.ts';

import { formatNumber } from '@/lib/numeral.ts';

export function formatCoordinate(value: Coordinate, options: CoordinateFormatOptions): string {
    if (!value) {
        return '';
    }

    value = getNormalizedCoordinate(value);

    const displayType = options.coordinateDisplayType;
    const numberFormatOptions = options.numberFormatOptions;

    const formattedLatitude = formatCoordinateValue(value.latitude, numberFormatOptions, 'N', 'S', displayType.displayFormat, displayType.directionFormat);
    const formattedLongitude = formatCoordinateValue(value.longitude, numberFormatOptions, 'E', 'W', displayType.displayFormat, displayType.directionFormat);
    const coordinateSeparator = numberFormatOptions.decimalSeparator === DecimalSeparator.Comma.symbol ? ';' : ',';

    if (displayType.displayOrder === CoordinateDisplayOrder.LatitudeLongitude) {
        return `${formattedLatitude}${coordinateSeparator} ${formattedLongitude}`;
    } else if (displayType.displayOrder === CoordinateDisplayOrder.LongitudeLatitude) {
        return `${formattedLongitude}${coordinateSeparator} ${formattedLatitude}`;
    } else {
        return '';
    }
}

function formatCoordinateValue(value: number, numberFormatOptions: NumberFormatOptions, positiveDirectionName: string, negativeDirectionName: string, displayFormat: CoordinateDisplayFormat, directionFormat: CoordinateDirectionFormat): string {
    let prefix = '';
    let suffix = '';

    if (directionFormat === CoordinateDirectionFormat.Signed) {
        prefix = value >= 0 ? '' : '-';
    } else if (directionFormat === CoordinateDirectionFormat.Directional) {
        suffix = value >= 0 ? positiveDirectionName : negativeDirectionName;
    }

    value = Math.abs(value);

    if (displayFormat === CoordinateDisplayFormat.DecimalDegrees) {
        return `${prefix}${formatNumber(value, numberFormatOptions, 6)}${suffix}`;
    } else if (displayFormat === CoordinateDisplayFormat.DecimalMinutes) {
        const degrees = Math.trunc(value);
        const minutes = (value - degrees) * 60;
        return `${prefix}${degrees}°${formatNumber(minutes, numberFormatOptions, 5)}'${suffix}`;
    } else if (displayFormat === CoordinateDisplayFormat.DegreesMinutesSeconds) {
        const degrees = Math.trunc(value);
        const minutes = Math.trunc((value - degrees) * 60);
        const seconds = (value - degrees - minutes / 60) * 3600;
        return `${prefix}${degrees}°${minutes}'${formatNumber(seconds, numberFormatOptions, 4)}"${suffix}`;
    } else {
        return '';
    }
}
