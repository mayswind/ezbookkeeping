import type { PartialRecord, NameValue } from './base.ts';

export type ImportTransactionReplaceRuleConditionValue = string | [number, number];
export type ImportTransactionReplaceRuleTargetValue = string | number;
export type ImportTransactionReplaceRuleApplyScope = 'all' | 'selected';

export enum ImportTransactionReplaceRuleConditionFieldType {
    ExpenseCategory = 'expenseCategory',
    IncomeCategory = 'incomeCategory',
    TransferCategory = 'transferCategory',
    SourceAccount = 'sourceAccount',
    DestinationAccount = 'destinationAccount',
    Tag = 'tag',
    Amount = 'amount',
    TransferInAmount = 'transferInAmount',
    Description = 'description',
    DescriptionCaseInsensitive = 'descriptionCaseInsensitive',
    DescriptionNormalized = 'descriptionNormalized'
}

export class ImportTransactionReplaceRuleConditionField implements NameValue {
    private static readonly allInstances: ImportTransactionReplaceRuleConditionField[] = [];
    private static readonly allInstancesByValue: PartialRecord<ImportTransactionReplaceRuleConditionFieldType, ImportTransactionReplaceRuleConditionField> = {};

    public static readonly ExpenseCategory = new ImportTransactionReplaceRuleConditionField('Expense Category', ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory);
    public static readonly IncomeCategory = new ImportTransactionReplaceRuleConditionField('Income Category', ImportTransactionReplaceRuleConditionFieldType.IncomeCategory);
    public static readonly TransferCategory = new ImportTransactionReplaceRuleConditionField('Transfer Category', ImportTransactionReplaceRuleConditionFieldType.TransferCategory);
    public static readonly SourceAccount = new ImportTransactionReplaceRuleConditionField('Source Account', ImportTransactionReplaceRuleConditionFieldType.SourceAccount);
    public static readonly DestinationAccount = new ImportTransactionReplaceRuleConditionField('Destination Account', ImportTransactionReplaceRuleConditionFieldType.DestinationAccount);
    public static readonly Tag = new ImportTransactionReplaceRuleConditionField('Transaction Tag', ImportTransactionReplaceRuleConditionFieldType.Tag);
    public static readonly Amount = new ImportTransactionReplaceRuleConditionField('Amount', ImportTransactionReplaceRuleConditionFieldType.Amount);
    public static readonly TransferInAmount = new ImportTransactionReplaceRuleConditionField('Transfer In Amount', ImportTransactionReplaceRuleConditionFieldType.TransferInAmount);
    public static readonly Description = new ImportTransactionReplaceRuleConditionField('Description (Exact)', ImportTransactionReplaceRuleConditionFieldType.Description);
    public static readonly DescriptionCaseInsensitive = new ImportTransactionReplaceRuleConditionField('Description (Ignore Case)', ImportTransactionReplaceRuleConditionFieldType.DescriptionCaseInsensitive);
    public static readonly DescriptionNormalized = new ImportTransactionReplaceRuleConditionField('Description (Normalized)', ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized);

    public readonly name: string;
    public readonly value: ImportTransactionReplaceRuleConditionFieldType;

    private constructor(name: string, value: ImportTransactionReplaceRuleConditionFieldType) {
        this.name = name;
        this.value = value;

        ImportTransactionReplaceRuleConditionField.allInstances.push(this);
        ImportTransactionReplaceRuleConditionField.allInstancesByValue[value] = this;
    }

    public static values(): ImportTransactionReplaceRuleConditionField[] {
        return ImportTransactionReplaceRuleConditionField.allInstances;
    }

    public static valueOf(value: ImportTransactionReplaceRuleConditionFieldType): ImportTransactionReplaceRuleConditionField | undefined {
        return ImportTransactionReplaceRuleConditionField.allInstancesByValue[value];
    }
}

export enum ImportTransactionReplaceRuleConditionOperatorType {
    Is = 'is',
    IsEmpty = 'isEmpty',
    IsNotEmpty = 'isNotEmpty',
    Equals = 'equals',
    NotEquals = 'notEquals',
    Contains = 'contains',
    NotContains = 'notContains',
    StartsWith = 'startsWith',
    NotStartsWith = 'notStartsWith',
    EndsWith = 'endsWith',
    NotEndsWith = 'notEndsWith',
    RegexMatch = 'regexMatch',
    NotRegexMatch = 'notRegexMatch',
    GreaterThan = 'greaterThan',
    LessThan = 'lessThan',
    Between = 'between',
    NotBetween = 'notBetween'
}

export class ImportTransactionReplaceRuleConditionOperator implements NameValue {
    private static readonly allInstancesByValue: PartialRecord<ImportTransactionReplaceRuleConditionOperatorType, ImportTransactionReplaceRuleConditionOperator> = {};

