import { type PartialRecord, itemAndIndex, keysIfValueEquals } from '@/core/base.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import {
    TransactionExploreConditionRelation,
    TransactionExploreConditionRelationPriority,
    TransactionExploreConditionFieldType,
    TransactionExploreConditionField,
    TransactionExploreConditionOperatorType,
    TransactionExploreConditionOperator
} from '@/core/explore.ts';

import { Account } from '@/models/account.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { type TransactionInsightDataItem } from '@/models/transaction.ts';

interface ExpressionNode {
    textualExpression: string;
    operator?: TransactionExploreConditionRelation;
}

export class TransactionExploreQuery {
    public name: string;
    public conditions: TransactionExploreConditionWithRelation[];

    private constructor(name: string, conditions: TransactionExploreConditionWithRelation[]) {
        this.name = name;
        this.conditions = conditions;
    }

    public addNewCondition(field: TransactionExploreConditionField, isFirst: boolean): TransactionExploreConditionWithRelation {
        let condition: TransactionExploreCondition;

        switch (field) {
            case TransactionExploreConditionField.TransactionType:
                condition = new TransactionExploreTransactionTypeCondition([ TransactionType.Expense, TransactionType.Income, TransactionType.Transfer ]);
                break;
            case TransactionExploreConditionField.TransactionCategory:
                condition = new TransactionExploreTransactionCategoryCondition([]);
                break;
            case TransactionExploreConditionField.SourceAccount:
                condition = new TransactionExploreSourceAccountCondition([]);
                break;
            case TransactionExploreConditionField.DestinationAccount:
                condition = new TransactionExploreDestinationAccountCondition([]);
                break;
            case TransactionExploreConditionField.SourceAmount:
                condition = new TransactionExploreSourceAmountCondition(TransactionExploreConditionOperatorType.Between, [0, 0]);
                break;
            case TransactionExploreConditionField.DestinationAmount:
                condition = new TransactionExploreDestinationAmountCondition(TransactionExploreConditionOperatorType.Between, [0, 0]);
                break;
            case TransactionExploreConditionField.TransactionTag:
                condition = new TransactionExploreTransactionTagCondition(TransactionExploreConditionOperatorType.HasAny, []);
                break;
            case TransactionExploreConditionField.Description:
                condition = new TransactionExploreDescriptionCondition(TransactionExploreConditionOperatorType.Contains, '');
                break;
            default:
                condition = new TransactionExploreTransactionTypeCondition([ TransactionType.Expense, TransactionType.Income, TransactionType.Transfer ]);
                break;
        }

        return new TransactionExploreConditionWithRelation(
            condition,
            isFirst ? TransactionExploreConditionRelation.First : TransactionExploreConditionRelation.And
        );
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        if (!this.conditions || this.conditions.length < 1) {
            return true;
        }

        const postfixExprTokens = this.getPostfixExprTokens();
        const stack: boolean[] = [];

        for (const token of postfixExprTokens) {
            if (token === TransactionExploreConditionRelation.And || token === TransactionExploreConditionRelation.Or) {
                const right = stack.pop();
                const left = stack.pop();

                if (left === undefined || right === undefined) {
                    throw new Error('invalid postfix expression');
                }

                if (token === TransactionExploreConditionRelation.And) {
                    stack.push(left && right);
                } else if (token === TransactionExploreConditionRelation.Or) {
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
            if (token === TransactionExploreConditionRelation.And || token === TransactionExploreConditionRelation.Or) {
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

    public getPostfixExprTokens(): (TransactionExploreCondition | TransactionExploreConditionRelation.And | TransactionExploreConditionRelation.Or)[] {
        const finalTokens: (TransactionExploreCondition | TransactionExploreConditionRelation.And | TransactionExploreConditionRelation.Or)[] = [];

        if (this.conditions.length < 1) {
            return finalTokens;
        }

        const operatorStack: TransactionExploreConditionRelation[] = [];
        const firstCondition = this.conditions[0] as TransactionExploreConditionWithRelation;

        if (firstCondition.relation !== TransactionExploreConditionRelation.First) {
            throw new Error('the first condition must have relation "first"');
        }

        finalTokens.push(firstCondition.condition);

        for (const [item, index] of itemAndIndex(this.conditions)) {
            if (index < 1) {
                continue;
            }

            if (item.relation === TransactionExploreConditionRelation.First) {
                throw new Error('only the first condition can have relation "first"');
            }

            const currentOperator = item.relation;

            while (operatorStack.length > 0) {
                const topOperator = operatorStack[operatorStack.length - 1];
                const isAndOrOperator = topOperator === TransactionExploreConditionRelation.And || topOperator === TransactionExploreConditionRelation.Or;

                if (isAndOrOperator && TransactionExploreConditionRelationPriority[topOperator] >= TransactionExploreConditionRelationPriority[currentOperator]) {
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

            if (topOperator !== TransactionExploreConditionRelation.And && topOperator !== TransactionExploreConditionRelation.Or) {
                throw new Error('invalid operator in stack');
            }

            finalTokens.push(topOperator);
        }

        return finalTokens;
    }

    public clone(): TransactionExploreQuery {
        const clonedConditions: TransactionExploreConditionWithRelation[] = [];

        for (const condition of this.conditions) {
            const clonedCondition = TransactionExploreConditionWithRelation.parse(condition.toJsonObject());

            if (!clonedCondition) {
                continue;
            }

            clonedConditions.push(clonedCondition);
        }

        return new TransactionExploreQuery(this.name, clonedConditions);
    }

    public toJson(): string {
        return JSON.stringify({
            name: this.name,
            conditions: this.conditions.map(condition => condition.toJsonObject())
        });
    }

    public static create(): TransactionExploreQuery {
        return new TransactionExploreQuery("", []);
    }

    public static parse(json: string): TransactionExploreQuery | null {
        const parsed = JSON.parse(json);
        const nameFieldValue = parsed['name'] as unknown;
        const conditionsFieldValue = parsed['conditions'] as unknown;

        if (typeof nameFieldValue !== 'string' || !Array.isArray(conditionsFieldValue)) {
            return null;
        }

        const name: string = nameFieldValue;
        const conditions: TransactionExploreConditionWithRelation[] = [];

        for (const [item, index] of itemAndIndex(conditionsFieldValue)) {
            const condition = TransactionExploreConditionWithRelation.parse(item);

            if (condition === null) {
                return null;
            }

            if (index === 0 && condition.relation !== TransactionExploreConditionRelation.First) {
                return null;
            } else if (index > 0 && condition.relation === TransactionExploreConditionRelation.First) {
                return null;
            }

            conditions.push(condition);
        }

        return new TransactionExploreQuery(name, conditions);
    }
}

export class TransactionExploreConditionWithRelation {
    public condition: TransactionExploreCondition;
    public relation: TransactionExploreConditionRelation;

    constructor(condition: TransactionExploreCondition, relation: TransactionExploreConditionRelation) {
        this.condition = condition;
        this.relation = relation;
    }

    public getSupportedOperators(): TransactionExploreConditionOperator[] {
        let operatorTypes: PartialRecord<TransactionExploreConditionOperatorType, true> = {};

        switch (this.condition.field) {
            case TransactionExploreConditionField.TransactionType.value:
                operatorTypes = TransactionExploreTransactionTypeCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.TransactionCategory.value:
                operatorTypes = TransactionExploreTransactionCategoryCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.SourceAccount.value:
                operatorTypes = TransactionExploreSourceAccountCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.DestinationAccount.value:
                operatorTypes = TransactionExploreDestinationAccountCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.SourceAmount.value:
                operatorTypes = TransactionExploreSourceAmountCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.DestinationAmount.value:
                operatorTypes = TransactionExploreDestinationAmountCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.TransactionTag.value:
                operatorTypes = TransactionExploreTransactionTagCondition.supportedOperators;
                break;
            case TransactionExploreConditionField.Description.value:
                operatorTypes = TransactionExploreDescriptionCondition.supportedOperators;
                break;
            default:
                return [];
        }

        const result: TransactionExploreConditionOperator[] = [];

        for (const key of keysIfValueEquals(operatorTypes, true)) {
            const operator = TransactionExploreConditionOperator.valueOf(key);

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

    public static parse(data: unknown): TransactionExploreConditionWithRelation | null {
        if (typeof data !== 'object' || !data || !('condition' in data) || !('relation' in data)) {
            return null;
        }

        const conditionObject = data['condition'];
        const relation = data['relation'];

        if (typeof conditionObject !== 'object' || !conditionObject || !('field' in conditionObject) || !('operator' in conditionObject) || !('value' in conditionObject) || typeof relation !== 'string') {
            return null;
        }

        const conditionField = conditionObject['field'];
        const conditionOperator = conditionObject['operator'] as TransactionExploreConditionOperatorType;
        const conditionValue = conditionObject['value'];

        let condition: TransactionExploreCondition | null = null;

        switch (conditionField) {
            case TransactionExploreConditionField.TransactionType.value:
                if (conditionOperator === TransactionExploreConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExploreTransactionTypeCondition(conditionValue as number[]);
                }
                break;
            case TransactionExploreConditionField.TransactionCategory.value:
                if (conditionOperator === TransactionExploreConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExploreTransactionCategoryCondition(conditionValue as string[]);
                }
                break;
            case TransactionExploreConditionField.SourceAccount.value:
                if (conditionOperator === TransactionExploreConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExploreSourceAccountCondition(conditionValue as string[]);
                }
                break;
            case TransactionExploreConditionField.DestinationAccount.value:
                if (conditionOperator === TransactionExploreConditionOperatorType.In && Array.isArray(conditionValue)) {
                    condition = new TransactionExploreDestinationAccountCondition(conditionValue as string[]);
                }
                break;
            case TransactionExploreConditionField.SourceAmount.value:
                if (TransactionExploreSourceAmountCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue) && conditionValue.length === 2) {
                    condition = new TransactionExploreSourceAmountCondition(conditionOperator as AmountConditionOperator, conditionValue as [number, number]);
                }
                break;
            case TransactionExploreConditionField.DestinationAmount.value:
                if (TransactionExploreDestinationAmountCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue) && conditionValue.length === 2) {
                    condition = new TransactionExploreDestinationAmountCondition(conditionOperator as AmountConditionOperator, conditionValue as [number, number]);
                }
                break;
            case TransactionExploreConditionField.TransactionTag.value:
                if (TransactionExploreTransactionTagCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue)) {
                    condition = new TransactionExploreTransactionTagCondition(conditionOperator as TransactionTagConditionOperator, conditionValue as string[]);
                }
                break;
            case TransactionExploreConditionField.Description.value:
                if (TransactionExploreDescriptionCondition.supportedOperators[conditionOperator] && typeof conditionValue === 'string') {
                    condition = new TransactionExploreDescriptionCondition(conditionOperator as DescriptionConditionOperator, conditionValue);
                }
                break;
            default:
                break;
        }

        if (condition === null) {
            return null;
        }

        if (relation !== TransactionExploreConditionRelation.First && relation !== TransactionExploreConditionRelation.And && relation !== TransactionExploreConditionRelation.Or) {
            return null;
        }

        return new TransactionExploreConditionWithRelation(condition, relation);
    }
}

export interface TransactionExploreCondition<T = TransactionExploreConditionFieldType, V = string | string[] | number[]> {
    readonly field: T;
    readonly operator: TransactionExploreConditionOperatorType;
    value: V;

    getValueForStore(): V;
    match(transaction: TransactionInsightDataItem): boolean;
    toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string;
}

export class TransactionExploreTransactionTypeCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.TransactionType, number[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.In]: true
    };
    public readonly field = TransactionExploreConditionFieldType.TransactionType;
    public readonly operator: TransactionExploreConditionOperatorType.In = TransactionExploreConditionOperatorType.In;
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

export class TransactionExploreTransactionCategoryCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.TransactionCategory, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.In]: true
    };
    public readonly field = TransactionExploreConditionFieldType.TransactionCategory;
    public readonly operator: TransactionExploreConditionOperatorType.In = TransactionExploreConditionOperatorType.In;
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

export class TransactionExploreSourceAccountCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.SourceAccount, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.In]: true
    };
    public readonly field = TransactionExploreConditionFieldType.SourceAccount;
    public readonly operator: TransactionExploreConditionOperatorType.In = TransactionExploreConditionOperatorType.In;
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

export class TransactionExploreDestinationAccountCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.DestinationAccount, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.In]: true
    };
    public readonly field = TransactionExploreConditionFieldType.DestinationAccount;
    public readonly operator: TransactionExploreConditionOperatorType.In = TransactionExploreConditionOperatorType.In;
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

