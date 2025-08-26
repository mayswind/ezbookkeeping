import type { TypeAndName } from '@/core/base.ts';

export enum CalendarType {
    Gregorian = 0,
    Buddhist = 1,
    Chinese = 2
}

export class CalendarDisplayType implements TypeAndName {
    private static readonly allInstances: CalendarDisplayType[] = [];
    private static readonly allInstancesByType: Record<number, CalendarDisplayType> = {};
    private static readonly allInstancesByTypeName: Record<string, CalendarDisplayType> = {};

    public static readonly LanguageDefaultType: number = 0;
    public static readonly Gregorian = new CalendarDisplayType(1, 'Gregorian', 'Gregorian', CalendarType.Gregorian);

    public static readonly Default = CalendarDisplayType.Gregorian;

    public readonly type: number;
    public readonly typeName: string;
    public readonly name: string;
    public readonly primaryCalendarType: CalendarType;
    public readonly secondaryCalendarType?: CalendarType;

    private constructor(type: number, typeName: string, name: string, primaryCalendarType: CalendarType, secondaryCalendarType?: CalendarType) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;
        this.primaryCalendarType = primaryCalendarType;
        this.secondaryCalendarType = secondaryCalendarType;

        CalendarDisplayType.allInstances.push(this);
        CalendarDisplayType.allInstancesByType[type] = this;
        CalendarDisplayType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): CalendarDisplayType[] {
        return CalendarDisplayType.allInstances;
    }

    public static valueOf(type: number): CalendarDisplayType | undefined {
        return CalendarDisplayType.allInstancesByType[type];
    }

    public static parse(typeName: string): CalendarDisplayType | undefined {
        return CalendarDisplayType.allInstancesByTypeName[typeName];
    }
}

export class DateDisplayType implements TypeAndName {
    private static readonly allInstances: DateDisplayType[] = [];
    private static readonly allInstancesByType: Record<number, DateDisplayType> = {};
    private static readonly allInstancesByTypeName: Record<string, DateDisplayType> = {};

    public static readonly LanguageDefaultType: number = 0;
    public static readonly Gregorian = new DateDisplayType(1, 'Gregorian', 'Gregorian', CalendarType.Gregorian);
    public static readonly Chinese = new DateDisplayType(2, 'Buddhist', 'Buddhist', CalendarType.Buddhist);

    public static readonly Default = DateDisplayType.Gregorian;

    public readonly type: number;
    public readonly typeName: string;
    public readonly name: string;
    public readonly calendarType: CalendarType;

    private constructor(type: number, typeName: string, name: string, calendarType: CalendarType) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;
        this.calendarType = calendarType;

        DateDisplayType.allInstances.push(this);
        DateDisplayType.allInstancesByType[type] = this;
        DateDisplayType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): DateDisplayType[] {
        return DateDisplayType.allInstances;
    }

    public static valueOf(type: number): DateDisplayType | undefined {
        return DateDisplayType.allInstancesByType[type];
    }

    public static parse(typeName: string): DateDisplayType | undefined {
        return DateDisplayType.allInstancesByTypeName[typeName];
    }
}