    public static readonly Is = new ImportTransactionReplaceRuleConditionOperator('Is', ImportTransactionReplaceRuleConditionOperatorType.Is);
    public static readonly IsEmpty = new ImportTransactionReplaceRuleConditionOperator('Is empty', ImportTransactionReplaceRuleConditionOperatorType.IsEmpty);
    public static readonly IsNotEmpty = new ImportTransactionReplaceRuleConditionOperator('Is not empty', ImportTransactionReplaceRuleConditionOperatorType.IsNotEmpty);
    public static readonly Equals = new ImportTransactionReplaceRuleConditionOperator('Equal to', ImportTransactionReplaceRuleConditionOperatorType.Equals);
    public static readonly NotEquals = new ImportTransactionReplaceRuleConditionOperator('Not equal to', ImportTransactionReplaceRuleConditionOperatorType.NotEquals);
    public static readonly Contains = new ImportTransactionReplaceRuleConditionOperator('Contains', ImportTransactionReplaceRuleConditionOperatorType.Contains);
    public static readonly NotContains = new ImportTransactionReplaceRuleConditionOperator('Does not contain', ImportTransactionReplaceRuleConditionOperatorType.NotContains);
    public static readonly StartsWith = new ImportTransactionReplaceRuleConditionOperator('Starts with', ImportTransactionReplaceRuleConditionOperatorType.StartsWith);
    public static readonly NotStartsWith = new ImportTransactionReplaceRuleConditionOperator('Does not start with', ImportTransactionReplaceRuleConditionOperatorType.NotStartsWith);
    public static readonly EndsWith = new ImportTransactionReplaceRuleConditionOperator('Ends with', ImportTransactionReplaceRuleConditionOperatorType.EndsWith);
    public static readonly NotEndsWith = new ImportTransactionReplaceRuleConditionOperator('Does not end with', ImportTransactionReplaceRuleConditionOperatorType.NotEndsWith);
    public static readonly RegexMatch = new ImportTransactionReplaceRuleConditionOperator('Matches regex', ImportTransactionReplaceRuleConditionOperatorType.RegexMatch);
    public static readonly NotRegexMatch = new ImportTransactionReplaceRuleConditionOperator('Does not match regex', ImportTransactionReplaceRuleConditionOperatorType.NotRegexMatch);
    public static readonly GreaterThan = new ImportTransactionReplaceRuleConditionOperator('Greater than', ImportTransactionReplaceRuleConditionOperatorType.GreaterThan);
    public static readonly LessThan = new ImportTransactionReplaceRuleConditionOperator('Less than', ImportTransactionReplaceRuleConditionOperatorType.LessThan);
    public static readonly Between = new ImportTransactionReplaceRuleConditionOperator('Between', ImportTransactionReplaceRuleConditionOperatorType.Between);
    public static readonly NotBetween = new ImportTransactionReplaceRuleConditionOperator('Not between', ImportTransactionReplaceRuleConditionOperatorType.NotBetween);

    public readonly name: string;
    public readonly value: ImportTransactionReplaceRuleConditionOperatorType;

    private constructor(name: string, value: ImportTransactionReplaceRuleConditionOperatorType) {
        this.name = name;
        this.value = value;

        ImportTransactionReplaceRuleConditionOperator.allInstancesByValue[value] = this;
    }

    public static valueOf(value: ImportTransactionReplaceRuleConditionOperatorType): ImportTransactionReplaceRuleConditionOperator | undefined {
        return ImportTransactionReplaceRuleConditionOperator.allInstancesByValue[value];
    }
}

export enum ImportTransactionReplaceRuleActionType {
    SetTransactionType = 'setTransactionType',
    SetExpenseCategory = 'setExpenseCategory',
    SetIncomeCategory = 'setIncomeCategory',
    SetTransferCategory = 'setTransferCategory',
    SetSourceAccount = 'setSourceAccount',
    SetDestinationAccount = 'setDestinationAccount',
    AddTag = 'addTag',
    ReplaceTag = 'replaceTag',
    DeleteTag = 'deleteTag',
    SetSourceAmount = 'setSourceAmount',
    SetDestinationAmount = 'setDestinationAmount',
    SetTimezone = 'setTimezone'
}

export class ImportTransactionReplaceRuleAction implements NameValue {
    private static readonly allInstances: ImportTransactionReplaceRuleAction[] = [];
    private static readonly allInstancesByValue: PartialRecord<ImportTransactionReplaceRuleActionType, ImportTransactionReplaceRuleAction> = {};