type AmountConditionField = TransactionExploreConditionFieldType.SourceAmount | TransactionExploreConditionFieldType.DestinationAmount;
type AmountConditionOperator = TransactionExploreConditionOperatorType.Equals |
    TransactionExploreConditionOperatorType.NotEquals |
    TransactionExploreConditionOperatorType.GreaterThan |
    TransactionExploreConditionOperatorType.LessThan |
    TransactionExploreConditionOperatorType.Between |
    TransactionExploreConditionOperatorType.NotBetween;

export abstract class AbstractTransactionExploreAmountCondition<T = AmountConditionField> implements TransactionExploreCondition<T, [number, number]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.Equals]: true,
        [TransactionExploreConditionOperatorType.NotEquals]: true,
        [TransactionExploreConditionOperatorType.GreaterThan]: true,
        [TransactionExploreConditionOperatorType.LessThan]: true,
        [TransactionExploreConditionOperatorType.Between]: true,
        [TransactionExploreConditionOperatorType.NotBetween]: true
    };
    public abstract readonly field: T;
    public readonly operator: AmountConditionOperator = TransactionExploreConditionOperatorType.Between;
    public value: [number, number];

    protected constructor(operator: AmountConditionOperator, value: [number, number]) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): [number, number] {
        if (this.operator === TransactionExploreConditionOperatorType.Between || this.operator === TransactionExploreConditionOperatorType.NotBetween) {
            return [this.value[0], this.value[1]];
        } else {
            return [this.value[0], this.value[0]];
        }
    }

    public abstract match(transaction: TransactionInsightDataItem): boolean;
    public abstract toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string;

    protected matchAmount(amount: number): boolean {
        switch (this.operator) {
            case TransactionExploreConditionOperatorType.GreaterThan:
                return amount > this.value[0];
            case TransactionExploreConditionOperatorType.LessThan:
                return amount < this.value[0];
            case TransactionExploreConditionOperatorType.Equals:
                return amount === this.value[0];
            case TransactionExploreConditionOperatorType.NotEquals:
                return amount !== this.value[0];
            case TransactionExploreConditionOperatorType.Between:
                return amount >= this.value[0] && amount <= this.value[1];
            case TransactionExploreConditionOperatorType.NotBetween:
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
            case TransactionExploreConditionOperatorType.GreaterThan:
                return `${amountFieldName} > ${expressionAmount1}`;
            case TransactionExploreConditionOperatorType.LessThan:
                return `${amountFieldName} < ${expressionAmount1}`;
            case TransactionExploreConditionOperatorType.Equals:
                return `${amountFieldName} = ${expressionAmount1}`;
            case TransactionExploreConditionOperatorType.NotEquals:
                return `${amountFieldName} <> ${expressionAmount1}`;
            case TransactionExploreConditionOperatorType.Between:
                return `(${amountFieldName} >= ${expressionAmount1} AND ${amountFieldName} <= ${expressionAmount2})`;
            case TransactionExploreConditionOperatorType.NotBetween:
                return `(${amountFieldName} < ${expressionAmount1} OR ${amountFieldName} > ${expressionAmount2})`;
            default:
                return '';
        }
    }
}

