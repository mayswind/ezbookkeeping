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