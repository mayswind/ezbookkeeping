import { type NameValue } from '@/core/base.ts';
import { DateRange } from '@/core/datetime.ts';

export enum TransactionExploreConditionRelation {
    First = 'first',
    And = 'and',
    Or = 'or'
}

export const TransactionExploreConditionRelationPriority: Record<TransactionExploreConditionRelation, number> = {
    [TransactionExploreConditionRelation.First]: 0,
    [TransactionExploreConditionRelation.Or]: 1,
    [TransactionExploreConditionRelation.And]: 2
};


export enum TransactionExploreConditionFieldType {
    TransactionType = 'transactionType',
    TransactionCategory = 'transactionCategory',
    SourceAccount = 'sourceAccount',
    DestinationAccount = 'destinationAccount',
    SourceAmount = 'sourceAmount',
    DestinationAmount = 'destinationAmount',
    TransactionTag = 'transactionTag',
    Description = 'description'
}

export class TransactionExploreConditionField implements NameValue {
    private static readonly allInstances: TransactionExploreConditionField[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreConditionField> = {};

    public static readonly TransactionType = new TransactionExploreConditionField('Transaction Type', TransactionExploreConditionFieldType.TransactionType);
    public static readonly TransactionCategory = new TransactionExploreConditionField('Category', TransactionExploreConditionFieldType.TransactionCategory);
    public static readonly SourceAccount = new TransactionExploreConditionField('Source Account', TransactionExploreConditionFieldType.SourceAccount);
    public static readonly DestinationAccount = new TransactionExploreConditionField('Destination Account', TransactionExploreConditionFieldType.DestinationAccount);
    public static readonly SourceAmount = new TransactionExploreConditionField('Amount', TransactionExploreConditionFieldType.SourceAmount);
    public static readonly DestinationAmount = new TransactionExploreConditionField('Transfer In Amount', TransactionExploreConditionFieldType.DestinationAmount);
    public static readonly TransactionTag = new TransactionExploreConditionField('Tags', TransactionExploreConditionFieldType.TransactionTag);
    public static readonly Description = new TransactionExploreConditionField('Description', TransactionExploreConditionFieldType.Description);

    public readonly name: string;
    public readonly value: TransactionExploreConditionFieldType;

    private constructor(name: string, value: TransactionExploreConditionFieldType) {
        this.name = name;
        this.value = value;

        TransactionExploreConditionField.allInstances.push(this);
        TransactionExploreConditionField.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreConditionField[] {
        return TransactionExploreConditionField.allInstances;
    }

    public static valueOf(type: string): TransactionExploreConditionField | undefined {
        return TransactionExploreConditionField.allInstancesByValue[type];
    }
}

export enum TransactionExploreConditionOperatorType {
    In = 'in',
    GreaterThan = 'greaterThan',
    LessThan = 'lessThan',
    Equals = 'equals',
    NotEquals = 'notEquals',
    Between = 'between',
    NotBetween = 'notBetween',
    HasAny = 'hasAny',
    HasAll = 'hasAll',
    NotHasAny = 'notHasAny',
    NotHasAll = 'notHasAll',
    IsEmpty = 'isEmpty',
    IsNotEmpty = 'isNotEmpty',
    Contains = 'contains',
    NotContains = 'notContains',
    StartsWith = 'startsWith',
    NotStartsWith = 'notStartsWith',
    EndsWith = 'endsWith',
    NotEndsWith = 'notEndsWith'
}

export class TransactionExploreConditionOperator implements NameValue {
    private static readonly allInstances: TransactionExploreConditionOperator[] = [];
    private static readonly allInstancesByValue: Record<string, TransactionExploreConditionOperator> = {};

    public static readonly In = new TransactionExploreConditionOperator('In', TransactionExploreConditionOperatorType.In);
    public static readonly GreaterThan = new TransactionExploreConditionOperator('Greater than', TransactionExploreConditionOperatorType.GreaterThan);
    public static readonly LessThan = new TransactionExploreConditionOperator('Less than', TransactionExploreConditionOperatorType.LessThan);
    public static readonly Equals = new TransactionExploreConditionOperator('Equal to', TransactionExploreConditionOperatorType.Equals);
    public static readonly NotEquals = new TransactionExploreConditionOperator('Not equal to', TransactionExploreConditionOperatorType.NotEquals);
    public static readonly Between = new TransactionExploreConditionOperator('Between', TransactionExploreConditionOperatorType.Between);
    public static readonly NotBetween = new TransactionExploreConditionOperator('Not between', TransactionExploreConditionOperatorType.NotBetween);
    public static readonly HasAny = new TransactionExploreConditionOperator('Has any', TransactionExploreConditionOperatorType.HasAny);
    public static readonly HasAll = new TransactionExploreConditionOperator('Has all', TransactionExploreConditionOperatorType.HasAll);
    public static readonly NotHasAny = new TransactionExploreConditionOperator('Not has any', TransactionExploreConditionOperatorType.NotHasAny);
    public static readonly NotHasAll = new TransactionExploreConditionOperator('Not has all', TransactionExploreConditionOperatorType.NotHasAll);
    public static readonly IsEmpty = new TransactionExploreConditionOperator('Is empty', TransactionExploreConditionOperatorType.IsEmpty);
    public static readonly IsNotEmpty = new TransactionExploreConditionOperator('Is not empty', TransactionExploreConditionOperatorType.IsNotEmpty);
    public static readonly Contains = new TransactionExploreConditionOperator('Contains', TransactionExploreConditionOperatorType.Contains);
    public static readonly NotContains = new TransactionExploreConditionOperator('Not contains', TransactionExploreConditionOperatorType.NotContains);
    public static readonly StartsWith = new TransactionExploreConditionOperator('Starts with', TransactionExploreConditionOperatorType.StartsWith);
    public static readonly NotStartsWith = new TransactionExploreConditionOperator('Not starts with', TransactionExploreConditionOperatorType.NotStartsWith);
    public static readonly EndsWith = new TransactionExploreConditionOperator('Ends with', TransactionExploreConditionOperatorType.EndsWith);
    public static readonly NotEndsWith = new TransactionExploreConditionOperator('Not ends with', TransactionExploreConditionOperatorType.NotEndsWith);

    public readonly name: string;
    public readonly value: TransactionExploreConditionOperatorType;

    private constructor(name: string, value: TransactionExploreConditionOperatorType) {
        this.name = name;
        this.value = value;

        TransactionExploreConditionOperator.allInstances.push(this);
        TransactionExploreConditionOperator.allInstancesByValue[value] = this;
    }

    public static values(): TransactionExploreConditionOperator[] {
        return TransactionExploreConditionOperator.allInstances;
    }

    public static valueOf(type: string): TransactionExploreConditionOperator | undefined {
        return TransactionExploreConditionOperator.allInstancesByValue[type];
    }
}

export const DEFAULT_TRANSACTION_EXPLORE_DATE_RANGE: DateRange = DateRange.ThisMonth;
