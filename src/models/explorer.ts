import { type PartialRecord, itemAndIndex, keysIfValueEquals } from '@/core/base.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    TransactionExplorerConditionRelation,
    TransactionExplorerConditionRelationPriority,
    TransactionExplorerConditionFieldType,
    TransactionExplorerConditionField,
    TransactionExplorerConditionOperatorType,
    TransactionExplorerConditionOperator
} from '@/core/explorer.ts';

import { Account } from '@/models/account.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { type TransactionInsightDataItem } from '@/models/transaction.ts';

interface ExpressionNode {
    textualExpression: string;
    operator?: TransactionExplorerConditionRelation;
}

export class TransactionExplorerQuery {
    public name: string;
    public conditions: TransactionExplorerConditionWithRelation[];

    private constructor(name: string, conditions: TransactionExplorerConditionWithRelation[]) {
        this.name = name;
        this.conditions = conditions;
    }

    public addNewCondition(field: TransactionExplorerConditionField, isFirst: boolean): TransactionExplorerConditionWithRelation {
        let condition: TransactionExplorerCondition;

        switch (field) {
            case TransactionExplorerConditionField.TransactionType:
                condition = new TransactionExplorerTransactionTypeCondition([ TransactionType.Expense, TransactionType.Income, TransactionType.Transfer ]);
                break;
            case TransactionExplorerConditionField.TransactionCategory:
                condition = new TransactionExplorerTransactionCategoryCondition([]);
                break;
            case TransactionExplorerConditionField.SourceAccount:
                condition = new TransactionExplorerSourceAccountCondition([]);
                break;
            case TransactionExplorerConditionField.DestinationAccount:
                condition = new TransactionExplorerDestinationAccountCondition([]);
                break;
            case TransactionExplorerConditionField.SourceAmount:
                condition = new TransactionExplorerSourceAmountCondition(TransactionExplorerConditionOperatorType.Between, [0, 0]);
                break;
            case TransactionExplorerConditionField.DestinationAmount:
                condition = new TransactionExplorerDestinationAmountCondition(TransactionExplorerConditionOperatorType.Between, [0, 0]);
                break;
            case TransactionExplorerConditionField.TransactionTag:
                condition = new TransactionExplorerTransactionTagCondition(TransactionExplorerConditionOperatorType.HasAny, []);
                break;
            case TransactionExplorerConditionField.Description:
                condition = new TransactionExplorerDescriptionCondition(TransactionExplorerConditionOperatorType.Contains, '');
                break;
            default:
                condition = new TransactionExplorerTransactionTypeCondition([ TransactionType.Expense, TransactionType.Income, TransactionType.Transfer ]);
                break;
        }

        return new TransactionExplorerConditionWithRelation(
            condition,
            isFirst ? TransactionExplorerConditionRelation.First : TransactionExplorerConditionRelation.And
        );
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        if (!this.conditions || this.conditions.length < 1) {
            return true;
        }

        const postfixExprTokens = this.getPostfixExprTokens();
        const stack: boolean[] = [];

        for (const token of postfixExprTokens) {
            if (token === TransactionExplorerConditionRelation.And || token === TransactionExplorerConditionRelation.Or) {
                const right = stack.pop();
                const left = stack.pop();

                if (left === undefined || right === undefined) {
                    throw new Error('invalid postfix expression');
                }

                if (token === TransactionExplorerConditionRelation.And) {
                    stack.push(left && right);
                } else if (token === TransactionExplorerConditionRelation.Or) {
                    stack.push(left || right);
                } else {
                    throw new Error('invalid postfix expression');
                }
            } else {
                stack.push(token.match(transaction));
            }
        }

        if (stack.length !== 1) {
            throw new Error('invalid postfix evaluation result');
        }

        return stack[0] as boolean;
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string {
        if (!this.conditions || this.conditions.length < 1) {
            return '';
        }

        const postfixExprTokens = this.getPostfixExprTokens();
        const stack: ExpressionNode[] = [];

        for (const token of postfixExprTokens) {
            if (token === TransactionExplorerConditionRelation.And || token === TransactionExplorerConditionRelation.Or) {
                const right = stack.pop();
                const left = stack.pop();

                if (left === undefined || right === undefined) {
                    throw new Error('invalid postfix expression');
                }

                let leftExpression = left.textualExpression;
                let rightExpression = right.textualExpression;

                if (left.operator && left.operator !== token) {
                    leftExpression = `(${leftExpression})`;
                }

                if (right.operator && right.operator !== token) {
                    rightExpression = `(${rightExpression})`;
                }

                stack.push({
                    textualExpression: `${leftExpression} ${token.toUpperCase()} ${rightExpression}`,
                    operator: token
                });
            } else {
                stack.push({
                    textualExpression: token.toExpression(allCategoriesMap, allAccountsMap, allTagsMap)
                });
            }
        }

        if (stack.length !== 1) {
            throw new Error('invalid postfix evaluation result');
        }

        const finalNode = stack[0];

        if (!finalNode) {
            throw new Error('invalid postfix evaluation result');
        }

        return finalNode.textualExpression;
    }

    public getPostfixExprTokens(): (TransactionExplorerCondition | TransactionExplorerConditionRelation.And | TransactionExplorerConditionRelation.Or)[] {
        const finalTokens: (TransactionExplorerCondition | TransactionExplorerConditionRelation.And | TransactionExplorerConditionRelation.Or)[] = [];

        if (this.conditions.length < 1) {
            return finalTokens;
        }

        const operatorStack: TransactionExplorerConditionRelation[] = [];
        const firstCondition = this.conditions[0] as TransactionExplorerConditionWithRelation;

        if (firstCondition.relation !== TransactionExplorerConditionRelation.First) {
            throw new Error('the first condition must have relation "first"');
        }

        finalTokens.push(firstCondition.condition);

        for (const [item, index] of itemAndIndex(this.conditions)) {
            if (index < 1) {
                continue;
            }

            if (item.relation === TransactionExplorerConditionRelation.First) {
                throw new Error('only the first condition can have relation "first"');
            }

            const currentOperator = item.relation;

            while (operatorStack.length > 0) {
                const topOperator = operatorStack[operatorStack.length - 1];
                const isAndOrOperator = topOperator === TransactionExplorerConditionRelation.And || topOperator === TransactionExplorerConditionRelation.Or;

                if (isAndOrOperator && TransactionExplorerConditionRelationPriority[topOperator] >= TransactionExplorerConditionRelationPriority[currentOperator]) {
                    finalTokens.push(topOperator);
                    operatorStack.pop();
                } else {
                    break;
                }
            }

            operatorStack.push(currentOperator);
            finalTokens.push(item.condition);
        }

        while (operatorStack.length > 0) {
            const topOperator = operatorStack.pop();

            if (topOperator !== TransactionExplorerConditionRelation.And && topOperator !== TransactionExplorerConditionRelation.Or) {
                throw new Error('invalid operator in stack');
            }

            finalTokens.push(topOperator);
        }

        return finalTokens;
    }

    public clone(): TransactionExplorerQuery {
        const clonedConditions: TransactionExplorerConditionWithRelation[] = [];

        for (const condition of this.conditions) {
            const clonedCondition = TransactionExplorerConditionWithRelation.parse(condition.toJsonObject());

            if (!clonedCondition) {
                continue;
            }

            clonedConditions.push(clonedCondition);
        }

        return new TransactionExplorerQuery(this.name, clonedConditions);
    }

    public toJson(): string {
        return JSON.stringify({
            name: this.name,
            conditions: this.conditions.map(condition => condition.toJsonObject())
        });
    }

    public static create(): TransactionExplorerQuery {
        return new TransactionExplorerQuery("", []);
    }

    public static parse(json: string): TransactionExplorerQuery | null {
        const parsed = JSON.parse(json);
        const nameFieldValue = parsed['name'] as unknown;
        const conditionsFieldValue = parsed['conditions'] as unknown;

        if (typeof nameFieldValue !== 'string' || !Array.isArray(conditionsFieldValue)) {
            return null;
        }

        const name: string = nameFieldValue;
        const conditions: TransactionExplorerConditionWithRelation[] = [];

        for (const [item, index] of itemAndIndex(conditionsFieldValue)) {
            const condition = TransactionExplorerConditionWithRelation.parse(item);

            if (condition === null) {
                return null;
            }

            if (index === 0 && condition.relation !== TransactionExplorerConditionRelation.First) {
                return null;
            } else if (index > 0 && condition.relation === TransactionExplorerConditionRelation.First) {
                return null;
            }

            conditions.push(condition);
        }

        return new TransactionExplorerQuery(name, conditions);
    }
}

export class TransactionExplorerConditionWithRelation {
    public condition: TransactionExplorerCondition;
    public relation: TransactionExplorerConditionRelation;

