export class TransactionTag implements TransactionTagInfoResponse {
    public id: string;
    public name: string;
    public displayOrder: number;
    public hidden: boolean;

    private constructor(id: string, name: string, displayOrder: number, hidden: boolean) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
        this.hidden = hidden;
    }

    public toCreateRequest(): TransactionTagCreateRequest {
        return {
            name: this.name
        };
    }

    public toModifyRequest(): TransactionTagModifyRequest {
        return {
            id: this.id,
            name: this.name
        };
    }

    public static of(tagResponse: TransactionTagInfoResponse): TransactionTag {
        return new TransactionTag(tagResponse.id, tagResponse.name, tagResponse.displayOrder, tagResponse.hidden);
    }

    public static ofMany(tagResponses: TransactionTagInfoResponse[]): TransactionTag[] {
        const tags: TransactionTag[] = [];

        for (const tagResponse of tagResponses) {
            tags.push(TransactionTag.of(tagResponse));
        }

        return tags;
    }

    public static createNewTag(): TransactionTag {
        return new TransactionTag('', '', 0, false);
    }
}

export interface TransactionTagCreateRequest {
    readonly name: string;
}

export interface TransactionTagModifyRequest {
    readonly id: string;
    readonly name: string;
}

export interface TransactionTagHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface TransactionTagMoveRequest {
    readonly newDisplayOrders: TransactionTagNewDisplayOrderRequest[];
}

export interface TransactionTagNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface TransactionTagDeleteRequest {
    readonly id: string;
}

export interface TransactionTagInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly displayOrder: number;
    readonly hidden: boolean;
}
