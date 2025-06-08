import type { UnixTimeRange } from './datetime.ts';

export class FiscalYearStart {
    public static readonly JanuaryFirstDay = new FiscalYearStart(1, 1);
    public static readonly Default = FiscalYearStart.JanuaryFirstDay;

    private static readonly MONTH_MAX_DAYS: number[] = [
        31, // January
        28, // February (Disallow fiscal year start on leap day)
        31, // March
        30, // April
        31, // May
        30, // June
        31, // July
        31, // August
        30, // September
        31, // October
        30, // November
        31 // December
    ];

    public readonly month: number; // 1-based (1 = January, 12 = December)
    public readonly day: number;
    public readonly value: number;

    private constructor(month: number, day: number) {
        this.month = month;
        this.day = day;
        this.value = (month << 8) | day;
    }

    public static of(month: number, day: number): FiscalYearStart | undefined {
        if (!FiscalYearStart.isValidFiscalYearMonthDay(month, day)) {
            return undefined;
        }

        return new FiscalYearStart(month, day);
    }

    /**
     * Create a FiscalYearStart from a uint16 value (two bytes - month high, day low)
     * @param value uint16 value (month in high byte, day in low byte)
     * @returns FiscalYearStart instance or undefined if the value is out of range
     */
    public static valueOf(value: number): FiscalYearStart | undefined {
        if (value < 0x0101 || value > 0x0C1F) {
            return undefined;
        }

        const month = (value >> 8) & 0xFF;  // high byte
        const day = value & 0xFF;           // low byte

        return FiscalYearStart.of(month, day);
    }

    /**
     * Create a FiscalYearStart from a month/day string
     * @param monthDay MM-dd string (e.g. "04-01" = 1 April)
     * @returns FiscalYearStart instance or undefined if the monthDay is invalid
     */
    public static parse(monthDay: string): FiscalYearStart | undefined {
        if (!monthDay || !monthDay.includes('-')) {
            return undefined;
        }

        const parts = monthDay.split('-');

        if (parts.length !== 2) {
            return undefined;
        }

        const month = parseInt(parts[0], 10);
        const day = parseInt(parts[1], 10);

        return FiscalYearStart.of(month, day);
    }

    public toMonthDashDayString(): string {
        return `${this.month.toString().padStart(2, '0')}-${this.day.toString().padStart(2, '0')}`;
    }

    private static isValidFiscalYearMonthDay(month: number, day: number): boolean {
        return 1 <= month && month <= 12 && 1 <= day && day <= FiscalYearStart.MONTH_MAX_DAYS[month - 1];
    }
}

export class FiscalYearUnixTime implements UnixTimeRange {
    public readonly year: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(fiscalYear: number, minUnixTime: number, maxUnixTime: number) {
        this.year = fiscalYear;
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(fiscalYear: number, minUnixTime: number, maxUnixTime: number): FiscalYearUnixTime {
        return new FiscalYearUnixTime(fiscalYear, minUnixTime, maxUnixTime);
    }
}

export const LANGUAGE_DEFAULT_FISCAL_YEAR_FORMAT_VALUE: number = 0;

export class FiscalYearFormat {
    private static readonly allInstances: FiscalYearFormat[] = [];
    private static readonly allInstancesByType: Record<number, FiscalYearFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, FiscalYearFormat> = {};

    public static readonly StartYYYY_EndYYYY = new FiscalYearFormat(1, 'StartYYYY_EndYYYY');
    public static readonly StartYYYY_EndYY = new FiscalYearFormat(2, 'StartYYYY_EndYY');
    public static readonly StartYY_EndYY = new FiscalYearFormat(3, 'StartYY_EndYY');
    public static readonly EndYYYY = new FiscalYearFormat(4, 'EndYYYY');
    public static readonly EndYY = new FiscalYearFormat(5, 'EndYY');

    public static readonly Default = FiscalYearFormat.EndYYYY;

    public readonly type: number;
    public readonly typeName: string;

    private constructor(type: number, typeName: string) {
        this.type = type;
        this.typeName = typeName;

        FiscalYearFormat.allInstances.push(this);
        FiscalYearFormat.allInstancesByType[type] = this;
        FiscalYearFormat.allInstancesByTypeName[typeName] = this;
    }

    public static values(): FiscalYearFormat[] {
        return FiscalYearFormat.allInstances;
    }

    public static valueOf(type: number): FiscalYearFormat | undefined {
        return FiscalYearFormat.allInstancesByType[type];
    }

    public static parse(typeName: string): FiscalYearFormat | undefined {
        return FiscalYearFormat.allInstancesByTypeName[typeName];
    }
}
