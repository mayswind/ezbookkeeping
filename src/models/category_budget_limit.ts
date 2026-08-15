// [PLUGIN:budget] Category budget limit feature - TypeScript model.
// Mirrors the backend DTOs in pkg/models/category_budget_limit.go.

export class CategoryBudgetLimit implements CategoryBudgetLimitInfoResponse {
    public id: string;
    public categoryId: string;
    public startDate: number;
    public endDate: number;
    public amount: number;
    public currency: string;

    public constructor(id: string, categoryId: string, startDate: number, endDate: number, amount: number, currency: string) {
        this.id = id;
        this.categoryId = categoryId;
        this.startDate = startDate;
        this.endDate = endDate;
        this.amount = amount;
        this.currency = currency;
    }

    public toCreateRequest(clientSessionId: string): CategoryBudgetLimitCreateRequest {
        return {
            categoryId: this.categoryId,
            // startDate is sent as a quoted string: the backend CreateRequest uses the
            // `,string` json tag (consistent with all large int64 fields in ezBookkeeping).
            startDate: String(this.startDate),
            amount: this.amount,
            currency: this.currency,
            clientSessionId
        };
    }

    public toModifyRequest(): CategoryBudgetLimitModifyRequest {
        return {
            id: this.id,
            amount: this.amount,
            currency: this.currency
        };
    }

    public static of(resp: CategoryBudgetLimitInfoResponse): CategoryBudgetLimit {
        return new CategoryBudgetLimit(resp.id, resp.categoryId, resp.startDate, resp.endDate, resp.amount, resp.currency);
    }

    public static ofMulti(responses: CategoryBudgetLimitInfoResponse[]): CategoryBudgetLimit[] {
        return responses.map(CategoryBudgetLimit.of);
    }

    public static createNew(categoryId: string, startDate: number, currency: string): CategoryBudgetLimit {
        return new CategoryBudgetLimit('', categoryId, startDate, 0, 0, currency);
    }
}

export interface CategoryBudgetLimitCreateRequest {
    readonly categoryId: string;
    readonly startDate: string;
    readonly amount: number;
    readonly currency: string;
    readonly clientSessionId: string;
}

export interface CategoryBudgetLimitModifyRequest {
    readonly id: string;
    readonly amount: number;
    readonly currency: string;
}

export interface CategoryBudgetLimitListByMonthRequest {
    readonly startDate: string;
}

export interface CategoryBudgetLimitInfoResponse {
    readonly id: string;
    readonly categoryId: string;
    readonly startDate: number;
    readonly endDate: number;
    readonly amount: number;
    readonly currency: string;
}

export interface CategoryBudgetOverviewItem extends CategoryBudgetLimitInfoResponse {
    readonly categoryName: string;
    readonly categoryParentId: string;
    readonly categoryType: number;
    readonly actualExpenseAmount: number;
    readonly availableAmount: number;
}

export interface CategoryBudgetOverviewResponse {
    readonly startDate: number;
    readonly endDate: number;
    readonly totalLimit: number;
    readonly totalActualExpenseAmount: number;
    readonly totalAvailableAmount: number;
    readonly items: CategoryBudgetOverviewItem[];
}