export class TransactionExploreSourceAmountCondition extends AbstractTransactionExploreAmountCondition<TransactionExploreConditionFieldType.SourceAmount> {
    public readonly field = TransactionExploreConditionFieldType.SourceAmount;

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

export class TransactionExploreDestinationAmountCondition extends AbstractTransactionExploreAmountCondition<TransactionExploreConditionFieldType.DestinationAmount> {
    public readonly field = TransactionExploreConditionFieldType.DestinationAmount;

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

type TransactionTagConditionOperator = TransactionExploreConditionOperatorType.IsEmpty |
    TransactionExploreConditionOperatorType.IsNotEmpty |
    TransactionExploreConditionOperatorType.Equals |
    TransactionExploreConditionOperatorType.NotEquals |
    TransactionExploreConditionOperatorType.HasAny |
    TransactionExploreConditionOperatorType.HasAll |
    TransactionExploreConditionOperatorType.NotHasAny |
    TransactionExploreConditionOperatorType.NotHasAll;

export class TransactionExploreTransactionTagCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.TransactionTag, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.IsEmpty]: true,
        [TransactionExploreConditionOperatorType.IsNotEmpty]: true,
        [TransactionExploreConditionOperatorType.Equals]: true,
        [TransactionExploreConditionOperatorType.NotEquals]: true,
        [TransactionExploreConditionOperatorType.HasAny]: true,
        [TransactionExploreConditionOperatorType.HasAll]: true,
        [TransactionExploreConditionOperatorType.NotHasAny]: true,
        [TransactionExploreConditionOperatorType.NotHasAll]: true
    };
    public readonly field = TransactionExploreConditionFieldType.TransactionTag;
    public readonly operator: TransactionTagConditionOperator = TransactionExploreConditionOperatorType.HasAny;
    public value: string[];

    constructor(operator: TransactionTagConditionOperator, value: string[]) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): string[] {
        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty || this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
            return [];
        }

        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        const transactionTags: Record<string, boolean> = {};

        for (const tagId of transaction.tagIds) {
            transactionTags[tagId] = true;
        }

        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty || this.value.length < 1) {
            return transaction.tagIds.length < 1;
        } else if (this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
            return transaction.tagIds.length > 0;
        } else if (this.operator === TransactionExploreConditionOperatorType.Equals || this.operator === TransactionExploreConditionOperatorType.NotEquals) {
            let hasAll = true;

            for (const tagId of this.value) {
                if (!transactionTags[tagId]) {
                    hasAll = false;
                    break;
                }
            }

            const hasSameCount = transaction.tagIds.length === this.value.length;

            if (this.operator === TransactionExploreConditionOperatorType.Equals && hasAll && hasSameCount) {
                return true;
            } else if (this.operator === TransactionExploreConditionOperatorType.NotEquals && (!hasAll || !hasSameCount)) {
                return true;
            }
        } else if (this.operator === TransactionExploreConditionOperatorType.HasAny || this.operator === TransactionExploreConditionOperatorType.NotHasAny) {
            let hasAny = false;

            for (const tagId of this.value) {
                if (transactionTags[tagId]) {
                    hasAny = true;
                    break;
                }
            }

            if (this.operator === TransactionExploreConditionOperatorType.HasAny && hasAny) {
                return true;
            } else if (this.operator === TransactionExploreConditionOperatorType.NotHasAny && !hasAny) {
                return true;
            }
        } else if (this.operator === TransactionExploreConditionOperatorType.HasAll || this.operator === TransactionExploreConditionOperatorType.NotHasAll) {
            let hasAll = true;

            for (const tagId of this.value) {
                if (!transactionTags[tagId]) {
                    hasAll = false;
                    break;
                }
            }

            if (this.operator === TransactionExploreConditionOperatorType.HasAll && hasAll) {
                return true;
            } else if (this.operator === TransactionExploreConditionOperatorType.NotHasAll && !hasAll) {
                return true;
            }
        }

        return false;
    }

    public toExpression(allCategoriesMap: Record<string, TransactionCategory>, allAccountsMap: Record<string, Account>, allTagsMap: Record<string, TransactionTag>): string {
        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty) {
            return `tags IS EMPTY`;
        } else if (this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
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

        if (this.operator === TransactionExploreConditionOperatorType.Equals) {
            return `tags FULL MATCHES (${textualTags})`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotEquals) {
            return `tags NOT FULL MATCHES (${textualTags})`;
        } else if (this.operator === TransactionExploreConditionOperatorType.HasAny) {
            return `tags HAS ANY (${textualTags})`;
        } else if (this.operator === TransactionExploreConditionOperatorType.HasAll) {
            return `tags HAS ALL (${textualTags})`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotHasAny) {
            return `tags NOT HAS ANY (${textualTags})`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotHasAll) {
            return `tags NOT HAS ALL (${textualTags})`;
        } else {
            return '';
        }
    }
}