    public static readonly SetTransactionType = new ImportTransactionReplaceRuleAction('Set Transaction Type', ImportTransactionReplaceRuleActionType.SetTransactionType);
    public static readonly SetExpenseCategory = new ImportTransactionReplaceRuleAction('Set Expense Category', ImportTransactionReplaceRuleActionType.SetExpenseCategory);
    public static readonly SetIncomeCategory = new ImportTransactionReplaceRuleAction('Set Income Category', ImportTransactionReplaceRuleActionType.SetIncomeCategory);
    public static readonly SetTransferCategory = new ImportTransactionReplaceRuleAction('Set Transfer Category', ImportTransactionReplaceRuleActionType.SetTransferCategory);
    public static readonly SetSourceAccount = new ImportTransactionReplaceRuleAction('Set Source Account', ImportTransactionReplaceRuleActionType.SetSourceAccount);
    public static readonly SetDestinationAccount = new ImportTransactionReplaceRuleAction('Set Destination Account', ImportTransactionReplaceRuleActionType.SetDestinationAccount);
    public static readonly AddTag = new ImportTransactionReplaceRuleAction('Add Tag', ImportTransactionReplaceRuleActionType.AddTag);
    public static readonly ReplaceTag = new ImportTransactionReplaceRuleAction('Replace Tag', ImportTransactionReplaceRuleActionType.ReplaceTag);
    public static readonly DeleteTag = new ImportTransactionReplaceRuleAction('Delete Tag', ImportTransactionReplaceRuleActionType.DeleteTag);
    public static readonly SetSourceAmount = new ImportTransactionReplaceRuleAction('Set Source Amount', ImportTransactionReplaceRuleActionType.SetSourceAmount);
    public static readonly SetDestinationAmount = new ImportTransactionReplaceRuleAction('Set Destination Amount', ImportTransactionReplaceRuleActionType.SetDestinationAmount);
    public static readonly SetTimezone = new ImportTransactionReplaceRuleAction('Set Timezone', ImportTransactionReplaceRuleActionType.SetTimezone);

    public readonly name: string;
    public readonly value: ImportTransactionReplaceRuleActionType;

    private constructor(name: string, value: ImportTransactionReplaceRuleActionType) {
        this.name = name;
        this.value = value;

        ImportTransactionReplaceRuleAction.allInstances.push(this);
        ImportTransactionReplaceRuleAction.allInstancesByValue[value] = this;
    }

    public static values(): ImportTransactionReplaceRuleAction[] {
        return ImportTransactionReplaceRuleAction.allInstances;
    }

    public static valueOf(value: ImportTransactionReplaceRuleActionType): ImportTransactionReplaceRuleAction | undefined {
        return ImportTransactionReplaceRuleAction.allInstancesByValue[value];
    }
}

export const TEXT_CONDITION_FIELDS: PartialRecord<ImportTransactionReplaceRuleConditionFieldType, true> = {
    [ImportTransactionReplaceRuleConditionFieldType.Description]: true,
    [ImportTransactionReplaceRuleConditionFieldType.DescriptionCaseInsensitive]: true,
    [ImportTransactionReplaceRuleConditionFieldType.DescriptionNormalized]: true
};

export const AMOUNT_CONDITION_FIELDS: PartialRecord<ImportTransactionReplaceRuleConditionFieldType, true> = {
    [ImportTransactionReplaceRuleConditionFieldType.Amount]: true,
    [ImportTransactionReplaceRuleConditionFieldType.TransferInAmount]: true
};

export const ITEM_CONDITION_FIELDS: PartialRecord<ImportTransactionReplaceRuleConditionFieldType, true> = {
    [ImportTransactionReplaceRuleConditionFieldType.ExpenseCategory]: true,
    [ImportTransactionReplaceRuleConditionFieldType.IncomeCategory]: true,
    [ImportTransactionReplaceRuleConditionFieldType.TransferCategory]: true,
    [ImportTransactionReplaceRuleConditionFieldType.SourceAccount]: true,
    [ImportTransactionReplaceRuleConditionFieldType.DestinationAccount]: true,
    [ImportTransactionReplaceRuleConditionFieldType.Tag]: true
};

export const TEXT_CONDITION_OPERATORS: PartialRecord<ImportTransactionReplaceRuleConditionOperatorType, true> = {
    [ImportTransactionReplaceRuleConditionOperatorType.IsEmpty]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.IsNotEmpty]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.Equals]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotEquals]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.Contains]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotContains]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.StartsWith]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotStartsWith]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.EndsWith]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotEndsWith]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.RegexMatch]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotRegexMatch]: true
};

export const AMOUNT_CONDITION_OPERATORS: PartialRecord<ImportTransactionReplaceRuleConditionOperatorType, true> = {
    [ImportTransactionReplaceRuleConditionOperatorType.Equals]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotEquals]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.GreaterThan]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.LessThan]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.Between]: true,
    [ImportTransactionReplaceRuleConditionOperatorType.NotBetween]: true
};