    constructor(condition: TransactionExplorerCondition, relation: TransactionExplorerConditionRelation) {
        this.condition = condition;
        this.relation = relation;
    }

    public getSupportedOperators(): TransactionExplorerConditionOperator[] {
        let operatorTypes: PartialRecord<TransactionExplorerConditionOperatorType, true> = {};

        switch (this.condition.field) {
            case TransactionExplorerConditionField.TransactionType.value:
                operatorTypes = TransactionExplorerTransactionTypeCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.TransactionCategory.value:
                operatorTypes = TransactionExplorerTransactionCategoryCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.SourceAccount.value:
                operatorTypes = TransactionExplorerSourceAccountCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.DestinationAccount.value:
                operatorTypes = TransactionExplorerDestinationAccountCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.SourceAmount.value:
                operatorTypes = TransactionExplorerSourceAmountCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.DestinationAmount.value:
                operatorTypes = TransactionExplorerDestinationAmountCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.TransactionTag.value:
                operatorTypes = TransactionExplorerTransactionTagCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.Description.value:
                operatorTypes = TransactionExplorerDescriptionCondition.supportedOperators;
                break;
            default:
                return [];
        }

        const result: TransactionExplorerConditionOperator[] = [];

        for (const key of keysIfValueEquals(operatorTypes, true)) {
            const operator = TransactionExplorerConditionOperator.valueOf(key);

            if (operator) {
                result.push(operator);
            }
        }

        return result;
    }

