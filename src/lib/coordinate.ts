import {
    type Coordinate,
    CoordinateDisplayOrder,
    CoordinateDisplayFormat,
    CoordinateDirectionFormat,
    CoordinateDisplayType,
    getNormalizedCoordinate
} from '@/core/coordinate.ts';

export function formatCoordinate(value: Coordinate, coordinateDisplayType: number): string {
    if (!value) {
        return '';
    }

    value = getNormalizedCoordinate(value);

    const displayType = CoordinateDisplayType.valueOf(coordinateDisplayType) || CoordinateDisplayType.Default;
    const formattedLatitude = formatCoordinateValue(value.latitude, 'N', 'S', displayType.displayFormat, displayType.directionFormat);
    const formattedLongitude = formatCoordinateValue(value.longitude, 'E', 'W', displayType.displayFormat, displayType.directionFormat);

    if (displayType.displayOrder === CoordinateDisplayOrder.LatitudeLongitude) {
        return `${formattedLatitude}, ${formattedLongitude}`;
    } else if (displayType.displayOrder === CoordinateDisplayOrder.LongitudeLatitude) {
        return `${formattedLongitude}, ${formattedLatitude}`;
    } else {
        return '';
    }
}

function formatCoordinateValue(value: number, positiveDirectionName: string, negativeDirectionName: string, displayFormat: CoordinateDisplayFormat, directionFormat: CoordinateDirectionFormat): string {
    let prefix = '';
    let suffix = '';

    if (directionFormat === CoordinateDirectionFormat.Signed) {
        prefix = value >= 0 ? '' : '-';
    } else if (directionFormat === CoordinateDirectionFormat.Directional) {
        suffix = value >= 0 ? positiveDirectionName : negativeDirectionName;
    }

    value = Math.abs(value);

    if (displayFormat === CoordinateDisplayFormat.DecimalDegrees) {
        return `${prefix}${value.toFixed(6)}${suffix}`;
    } else if (displayFormat === CoordinateDisplayFormat.DecimalMinutes) {
        const degrees = Math.trunc(value);
        const minutes = (value - degrees) * 60;
        return `${prefix}${degrees}°${minutes.toFixed(5)}'${suffix}`;
    } else if (displayFormat === CoordinateDisplayFormat.DegreesMinutesSeconds) {
        const degrees = Math.trunc(value);
        const minutes = Math.trunc((value - degrees) * 60);
        const seconds = (value - degrees - minutes / 60) * 3600;
        return `${prefix}${degrees}°${minutes}'${seconds.toFixed(4)}"${suffix}`;
    } else {
        return '';
    }
}
