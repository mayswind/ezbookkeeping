export interface LatestExchangeRate {
    readonly currency: string;
    readonly rate: string;
}

export interface LatestExchangeRateResponse {
    readonly dataSource: string;
    readonly referenceUrl: string;
    readonly updateTime: number;
    readonly baseCurrency: string;
    readonly exchangeRates: LatestExchangeRate[];
}
