export interface UserCustomExchangeRateUpdateRequest {
    readonly currency: string;
    readonly rate: string;
}

export interface UserCustomExchangeRateDeleteRequest {
    readonly currency: string;
}

export interface UserCustomExchangeRateUpdateResponse extends LatestExchangeRate {
    readonly updateTime: number;
}

export interface LatestExchangeRate {
    readonly currency: string;
    readonly rate: string;
}

export interface LatestExchangeRateResponse {
    readonly dataSource: string;
    readonly referenceUrl: string;
    updateTime: number;
    readonly baseCurrency: string;
    readonly exchangeRates: LatestExchangeRate[];
}

export interface LocalizedLatestExchangeRate {
    readonly currencyCode: string;
    readonly currencyDisplayName: string;
    readonly rate: string;
}
