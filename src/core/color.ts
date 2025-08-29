import type { TypeAndName } from './base.ts';

export type ColorValue = string;

export type ColorStyleValue = `#${ColorValue}` | 'var(--default-icon-color)';

export interface ColorInfo extends Record<string, unknown> {
    readonly color: ColorValue;
}

export interface AmountColor {
    readonly expenseAmountColor: ColorStyleValue;
    readonly incomeAmountColor: ColorStyleValue;
}

export class PresetAmountColor implements TypeAndName {
    private static readonly allInstances: PresetAmountColor[] = [];
    private static readonly allInstancesByType: Record<number, PresetAmountColor> = {};

    public static readonly SystemDefaultType: number = 0;
    public static readonly Green = new PresetAmountColor(1, 'Green', '#009688', '#009688', 'expense-amount-color-green', 'income-amount-color-green');
    public static readonly Red = new PresetAmountColor(2, 'Red', '#d43f3f', '#d43f3f', 'expense-amount-color-red', 'income-amount-color-red');
    public static readonly Yellow = new PresetAmountColor(3, 'Yellow', '#e2b60a', '#e2b60a', 'expense-amount-color-yellow', 'income-amount-color-yellow');
    public static readonly BlackOrWhite = new PresetAmountColor(4, 'Black or White', '#413935', '#fcf0e3', 'expense-amount-color-blackorwhite', 'income-amount-color-blackorwhite');

    public static readonly DefaultExpenseColor = PresetAmountColor.Green;
    public static readonly DefaultIncomeColor = PresetAmountColor.Red;

    public readonly type: number;
    public readonly name: string;
    public readonly lightThemeColor: ColorStyleValue;
    public readonly darkThemeColor: ColorStyleValue;
    public readonly expenseClassName: string;
    public readonly incomeClassName: string;

    private constructor(type: number, name: string, lightThemeColor: ColorStyleValue, darkThemeColor: ColorStyleValue, expenseClassName: string, incomeClassName: string) {
        this.type = type;
        this.name = name;
        this.lightThemeColor = lightThemeColor;
        this.darkThemeColor = darkThemeColor;
        this.expenseClassName = expenseClassName;
        this.incomeClassName = incomeClassName;

        PresetAmountColor.allInstances.push(this);
        PresetAmountColor.allInstancesByType[type] = this;
    }

    public static values(): PresetAmountColor[] {
        return PresetAmountColor.allInstances;
    }

    public static valueOf(type: number): PresetAmountColor | undefined {
        return PresetAmountColor.allInstancesByType[type];
    }
}