type DescriptionConditionOperator = TransactionExploreConditionOperatorType.IsEmpty |
    TransactionExploreConditionOperatorType.IsNotEmpty |
    TransactionExploreConditionOperatorType.Equals |
    TransactionExploreConditionOperatorType.NotEquals |
    TransactionExploreConditionOperatorType.Contains |
    TransactionExploreConditionOperatorType.NotContains |
    TransactionExploreConditionOperatorType.StartsWith |
    TransactionExploreConditionOperatorType.NotStartsWith |
    TransactionExploreConditionOperatorType.EndsWith |
    TransactionExploreConditionOperatorType.NotEndsWith;

export class TransactionExploreDescriptionCondition implements TransactionExploreCondition<TransactionExploreConditionFieldType.Description, string> {
    public static readonly supportedOperators: PartialRecord<TransactionExploreConditionOperatorType, true> = {
        [TransactionExploreConditionOperatorType.IsEmpty]: true,
        [TransactionExploreConditionOperatorType.IsNotEmpty]: true,
        [TransactionExploreConditionOperatorType.Equals]: true,
        [TransactionExploreConditionOperatorType.NotEquals]: true,
        [TransactionExploreConditionOperatorType.Contains]: true,
        [TransactionExploreConditionOperatorType.NotContains]: true,
        [TransactionExploreConditionOperatorType.StartsWith]: true,
        [TransactionExploreConditionOperatorType.NotStartsWith]: true,
        [TransactionExploreConditionOperatorType.EndsWith]: true,
        [TransactionExploreConditionOperatorType.NotEndsWith]: true
    };
    public readonly field = TransactionExploreConditionFieldType.Description;
    public readonly operator: DescriptionConditionOperator = TransactionExploreConditionOperatorType.Contains;
    public value: string;

