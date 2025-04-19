import type { TypeAndDisplayName, TypeAndName } from '@/core/base.ts';
import type { UnixTimeRange } from './datetime';

export class FiscalYearStart {
    public static readonly Default = new FiscalYearStart(1, 1);

    public readonly day: number;
    public readonly month: number;
    public readonly value: number;

    private constructor(month: number, day: number) {
        const [validMonth, validDay] = validateMonthDay(month, day);
        this.day = validDay;
        this.month = validMonth;
        this.value = (validMonth << 8) | validDay;
    }

    public static of(month: number, day: number): FiscalYearStart {
        return new FiscalYearStart(month, day);
    }

    public static valueOf(value: number): FiscalYearStart {
        return FiscalYearStart.strictFromNumber(value);
    }

    public static valuesFromNumber(value: number): number[] {
        return FiscalYearStart.strictFromNumber(value).values();
    }

    public values(): number[] {
        return [
            this.month,
            this.day
        ];
    }

    public static parse(valueString: string): FiscalYearStart | undefined {
        return FiscalYearStart.strictFromMonthDashDayString(valueString);
    }

    public static isValidType(value: number): boolean {
        if (value < 0x0101 || value > 0x0C1F) {
            return false;
        }
        
        const month = (value >> 8) & 0xFF;
        const day = value & 0xFF;
        
        try {
            validateMonthDay(month, day);
            return true;
        } catch (error) {
            return false;
        }
    }

    public isValid(): boolean {
        try {
            FiscalYearStart.validateMonthDay(this.month, this.day);
            return true;
        } catch (error) {
            return false;
        }
    }

    public isDefault(): boolean {
        return this.month === 1 && this.day === 1;
    }

    public static validateMonthDay(month: number, day: number): [number, number] {
        return validateMonthDay(month, day);
    }

    public static strictFromMonthDayValues(month: number, day: number): FiscalYearStart {
        return FiscalYearStart.of(month, day);
    }

    /**
     * Create a FiscalYearStart from a uint16 value (two bytes - month high, day low)
     * @param value uint16 value (month in high byte, day in low byte)
     * @returns FiscalYearStart instance
     */
    public static strictFromNumber(value: number): FiscalYearStart {
        if (value < 0 || value > 0xFFFF) {
            throw new Error("Invalid uint16 value");
        }

        const month = (value >> 8) & 0xFF;  // high byte
        const day = value & 0xFF;           // low byte

        try {
            const [validMonth, validDay] = validateMonthDay(month, day);
            return FiscalYearStart.of(validMonth, validDay);
        } catch (error) {
            throw new Error("Invalid uint16 value");
        }
    }

    /**
     * Create a FiscalYearStart from a month/day string
     * @param input MM-dd string (e.g. "04-01" = 1 April)
     * @returns FiscalYearStart instance
     */
    public static strictFromMonthDashDayString(input: string): FiscalYearStart {
        if (!input || !input.includes('-')) {
            throw new Error("Invalid input string");
        }

        const parts = input.split('-');
        if (parts.length !== 2) {
            throw new Error("Invalid input string");
        }

        const month = parseInt(parts[0], 10);
        const day = parseInt(parts[1], 10);

        if (isNaN(month) || isNaN(day)) {
            throw new Error("Invalid input string");
        }

        try {
            const [validMonth, validDay] = validateMonthDay(month, day);
            return FiscalYearStart.of(validMonth, validDay);
        } catch (error) {
            throw new Error("Invalid input string");
        }
    }

    public static fromMonthDashDayString(input: string): FiscalYearStart | null {
        try {
            return FiscalYearStart.strictFromMonthDashDayString(input);
        } catch (error) {
            return null;
        }
    }

    public static fromNumber(value: number): FiscalYearStart | null {
        try {
            return FiscalYearStart.strictFromNumber(value);
        } catch (error) {
            return null;
        }
    }

    public static fromMonthDayValues(month: number, day: number): FiscalYearStart | null {
        try {
            return FiscalYearStart.strictFromMonthDayValues(month, day);
        } catch (error) {
            return null;
        }
    }

