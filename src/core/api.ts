export interface ApiResponse<T> {
    readonly success: boolean;
    readonly result: T;
}

export interface ErrorResponse {
    readonly success: boolean;
    readonly errorCode: number;
    readonly errorMessage: string;
    readonly path: string;
}