    public toJsonObject(): unknown {
        return {
            condition: {
                field: this.condition.field,
                operator: this.condition.operator,
                value: this.condition.getValueForStore()
            },
            relation: this.relation
        };
    }

    public static parse(data: unknown): TransactionExplorerConditionWithRelation | null {
        if (typeof data !== 'object' || !data || !('condition' in data) || !('relation' in data)) {
            return null;
        }

        const conditionObject = data['condition'];
        const relation = data['relation'];

        if (typeof conditionObject !== 'object' || !conditionObject || !('field' in conditionObject) || !('operator' in conditionObject) || !('value' in conditionObject) || typeof relation !== 'string') {
            return null;
        }

        const conditionField = conditionObject['field'];
        const conditionOperator = conditionObject['operator'] as TransactionExplorerConditionOperatorType;
        const conditionValue = conditionObject['value'];

        let condition: TransactionExplorerCondition | null = null;

        switch (conditionField) {
            case TransactionExplorerConditionField.TransactionType.value:
                if (conditionOperator === TransactionExplorerConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExplorerTransactionTypeCondition(conditionValue as number[]);
                }
                break;
            case TransactionExplorerConditionField.TransactionCategory.value:
                if (conditionOperator === TransactionExplorerConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExplorerTransactionCategoryCondition(conditionValue as string[]);
                }
                break;
            case TransactionExplorerConditionField.SourceAccount.value:
                if (conditionOperator === TransactionExplorerConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExplorerSourceAccountCondition(conditionValue as string[]);
                }
                break;
            case TransactionExplorerConditionField.DestinationAccount.value:
                if (conditionOperator === TransactionExplorerConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExplorerDestinationAccountCondition(conditionValue as string[]);
                }
                break;
            case TransactionExplorerConditionField.SourceAmount.value:
                if (TransactionExplorerSourceAmountCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue) && conditionValue.length === 2) {
                    condition = new TransactionExplorerSourceAmountCondition(conditionOperator as AmountConditionOperator, conditionValue as [number, number]);
                }
                break;
            case TransactionExplorerConditionField.DestinationAmount.value:
                if (TransactionExplorerDestinationAmountCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue) && conditionValue.length === 2) {
                    condition = new TransactionExplorerDestinationAmountCondition(conditionOperator as AmountConditionOperator, conditionValue as [number, number]);
                }
                break;
            case TransactionExplorerConditionField.TransactionTag.value:
                if (TransactionExplorerTransactionTagCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue)) {
                    condition = new TransactionExplorerTransactionTagCondition(conditionOperator as TransactionTagConditionOperator, conditionValue as string[]);
                }
                break;
            case TransactionExplorerConditionField.Description.value:
                if (TransactionExplorerDescriptionCondition.supportedOperators[conditionOperator] && typeof conditionValue === 'string') {
                    condition = new TransactionExplorerDescriptionCondition(conditionOperator as DescriptionConditionOperator, conditionValue);
                }
                break;
            default:
                break;
        }

        if (condition === null) {
            return null;
        }

        if (relation !== TransactionExplorerConditionRelation.First && relation !== TransactionExplorerConditionRelation.And && relation !== TransactionExplorerConditionRelation.Or) {
            return null;
        }

        return new TransactionExplorerConditionWithRelation(condition, relation);
    }
}

