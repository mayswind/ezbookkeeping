export interface ApiResponse<T> {
    readonly success: boolean;
    readonly result: T;
}

export interface ErrorResponse {
    readonly success: boolean;
    readonly errorCode: number;
    readonly errorMessage: string;
    readonly path: string;
    readonly context?: {
        readonly email?: string;
        readonly hasValidEmailVerifyToken?: boolean;
        readonly recordId?: string;
        readonly originalData?: unknown;
    };
}

export function buildErrorResponse(errorCode: number, errorMessage: string): ErrorResponse {
    const errorResponse: ErrorResponse = {
        success: false,
        errorCode: errorCode,
        errorMessage: errorMessage,
        path: ''
    };

    return errorResponse;
}