    public toMonthDashDayString(): string {
        return `${this.month.toString().padStart(2, '0')}-${this.day.toString().padStart(2, '0')}`;
    }

    public toMonthDayValues(): [string, string] {
        return [
            `${this.month.toString().padStart(2, '0')}`,
            `${this.day.toString().padStart(2, '0')}`
        ]
    }

    public toString(): string {
        return this.toMonthDashDayString();
    }
}

function validateMonthDay(month: number, day: number): [number, number] {
    if (month < 1 || month > 12 || day < 1) {
        throw new Error("Invalid month or day");
    }

    let maxDays = 31;
    switch (month) {
        // January, March, May, July, August, October, December
        case 1: case 3: case 5: case 7: case 8: case 10: case 12: 
            maxDays = 31;
            break;
        // April, June, September, November
        case 4: case 6: case 9: case 11: 
            maxDays = 30;
            break;
        // February
        case 2: 
            maxDays = 28; // Disallow fiscal year start on leap day
            break;
    }

    if (day > maxDays) {
        throw new Error("Invalid day for given month");
    }

    return [month, day];
}

export class FiscalYearUnixTime implements UnixTimeRange {
    public readonly fiscalYear: number;
    public readonly minUnixTime: number;
    public readonly maxUnixTime: number;

    private constructor(fiscalYear: number, minUnixTime: number, maxUnixTime: number) {
        this.fiscalYear = fiscalYear;
        this.minUnixTime = minUnixTime;
        this.maxUnixTime = maxUnixTime;
    }

    public static of(fiscalYear: number, minUnixTime: number, maxUnixTime: number): FiscalYearUnixTime {
        return new FiscalYearUnixTime(fiscalYear, minUnixTime, maxUnixTime);
    }
}

export const LANGUAGE_DEFAULT_FISCAL_YEAR_FORMAT_VALUE: number = 0;

export type FiscalYearFormatTypeName = 'StartYYYY_EndYYYY' | 'StartYYYY_EndYY' | 'StartYY_EndYY' | 'EndYYYY' | 'EndYY';

export class FiscalYearFormat implements TypeAndName {
    private static readonly allInstances: FiscalYearFormat[] = [];
    private static readonly allInstancesByType: Record<number, FiscalYearFormat> = {};
    private static readonly allInstancesByTypeName: Record<string, FiscalYearFormat> = {};

    public static readonly StartYYYY_EndYYYY = new FiscalYearFormat(1, 'StartYYYY_EndYYYY');
    public static readonly StartYYYY_EndYY = new FiscalYearFormat(2, 'StartYYYY_EndYY');
    public static readonly StartYY_EndYY = new FiscalYearFormat(3, 'StartYY_EndYY');
    public static readonly EndYYYY = new FiscalYearFormat(4, 'EndYYYY');
    public static readonly EndYY = new FiscalYearFormat(5, 'EndYY');

    public static readonly Default = FiscalYearFormat.StartYYYY_EndYYYY;

    public readonly type: number;
    public readonly name: FiscalYearFormatTypeName;

    private constructor(type: number, name: FiscalYearFormatTypeName) {
        this.type = type;
        this.name = name;
        
        FiscalYearFormat.allInstances.push(this);
        FiscalYearFormat.allInstancesByType[type] = this;
        FiscalYearFormat.allInstancesByTypeName[name] = this;
    }

    public static values(): FiscalYearFormat[] {
        return FiscalYearFormat.allInstances;
    }

    public static all(): Record<FiscalYearFormatTypeName, FiscalYearFormat> {
        return FiscalYearFormat.allInstancesByTypeName;
    }

    public static valueOf(type: number): FiscalYearFormat | undefined {
        return FiscalYearFormat.allInstancesByType[type];
    }
}

export interface LocalizedFiscalYearFormat extends TypeAndDisplayName {
    readonly type: number;
    readonly format: string;
    readonly displayName: string;
}