export interface TransactionExplorerCondition<T = TransactionExplorerConditionFieldType, V = string | string[] | number[]> {
    readonly field: T;
    readonly operator: TransactionExplorerConditionOperatorType;
    value: V;

    getValueForStore(): V;
    match(transaction: TransactionInsightDataItem): boolean;
    toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string;
}

export class TransactionExplorerTransactionTypeCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.TransactionType, number[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.In]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.TransactionType;
    public readonly operator: TransactionExplorerConditionOperatorType.In = TransactionExplorerConditionOperatorType.In;
    public value: number[];

    constructor(value: number[]) {
        this.value = value;
    }

    public getValueForStore(): number[] {
        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return this.value.includes(transaction.type);
    }

    public toExpression(): string {
        const textualTypes = this.value.map(type => {
            if (type === TransactionType.Income) {
                return 'Income';
            } else if (type === TransactionType.Expense) {
                return 'Expense';
            } else if (type === TransactionType.Transfer) {
                return 'Transfer';
            } else {
                return type.toString();
            }
        }).join(', ');
        return `type IN (${textualTypes})`;
    }
}

export class TransactionExplorerTransactionCategoryCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.TransactionCategory, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.In]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.TransactionCategory;
    public readonly operator: TransactionExplorerConditionOperatorType.In = TransactionExplorerConditionOperatorType.In;
    public value: string[];

    constructor(value: string[]) {
        this.value = value;
    }

    public getValueForStore(): string[] {
        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return this.value.includes(transaction.primaryCategory?.id ?? '') || this.value.includes(transaction.secondaryCategory?.id ?? transaction.categoryId);
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>): string {
        const textualCategories = this.value.map(id => {
            const category = allCategoriesMap[id];

            if (category) {
                if (!category.parentId || category.parentId === '0') {
                    return '';
                }

                return `'${category.name}'`;
            } else {
                return `'${id}'`;
            }
        }).filter(item => !!item).join(', ');
        return `category IN (${textualCategories})`;
    }
}

export class TransactionExplorerSourceAccountCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.SourceAccount, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.In]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.SourceAccount;
    public readonly operator: TransactionExplorerConditionOperatorType.In = TransactionExplorerConditionOperatorType.In;
    public value: string[];

    constructor(value: string[]) {
        this.value = value;
    }

    public getValueForStore(): string[] {
        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return this.value.includes(transaction.sourceAccountId);
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>): string {
        const textualAccounts = this.value.map(id => {
            const account = allAccountsMap[id];

            if (account) {
                if (account.type === AccountType.MultiSubAccounts.type) {
                    return '';
                }

                return `'${account.name}'`;
            } else {
                return `'${id}'`;
            }
        }).filter(item => !!item).join(', ');
        return `source_account IN (${textualAccounts})`;
    }
}

export class TransactionExplorerDestinationAccountCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.DestinationAccount, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.In]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.DestinationAccount;
    public readonly operator: TransactionExplorerConditionOperatorType.In = TransactionExplorerConditionOperatorType.In;
    public value: string[];

    constructor(value: string[]) {
        this.value = value;
    }

    public getValueForStore(): string[] {
        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return this.value.includes(transaction.destinationAccountId);
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>): string {
        const textualAccounts = this.value.map(id => {
            const account = allAccountsMap[id];

            if (account) {
                if (account.type === AccountType.MultiSubAccounts.type) {
                    return '';
                }

                return `'${account.name}'`;
            } else {
                return `'${id}'`;
            }
        }).filter(item => !!item).join(', ');
        return `destination_account IN (${textualAccounts})`;
    }
}

