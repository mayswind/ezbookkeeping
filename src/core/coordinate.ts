import type { TypeAndName } from '@/core/base.ts';

export interface Coordinate {
    latitude: number;
    longitude: number;
}

export enum CoordinateDisplayOrder {
    LatitudeLongitude = 0,
    LongitudeLatitude = 1
}

export enum CoordinateDisplayFormat {
    DecimalDegrees = 0,
    DecimalMinutes = 1,
    DegreesMinutesSeconds = 2
}

export enum CoordinateDirectionFormat {
    Signed = 0,
    Directional = 1
}

export class CoordinateDisplayType implements TypeAndName {
    private static readonly allInstances: CoordinateDisplayType[] = [];
    private static readonly allInstancesByType: Record<number, CoordinateDisplayType> = {};

    public static readonly SystemDefaultType: number = 0;
    public static readonly LatitudeLongitudeDecimalDegrees = new CoordinateDisplayType(1, 'Latitude Longitude D.D°', CoordinateDisplayOrder.LatitudeLongitude, CoordinateDisplayFormat.DecimalDegrees, CoordinateDirectionFormat.Signed);
    public static readonly LongitudeLatitudeDecimalDegrees = new CoordinateDisplayType(2, 'Longitude Latitude D.D°', CoordinateDisplayOrder.LongitudeLatitude, CoordinateDisplayFormat.DecimalDegrees, CoordinateDirectionFormat.Signed);
    public static readonly LatitudeLongitudeDecimalMinutes = new CoordinateDisplayType(3, 'Latitude Longitude D°M.M\'', CoordinateDisplayOrder.LatitudeLongitude, CoordinateDisplayFormat.DecimalMinutes, CoordinateDirectionFormat.Directional);
    public static readonly LongitudeLatitudeDecimalMinutes = new CoordinateDisplayType(4, 'Longitude Latitude D°M.M\'', CoordinateDisplayOrder.LongitudeLatitude, CoordinateDisplayFormat.DecimalMinutes, CoordinateDirectionFormat.Directional);
    public static readonly LatitudeLongitudeDegreesMinutesSeconds = new CoordinateDisplayType(5, 'Latitude Longitude D°M\'S"', CoordinateDisplayOrder.LatitudeLongitude, CoordinateDisplayFormat.DegreesMinutesSeconds, CoordinateDirectionFormat.Directional);
    public static readonly LongitudeLatitudeDegreesMinutesSeconds = new CoordinateDisplayType(6, 'Longitude Latitude D°M\'S"', CoordinateDisplayOrder.LongitudeLatitude, CoordinateDisplayFormat.DegreesMinutesSeconds, CoordinateDirectionFormat.Directional);

    public static readonly Default = CoordinateDisplayType.LatitudeLongitudeDecimalDegrees;

    public readonly type: number;
    public readonly name: string;
    public readonly displayOrder: CoordinateDisplayOrder;
    public readonly displayFormat: CoordinateDisplayFormat;
    public readonly directionFormat: CoordinateDirectionFormat;

    private constructor(type: number, name: string, displayOrder: CoordinateDisplayOrder, displayFormat: CoordinateDisplayFormat, directionFormat: CoordinateDirectionFormat) {
        this.type = type;
        this.name = name;
        this.displayOrder = displayOrder;
        this.displayFormat = displayFormat;
        this.directionFormat = directionFormat;

        CoordinateDisplayType.allInstances.push(this);
        CoordinateDisplayType.allInstancesByType[type] = this;
    }

    public static values(): CoordinateDisplayType[] {
        return CoordinateDisplayType.allInstances;
    }

    public static valueOf(type: number): CoordinateDisplayType | undefined {
        return CoordinateDisplayType.allInstancesByType[type];
    }
}

export function getNormalizedCoordinate(value: Coordinate): Coordinate {
    if (!value) {
        return value;
    }

    const normalizedLatitude = Math.max(-90, Math.min(90, value.latitude));
    const normalizedLongitude = ((value.longitude + 180) % 360 + 360) % 360 - 180;

    return {
        latitude: normalizedLatitude,
        longitude: normalizedLongitude
    };
}