    constructor(operator: DescriptionConditionOperator, value: string) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): string {
        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty || this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
            return '';
        }

        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        const description = transaction.comment || '';

        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty) {
            return description.length === 0;
        } else if (this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
            return description.length > 0;
        } else if (this.operator === TransactionExploreConditionOperatorType.Equals) {
            return description === this.value;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotEquals) {
            return description !== this.value;
        } else if (this.operator === TransactionExploreConditionOperatorType.Contains) {
            return description.includes(this.value);
        } else if (this.operator === TransactionExploreConditionOperatorType.NotContains) {
            return !description.includes(this.value);
        } else if (this.operator === TransactionExploreConditionOperatorType.StartsWith) {
            return description.startsWith(this.value);
        } else if (this.operator === TransactionExploreConditionOperatorType.NotStartsWith) {
            return !description.startsWith(this.value);
        } else if (this.operator === TransactionExploreConditionOperatorType.EndsWith) {
            return description.endsWith(this.value);
        } else if (this.operator === TransactionExploreConditionOperatorType.NotEndsWith) {
            return !description.endsWith(this.value);
        }

        return false;
    }

    public toExpression(): string {
        if (this.operator === TransactionExploreConditionOperatorType.IsEmpty) {
            return `description IS EMPTY`;
        } else if (this.operator === TransactionExploreConditionOperatorType.IsNotEmpty) {
            return `description IS NOT EMPTY`;
        } else if (this.operator === TransactionExploreConditionOperatorType.Equals) {
            return `description = '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotEquals) {
            return `description <> '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.Contains) {
            return `description CONTAINS '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotContains) {
            return `description NOT CONTAINS '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.StartsWith) {
            return `description STARTS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotStartsWith) {
            return `description NOT STARTS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.EndsWith) {
            return `description ENDS WITH '${this.value.replace(/'/g, "''")}'`;
        } else if (this.operator === TransactionExploreConditionOperatorType.NotEndsWith) {
            return `description NOT ENDS WITH '${this.value.replace(/'/g, "''")}'`;
        }

        return '';
    }
}