type AmountConditionField = TransactionExplorerConditionFieldType.SourceAmount | TransactionExplorerConditionFieldType.DestinationAmount;
type AmountConditionOperator = TransactionExplorerConditionOperatorType.Equals |
    TransactionExplorerConditionOperatorType.NotEquals |
    TransactionExplorerConditionOperatorType.GreaterThan |
    TransactionExplorerConditionOperatorType.LessThan |
    TransactionExplorerConditionOperatorType.Between |
    TransactionExplorerConditionOperatorType.NotBetween;

export abstract class AbstractTransactionExplorerAmountCondition<T = AmountConditionField> implements TransactionExplorerCondition<T, [number, number]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.Equals]: true,
        [TransactionExplorerConditionOperatorType.NotEquals]: true,
        [TransactionExplorerConditionOperatorType.GreaterThan]: true,
        [TransactionExplorerConditionOperatorType.LessThan]: true,
        [TransactionExplorerConditionOperatorType.Between]: true,
        [TransactionExplorerConditionOperatorType.NotBetween]: true
    };
    public abstract readonly field: T;
    public readonly operator: AmountConditionOperator = TransactionExplorerConditionOperatorType.Between;
    public value: [number, number];

    protected constructor(operator: AmountConditionOperator, value: [number, number]) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): [number, number] {
        if (this.operator === TransactionExplorerConditionOperatorType.Between || this.operator === TransactionExplorerConditionOperatorType.NotBetween) {
            return [this.value[0], this.value[1]];
        } else {
            return [this.value[0], this.value[0]];
        }
    }

    public abstract match(transaction: TransactionInsightDataItem): boolean;
    public abstract toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string;

    protected matchAmount(amount: number): boolean {
        switch (this.operator) {
            case TransactionExplorerConditionOperatorType.GreaterThan:
                return amount > this.value[0];
            case TransactionExplorerConditionOperatorType.LessThan:
                return amount < this.value[0];
            case TransactionExplorerConditionOperatorType.Equals:
                return amount === this.value[0];
            case TransactionExplorerConditionOperatorType.NotEquals:
                return amount !== this.value[0];
            case TransactionExplorerConditionOperatorType.Between:
                return amount >= this.value[0] && amount <= this.value[1];
            case TransactionExplorerConditionOperatorType.NotBetween:
                return amount < this.value[0] || amount > this.value[1];
            default:
                return false;
        }
    }

    protected getExpression(amountFieldName: string): string {
        let expressionAmount1 = this.value[0].toString(10);
        let expressionAmount2 = this.value[1].toString(10);

        if (expressionAmount1.length > 2) {
            expressionAmount1 = `${expressionAmount1.substring(0, expressionAmount1.length - 2)}.${expressionAmount1.substring(expressionAmount1.length - 2)}`;
        }

        if (expressionAmount2.length > 2) {
            expressionAmount2 = `${expressionAmount2.substring(0, expressionAmount2.length - 2)}.${expressionAmount2.substring(expressionAmount2.length - 2)}`;
        }

        switch (this.operator) {
            case TransactionExplorerConditionOperatorType.GreaterThan:
                return `${amountFieldName} > ${expressionAmount1}`;
            case TransactionExplorerConditionOperatorType.LessThan:
                return `${amountFieldName} < ${expressionAmount1}`;
            case TransactionExplorerConditionOperatorType.Equals:
                return `${amountFieldName} = ${expressionAmount1}`;
            case TransactionExplorerConditionOperatorType.NotEquals:
                return `${amountFieldName} <> ${expressionAmount1}`;
            case TransactionExplorerConditionOperatorType.Between:
                return `(${amountFieldName} >= ${expressionAmount1} AND ${amountFieldName} <= ${expressionAmount2})`;
            case TransactionExplorerConditionOperatorType.NotBetween:
                return `(${amountFieldName} < ${expressionAmount1} OR ${amountFieldName} > ${expressionAmount2})`;
            default:
                return '';
        }
    }
}

