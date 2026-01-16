export class TransactionTagGroup implements TransactionTagGroupInfoResponse {
    public id: string;
    public name: string;
    public displayOrder: number;

    private constructor(id: string, name: string, displayOrder: number) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
    }

    public toCreateRequest(): TransactionTagGroupCreateRequest {
        return {
            name: this.name
        };
    }

    public toModifyRequest(): TransactionTagGroupModifyRequest {
        return {
            id: this.id,
            name: this.name
        };
    }

    public clone(): TransactionTagGroup {
        return new TransactionTagGroup(this.id, this.name, this.displayOrder);
    }

    public static of(tagGroupResponse: TransactionTagGroupInfoResponse): TransactionTagGroup {
        return new TransactionTagGroup(tagGroupResponse.id, tagGroupResponse.name, tagGroupResponse.displayOrder);
    }

    public static ofMulti(tagGroupResponses: TransactionTagGroupInfoResponse[]): TransactionTagGroup[] {
        const tagGroups: TransactionTagGroup[] = [];

        for (const tagGroupResponse of tagGroupResponses) {
            tagGroups.push(TransactionTagGroup.of(tagGroupResponse));
        }

        return tagGroups;
    }

    public static createNewTagGroup(name?: string): TransactionTagGroup {
        return new TransactionTagGroup('', name || '', 0);
    }
}

export interface TransactionTagGroupCreateRequest {
    readonly name: string;
}

export interface TransactionTagGroupModifyRequest {
    readonly id: string;
    readonly name: string;
}

export interface TransactionTagGroupMoveRequest {
    readonly newDisplayOrders: TransactionTagGroupNewDisplayOrderRequest[];
}

export interface TransactionTagGroupNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface TransactionTagGroupDeleteRequest {
    readonly id: string;
}

export interface TransactionTagGroupInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly displayOrder: number;
}
