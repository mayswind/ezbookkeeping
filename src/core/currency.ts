import type { TypeAndName } from './base.ts';

export enum CurrencyDisplaySymbol {
    None = 0,
    Symbol = 1,
    Code = 2,
    Unit = 3,
    Name = 4
}

export enum CurrencyDisplayLocation {
    BeforeAmount = 0,
    AfterAmount = 1
}

export interface CurrencyPrependAndAppendText {
    prependText?: string;
    appendText?: string;
}

export class CurrencyDisplayType implements TypeAndName {
    private static readonly allInstances: CurrencyDisplayType[] = [];
    private static readonly allInstancesByType: Record<number, CurrencyDisplayType> = {};
    private static readonly allInstancesByTypeName: Record<string, CurrencyDisplayType> = {};

    public static readonly LanguageDefaultType: number = 0;
    public static readonly None = new CurrencyDisplayType(1, 'None', 'None', 2, CurrencyDisplaySymbol.None, undefined, '');
    public static readonly SymbolBeforeAmount = new CurrencyDisplayType(2, 'SymbolBeforeAmount', 'Currency Symbol', 2, CurrencyDisplaySymbol.Symbol, CurrencyDisplayLocation.BeforeAmount, ' ');
    public static readonly SymbolAfterAmount = new CurrencyDisplayType(3, 'SymbolAfterAmount', 'Currency Symbol', 2, CurrencyDisplaySymbol.Symbol, CurrencyDisplayLocation.AfterAmount, ' ');
    public static readonly SymbolBeforeAmountWithoutSpace = new CurrencyDisplayType(4, 'SymbolBeforeAmountWithoutSpace', 'Currency Symbol', 2, CurrencyDisplaySymbol.Symbol, CurrencyDisplayLocation.BeforeAmount, '');
    public static readonly SymbolAfterAmountWithoutSpace = new CurrencyDisplayType(5, 'SymbolAfterAmountWithoutSpace', 'Currency Symbol', 2, CurrencyDisplaySymbol.Symbol, CurrencyDisplayLocation.AfterAmount, '');
    public static readonly CodeBeforeAmount = new CurrencyDisplayType(6, 'CodeBeforeAmount', 'Currency Code', 2, CurrencyDisplaySymbol.Code, CurrencyDisplayLocation.BeforeAmount, ' ');
    public static readonly CodeAfterAmount = new CurrencyDisplayType(7, 'CodeAfterAmount', 'Currency Code', 2, CurrencyDisplaySymbol.Code, CurrencyDisplayLocation.AfterAmount, ' ');
    public static readonly UnitBeforeAmount = new CurrencyDisplayType(8, 'UnitBeforeAmount', 'Currency Unit', 2, CurrencyDisplaySymbol.Unit, CurrencyDisplayLocation.BeforeAmount, ' ');
    public static readonly UnitAfterAmount = new CurrencyDisplayType(9, 'UnitAfterAmount', 'Currency Unit', 2, CurrencyDisplaySymbol.Unit, CurrencyDisplayLocation.AfterAmount, ' ');
    public static readonly NameBeforeAmount = new CurrencyDisplayType(10, 'NameBeforeAmount', 'Currency Name', 2, CurrencyDisplaySymbol.Name, CurrencyDisplayLocation.BeforeAmount, ' ');
    public static readonly NameAfterAmount = new CurrencyDisplayType(11, 'NameAfterAmount', 'Currency Name', 2, CurrencyDisplaySymbol.Name, CurrencyDisplayLocation.AfterAmount, ' ');

    public static readonly Default = CurrencyDisplayType.SymbolBeforeAmount;

    public readonly type: number;
    public readonly typeName: string;
    public readonly name: string;
    public readonly fraction: number;
    public readonly symbol: CurrencyDisplaySymbol;
    public readonly location: CurrencyDisplayLocation | undefined;
    public readonly separator: string;

    private constructor(type: number, typeName: string, name: string, fraction: number, symbol: CurrencyDisplaySymbol, location: CurrencyDisplayLocation | undefined, separator: string) {
        this.type = type;
        this.typeName = typeName;
        this.name = name;
        this.fraction = fraction;
        this.symbol = symbol;
        this.location = location;
        this.separator = separator;

        CurrencyDisplayType.allInstances.push(this);
        CurrencyDisplayType.allInstancesByType[type] = this;
        CurrencyDisplayType.allInstancesByTypeName[typeName] = this;
    }

    public static values(): CurrencyDisplayType[] {
        return CurrencyDisplayType.allInstances;
    }

    public static valueOf(type: number): CurrencyDisplayType {
        return CurrencyDisplayType.allInstancesByType[type];
    }

    public static parse(typeName: string): CurrencyDisplayType {
        return CurrencyDisplayType.allInstancesByTypeName[typeName];
    }
}

export class CurrencySortingType implements TypeAndName {
    private static readonly allInstances: CurrencySortingType[] = [];

    public static readonly Name = new CurrencySortingType(0, 'Currency Name');
    public static readonly CurrencyCode = new CurrencySortingType(1, 'Currency Code');
    public static readonly ExchangeRate = new CurrencySortingType(2, 'Exchange Rate');

    public static readonly Default = CurrencySortingType.Name;

    public readonly type: number;
    public readonly name: string;

    private constructor(type: number, name: string) {
        this.type = type;
        this.name = name;

        CurrencySortingType.allInstances.push(this);
    }

    public static values(): CurrencySortingType[] {
        return CurrencySortingType.allInstances;
    }
}