export class TransactionExplorerSourceAmountCondition extends AbstractTransactionExplorerAmountCondition<TransactionExplorerConditionFieldType.SourceAmount> {
    public readonly field = TransactionExplorerConditionFieldType.SourceAmount;

    constructor(operator: AmountConditionOperator, value: [number, number]) {
        super(operator, value);
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return super.matchAmount(transaction.sourceAmount);
    }

    public toExpression(): string {
        return this.getExpression('source_amount');
    }
}

export class TransactionExplorerDestinationAmountCondition extends AbstractTransactionExplorerAmountCondition<TransactionExplorerConditionFieldType.DestinationAmount> {
    public readonly field = TransactionExplorerConditionFieldType.DestinationAmount;

    constructor(operator: AmountConditionOperator, value: [number, number]) {
        super(operator, value);
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return super.matchAmount(transaction.destinationAmount);
    }

    public toExpression(): string {
        return this.getExpression('destination_amount');
    }
}

type TransactionTagConditionOperator = TransactionExplorerConditionOperatorType.IsEmpty |
    TransactionExplorerConditionOperatorType.IsNotEmpty |
    TransactionExplorerConditionOperatorType.Equals |
    TransactionExplorerConditionOperatorType.NotEquals |
    TransactionExplorerConditionOperatorType.HasAny |
    TransactionExplorerConditionOperatorType.HasAll |
    TransactionExplorerConditionOperatorType.NotHasAny |
    TransactionExplorerConditionOperatorType.NotHasAll;

export class TransactionExplorerTransactionTagCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.TransactionTag, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.IsEmpty]: true,
        [TransactionExplorerConditionOperatorType.IsNotEmpty]: true,
        [TransactionExplorerConditionOperatorType.Equals]: true,
        [TransactionExplorerConditionOperatorType.NotEquals]: true,
        [TransactionExplorerConditionOperatorType.HasAny]: true,
        [TransactionExplorerConditionOperatorType.HasAll]: true,
        [TransactionExplorerConditionOperatorType.NotHasAny]: true,
        [TransactionExplorerConditionOperatorType.NotHasAll]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.TransactionTag;
    public readonly operator: TransactionTagConditionOperator = TransactionExplorerConditionOperatorType.HasAny;
    public value: string[];

    constructor(operator: TransactionTagConditionOperator, value: string[]) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): string[] {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty || this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return [];
        }

        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        const transactionTags: Record<string, boolean> = {};

        for (const tagId of transaction.tagIds) {
            transactionTags[tagId] = true;
        }

        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty || this.value.length < 1) {
            return transaction.tagIds.length < 1;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return transaction.tagIds.length > 0;
        } else if (this.operator === TransactionExplorerConditionOperatorType.Equals || this.operator === TransactionExplorerConditionOperatorType.NotEquals) {
            let hasAll = true;

            for (const tagId of this.value) {
                if (!transactionTags[tagId]) {
                    hasAll = false;
                    break;
                }
            }

            const hasSameCount = transaction.tagIds.length === this.value.length;

            if (this.operator === TransactionExplorerConditionOperatorType.Equals && hasAll && hasSameCount) {
                return true;
            } else if (this.operator === TransactionExplorerConditionOperatorType.NotEquals && (!hasAll || !hasSameCount)) {
                return true;
            }
        } else if (this.operator === TransactionExplorerConditionOperatorType.HasAny || this.operator === TransactionExplorerConditionOperatorType.NotHasAny) {
            let hasAny = false;

            for (const tagId of this.value) {
                if (transactionTags[tagId]) {
                    hasAny = true;
                    break;
                }
            }

            if (this.operator === TransactionExplorerConditionOperatorType.HasAny && hasAny) {
                return true;
            } else if (this.operator === TransactionExplorerConditionOperatorType.NotHasAny && !hasAny) {
                return true;
            }
        } else if (this.operator === TransactionExplorerConditionOperatorType.HasAll || this.operator === TransactionExplorerConditionOperatorType.NotHasAll) {
            let hasAll = true;

            for (const tagId of this.value) {
                if (!transactionTags[tagId]) {
                    hasAll = false;
                    break;
                }
            }

            if (this.operator === TransactionExplorerConditionOperatorType.HasAll && hasAll) {
                return true;
            } else if (this.operator === TransactionExplorerConditionOperatorType.NotHasAll && !hasAll) {
                return true;
            }
        }

        return false;
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return `tags IS EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return `tags IS NOT EMPTY`;
        }

        const textualTags = this.value.map(id => {
            const tag = allTagsMap[id];

            if (tag) {
                return `'${tag.name}'`;
            } else {
                return `'${id}'`;
            }
        }).join(', ');

        if (this.operator === TransactionExplorerConditionOperatorType.Equals) {
            return `tags FULL MATCHES (${textualTags})`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotEquals) {
            return `tags NOT FULL MATCHES (${textualTags})`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.HasAny) {
            return `tags HAS ANY (${textualTags})`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.HasAll) {
            return `tags HAS ALL (${textualTags})`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotHasAny) {
            return `tags NOT HAS ANY (${textualTags})`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotHasAll) {
            return `tags NOT HAS ALL (${textualTags})`;
        } else {
            return '';
        }
    }
}

type DescriptionConditionOperator = TransactionExplorerConditionOperatorType.IsEmpty |
    TransactionExplorerConditionOperatorType.IsNotEmpty |
    TransactionExplorerConditionOperatorType.Equals |
    TransactionExplorerConditionOperatorType.NotEquals |
    TransactionExplorerConditionOperatorType.Contains |
    TransactionExplorerConditionOperatorType.NotContains |
    TransactionExplorerConditionOperatorType.StartsWith |
    TransactionExplorerConditionOperatorType.NotStartsWith |
    TransactionExplorerConditionOperatorType.EndsWith |
    TransactionExplorerConditionOperatorType.NotEndsWith;

export class TransactionExplorerDescriptionCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.Description, string> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.IsEmpty]: true,
        [TransactionExplorerConditionOperatorType.IsNotEmpty]: true,
        [TransactionExplorerConditionOperatorType.Equals]: true,
        [TransactionExplorerConditionOperatorType.NotEquals]: true,
        [TransactionExplorerConditionOperatorType.Contains]: true,
        [TransactionExplorerConditionOperatorType.NotContains]: true,
        [TransactionExplorerConditionOperatorType.StartsWith]: true,
        [TransactionExplorerConditionOperatorType.NotStartsWith]: true,
        [TransactionExplorerConditionOperatorType.EndsWith]: true,
        [TransactionExplorerConditionOperatorType.NotEndsWith]: true
    };
    public readonly field = TransactionExplorerConditionFieldType.Description;
    public readonly operator: DescriptionConditionOperator = TransactionExplorerConditionOperatorType.Contains;
    public value: string;

    constructor(operator: DescriptionConditionOperator, value: string) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): string {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty || this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return '';
        }

        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        const description = transaction.comment || '';

        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return description.length === 0;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return description.length > 0;
        } else if (this.operator === TransactionExplorerConditionOperatorType.Equals) {
            return description === this.value;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotEquals) {
            return description !== this.value;
        } else if (this.operator === TransactionExplorerConditionOperatorType.Contains) {
            return description.includes(this.value);
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotContains) {
            return !description.includes(this.value);
        } else if (this.operator === TransactionExplorerConditionOperatorType.StartsWith) {
            return description.startsWith(this.value);
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotStartsWith) {
            return !description.startsWith(this.value);
        } else if (this.operator === TransactionExplorerConditionOperatorType.EndsWith) {
            return description.endsWith(this.value);
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotEndsWith) {
            return !description.endsWith(this.value);
        }

        return false;
    }

    public toExpression(): string {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return `description IS EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return `description IS NOT EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.Equals) {
            return `description = '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotEquals) {
            return `description <> '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.Contains) {
            return `description CONTAINS '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotContains) {
            return `description NOT CONTAINS '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.StartsWith) {
            return `description STARTS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotStartsWith) {
            return `description NOT STARTS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.EndsWith) {
            return `description ENDS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.NotEndsWith) {
            return `description NOT ENDS WITH '${this.value.replace(/'/g, "''")}'`;
        }

        return '';
    }
}
