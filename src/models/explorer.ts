import { type PartialRecord, itemAndIndex, keysIfValueEquals } from '@/core/base.ts';
import { TimezoneTypeForStatistics } from '@/core/timezone.ts';
import { AccountType } from '@/core/account.ts';
import { TransactionType } from '@/core/transaction.ts';
import { ChartSortingType } from '@/core/statistics.ts';
import {
    type TransactionExplorerSubConditionStartRelation,
    TransactionExplorerConditionRelation,
    TransactionExplorerSubConditionStartRelationPlaceholder,
    TransactionExplorerConditionRelationPriority,
    TransactionExplorerConditionFieldType,
    TransactionExplorerConditionField,
    TransactionExplorerConditionOperatorType,
    TransactionExplorerConditionOperator,
    TransactionExplorerChartTypeValue,
    TransactionExplorerChartType,
    TransactionExplorerDataDimensionType,
    TransactionExplorerDataDimension,
    TransactionExplorerValueMetricType,
    TransactionExplorerValueMetric
} from '@/core/explorer.ts';

import { Account } from '@/models/account.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';
import { type TransactionInsightDataItem } from '@/models/transaction.ts';

export class InsightsExplorerBasicInfo implements InsightsExplorerInfoResponse {
    public id: string;
    public name: string;
    public displayOrder: number;
    public hidden: boolean;
    public data: Record<string, string | number | object[]> = {};

    private constructor(id: string, name: string, displayOrder: number, hidden: boolean) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
        this.hidden = hidden;
    }

    public static of(explorerResponse: InsightsExplorerInfoResponse): InsightsExplorerBasicInfo {
        return new InsightsExplorerBasicInfo(
            explorerResponse.id,
            explorerResponse.name,
            explorerResponse.displayOrder,
            explorerResponse.hidden
        );
    }

    public static ofMulti(explorerResponses: InsightsExplorerInfoResponse[]): InsightsExplorerBasicInfo[] {
        const explorers: InsightsExplorerBasicInfo[] = [];

        for (const explorerResponse of explorerResponses) {
            explorers.push(InsightsExplorerBasicInfo.of(explorerResponse));
        }

        return explorers;
    }
}

export class InsightsExplorer implements InsightsExplorerInfoResponse {
    public id: string;
    public name: string;
    public displayOrder: number;
    public hidden: boolean;
    public queries: TransactionExplorerQuery[];
    public timezoneUsedForDateRange: number;
    public datatableQuerySource: string;
    public countPerPage: number;
    public chartType: TransactionExplorerChartTypeValue;
    public categoryDimension: TransactionExplorerDataDimensionType;
    public seriesDimension: TransactionExplorerDataDimensionType;
    public valueMetric: TransactionExplorerValueMetricType;
    public chartSortingType: number;

    public static readonly Default: InsightsExplorer = new InsightsExplorer(
        '',
        '',
        0,
        false,
        [],
        TimezoneTypeForStatistics.Default.type,
        '',
        15,
        TransactionExplorerChartType.Default.value,
        TransactionExplorerDataDimension.CategoryDimensionDefault.value,
        TransactionExplorerDataDimension.SeriesDimensionDefault.value,
        TransactionExplorerValueMetric.Default.value,
        ChartSortingType.Default.type
    );

    private constructor(id: string, name: string, displayOrder: number, hidden: boolean, queries: TransactionExplorerQuery[], timezoneUsedForDateRange: number, datatableQuerySource: string, countPerPage: number, chartType: TransactionExplorerChartTypeValue, categoryDimension: TransactionExplorerDataDimensionType, seriesDimension: TransactionExplorerDataDimensionType, valueMetric: TransactionExplorerValueMetricType, chartSortingType: number) {
        this.id = id;
        this.name = name;
        this.displayOrder = displayOrder;
        this.hidden = hidden;
        this.queries = queries;
        this.timezoneUsedForDateRange = timezoneUsedForDateRange;
        this.datatableQuerySource = datatableQuerySource;
        this.countPerPage = countPerPage;
        this.chartType = chartType;
        this.categoryDimension = categoryDimension;
        this.seriesDimension = seriesDimension;
        this.valueMetric = valueMetric;
        this.chartSortingType = chartSortingType;
    }

    public get data(): Record<string, string | number | object[]> {
        return {
            queries: this.queries.map(q => q.toJsonObject()),
            timezoneUsedForDateRange: this.timezoneUsedForDateRange,
            datatableQuerySource: this.datatableQuerySource,
            countPerPage: this.countPerPage,
            chartType: this.chartType,
            categoryDimension: this.categoryDimension,
            seriesDimension: this.seriesDimension,
            valueMetric: this.valueMetric,
            chartSortingType: this.chartSortingType
        };
    }

    public toCreateRequest(clientSessionId: string): InsightsExplorerCreateRequest {
        return {
            name: this.name,
            data: this.data,
            clientSessionId: clientSessionId
        };
    }

    public toModifyRequest(): InsightsExplorerModifyRequest {
        return {
            id: this.id,
            name: this.name,
            data: this.data,
            hidden: this.hidden
        };
    }

    public static of(explorerResponse: InsightsExplorerInfoResponse): InsightsExplorer {
        const data = explorerResponse.data;
        const queries: TransactionExplorerQuery[] = [];
        let timezoneUsedForDateRange = InsightsExplorer.Default.timezoneUsedForDateRange;
        let datatableQuerySource = InsightsExplorer.Default.datatableQuerySource;
        let countPerPage = InsightsExplorer.Default.countPerPage;
        let chartType = InsightsExplorer.Default.chartType;
        let categoryDimension = InsightsExplorer.Default.categoryDimension;
        let seriesDimension = InsightsExplorer.Default.seriesDimension;
        let valueMetric = InsightsExplorer.Default.valueMetric;
        let chartSortingType = InsightsExplorer.Default.chartSortingType;
        let hasDatatableQuerySource = false;

        if (data) {
            if (typeof data['timezoneUsedForDateRange'] === 'number') {
                timezoneUsedForDateRange = data['timezoneUsedForDateRange'] as number;
            }

            if (typeof data['datatableQuerySource'] === 'string') {
                datatableQuerySource = data['datatableQuerySource'] as string;
            }

            if (typeof data['countPerPage'] === 'number') {
                countPerPage = data['countPerPage'] as number;
            }

            if (typeof data['chartType'] === 'string') {
                chartType = data['chartType'] as TransactionExplorerChartTypeValue;
            }

            if (typeof data['categoryDimension'] === 'string') {
                categoryDimension = data['categoryDimension'] as TransactionExplorerDataDimensionType;
            }

            if (typeof data['seriesDimension'] === 'string') {
                seriesDimension = data['seriesDimension'] as TransactionExplorerDataDimensionType;
            }

            if (typeof data['valueMetric'] === 'string') {
                valueMetric = data['valueMetric'] as TransactionExplorerValueMetricType;
            }

            if (typeof data['chartSortingType'] === 'number') {
                chartSortingType = data['chartSortingType'] as number;
            }

            if (Array.isArray(data['queries'])) {
                const queryItems = data['queries'] as object[];

                for (const queryItem of queryItems) {
                    const query = TransactionExplorerQuery.parse(queryItem);

                    if (query) {
                        queries.push(query);
                    }

                    if (query && query.id === datatableQuerySource) {
                        hasDatatableQuerySource = true;
                    }
                }
            }

            if (!hasDatatableQuerySource) {
                datatableQuerySource = '';
            }
        }

        return new InsightsExplorer(
            explorerResponse.id,
            explorerResponse.name,
            explorerResponse.displayOrder,
            explorerResponse.hidden,
            queries,
            timezoneUsedForDateRange,
            datatableQuerySource,
            countPerPage,
            chartType,
            categoryDimension,
            seriesDimension,
            valueMetric,
            chartSortingType
        );
    }

    public static createNewExplorer(newQueryId: string): InsightsExplorer {
        return new InsightsExplorer(
            '',
            '',
            0,
            false,
            [TransactionExplorerQuery.create(newQueryId)],
            InsightsExplorer.Default.timezoneUsedForDateRange,
            InsightsExplorer.Default.datatableQuerySource,
            InsightsExplorer.Default.countPerPage,
            InsightsExplorer.Default.chartType,
            InsightsExplorer.Default.categoryDimension,
            InsightsExplorer.Default.seriesDimension,
            InsightsExplorer.Default.valueMetric,
            InsightsExplorer.Default.chartSortingType
        );
    }
}

export interface InsightsExplorerCreateRequest {
    readonly name: string;
    readonly data: Record<string, string | number | object[]>;
    readonly clientSessionId?: string;
}

export interface InsightsExplorerModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly data: Record<string, string | number | object[]>;
    readonly hidden: boolean;
    readonly clientSessionId?: string;
}

export interface InsightsExplorerHideRequest {
    readonly id: string;
    readonly hidden: boolean;
}

export interface InsightsExplorerMoveRequest {
    readonly newDisplayOrders: InsightsExplorerNewDisplayOrderRequest[];
}

export interface InsightsExplorerNewDisplayOrderRequest {
    readonly id: string;
    readonly displayOrder: number;
}

export interface InsightsExplorerDeleteRequest {
    readonly id: string;
}

export interface InsightsExplorerInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly displayOrder: number;
    readonly hidden: boolean;
    readonly data: Record<string, string | number | object[]>;
}

interface ExpressionNode {
    textualExpression: string;
    operator?: TransactionExplorerConditionRelation;
}

export class TransactionExplorerQuery {
    public id: string;
    public name: string;
    public conditions: TransactionExplorerConditionWithRelation[];

    private constructor(id: string, name: string, conditions: TransactionExplorerConditionWithRelation[]) {
        this.id = id;
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
            case TransactionExplorerConditionField.GeoLocation:
                condition = new TransactionExplorerGeoLocationCondition(TransactionExplorerConditionOperatorType.IsNotEmpty, [TransactionExplorerGeoLocationCondition.MIN_LATITUDE, TransactionExplorerGeoLocationCondition.MAX_LATITUDE, TransactionExplorerGeoLocationCondition.MIN_LONGITUDE, TransactionExplorerGeoLocationCondition.MAX_LONGITUDE]);
                break;
            case TransactionExplorerConditionField.TransactionTag:
                condition = new TransactionExplorerTransactionTagCondition(TransactionExplorerConditionOperatorType.HasAny, []);
                break;
            case TransactionExplorerConditionField.Pictures:
                condition = new TransactionExplorerPicturesCondition(TransactionExplorerConditionOperatorType.IsNotEmpty, []);
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

    public addSubConditionEnd(): TransactionExplorerConditionWithRelation {
        return new TransactionExplorerConditionWithRelation(new TransactionExplorerUndefinedCondition(), TransactionExplorerConditionRelation.SubEnd);
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

        const operatorStack: (TransactionExplorerConditionRelation | TransactionExplorerSubConditionStartRelation)[] = [];
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

            if (item.relation === TransactionExplorerConditionRelation.SubEnd) {
                while (operatorStack.length > 0) {
                    const topOperator = operatorStack.pop();

                    if (topOperator === TransactionExplorerSubConditionStartRelationPlaceholder) {
                        break;
                    }

                    const isAndOrOperator = topOperator === TransactionExplorerConditionRelation.And || topOperator === TransactionExplorerConditionRelation.Or;

                    if (isAndOrOperator) {
                        finalTokens.push(topOperator);
                    } else {
                        throw new Error('invalid operator in stack');
                    }
                }
            } else { // And, Or, AndSub, OrSub
                let currentOperator: TransactionExplorerConditionRelation.And | TransactionExplorerConditionRelation.Or;
                let startNewSubCondition = false;

                switch (item.relation) {
                    case TransactionExplorerConditionRelation.AndSub:
                        currentOperator = TransactionExplorerConditionRelation.And;
                        startNewSubCondition = true;
                        break;
                    case TransactionExplorerConditionRelation.OrSub:
                        currentOperator = TransactionExplorerConditionRelation.Or;
                        startNewSubCondition = true;
                        break;
                    case TransactionExplorerConditionRelation.And:
                        currentOperator = item.relation;
                        break;
                    case TransactionExplorerConditionRelation.Or:
                        currentOperator = item.relation;
                        break;
                    default:
                        throw new Error('invalid operator in stack');
                }

                while (operatorStack.length > 0) {
                    const topOperator = operatorStack[operatorStack.length - 1];

                    if (topOperator === TransactionExplorerSubConditionStartRelationPlaceholder) {
                        break;
                    }

                    const isAndOrOperator = topOperator === TransactionExplorerConditionRelation.And || topOperator === TransactionExplorerConditionRelation.Or;

                    if (isAndOrOperator && TransactionExplorerConditionRelationPriority[topOperator] >= TransactionExplorerConditionRelationPriority[currentOperator]) {
                        finalTokens.push(topOperator);
                        operatorStack.pop();
                    } else {
                        break;
                    }
                }

                operatorStack.push(currentOperator);

                if (startNewSubCondition) {
                    operatorStack.push(TransactionExplorerSubConditionStartRelationPlaceholder);
                }

                finalTokens.push(item.condition);
            }
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

    public getConditionNestingDepths(): number[] {
        const depths: number[] = [];
        let depth = 0;

        for (const item of this.conditions) {
            if (item.relation === TransactionExplorerConditionRelation.SubEnd) {
                depth--;
                depths.push(depth);
            } else if (item.relation === TransactionExplorerConditionRelation.AndSub || item.relation === TransactionExplorerConditionRelation.OrSub) {
                depth++;
                depths.push(depth);
            } else {
                depths.push(depth);
            }
        }

        return depths;
    }

    public clone(newId: string): TransactionExplorerQuery {
        const clonedConditions: TransactionExplorerConditionWithRelation[] = [];

        for (const condition of this.conditions) {
            const clonedCondition = TransactionExplorerConditionWithRelation.parse(condition.toJsonObject());

            if (!clonedCondition) {
                continue;
            }

            clonedConditions.push(clonedCondition);
        }

        return new TransactionExplorerQuery(newId, this.name, clonedConditions);
    }

    public toJsonObject(): object {
        return {
            id: this.id,
            name: this.name,
            conditions: this.conditions.map(condition => condition.toJsonObject())
        };
    }

    public static create(id: string): TransactionExplorerQuery {
        return new TransactionExplorerQuery(id, '', []);
    }

    public static parse(obj: object): TransactionExplorerQuery | null {
        if (!('id' in obj) || !('name' in obj) || !('conditions' in obj)) {
            return null;
        }

        const idFieldValue = obj['id'] as unknown;
        const nameFieldValue = obj['name'] as unknown;
        const conditionsFieldValue = obj['conditions'] as unknown;

        if (typeof idFieldValue !== 'string' || typeof nameFieldValue !== 'string' || !Array.isArray(conditionsFieldValue)) {
            return null;
        }

        const id: string = idFieldValue;
        const name: string = nameFieldValue;
        const conditions: TransactionExplorerConditionWithRelation[] = [];
        let conditionDepth = 0;

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

            if (condition.relation === TransactionExplorerConditionRelation.AndSub ||
                condition.relation === TransactionExplorerConditionRelation.OrSub) {
                conditionDepth++;
            } else if (condition.relation === TransactionExplorerConditionRelation.SubEnd) {
                conditionDepth--;
            }

            conditions.push(condition);
        }

        if (conditionDepth !== 0) {
            return null; // unbalanced parentheses
        }

        return new TransactionExplorerQuery(id, name, conditions);
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
            case TransactionExplorerConditionField.GeoLocation.value:
                operatorTypes = TransactionExplorerGeoLocationCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.TransactionTag.value:
                operatorTypes = TransactionExplorerTransactionTagCondition.supportedOperators;
                break;
            case TransactionExplorerConditionField.Pictures.value:
                operatorTypes = TransactionExplorerPicturesCondition.supportedOperators;
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
        if (this.relation === TransactionExplorerConditionRelation.SubEnd) {
            return {
                relation: this.relation
            };
        }

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
        if (typeof data !== 'object' || !data || !('relation' in data)) {
            return null;
        }

        const relation = data['relation'];
        let condition: TransactionExplorerCondition | null = null;

        if (relation === TransactionExplorerConditionRelation.First ||
            relation === TransactionExplorerConditionRelation.And || relation === TransactionExplorerConditionRelation.Or ||
            relation === TransactionExplorerConditionRelation.AndSub || relation === TransactionExplorerConditionRelation.OrSub) {
            if (!('condition' in data)) {
                return null;
            }

            const conditionObject = data['condition'];

            if (typeof conditionObject !== 'object' || !conditionObject || !('field' in conditionObject) || !('operator' in conditionObject) || !('value' in conditionObject) || typeof relation !== 'string') {
                return null;
            }

            const conditionField = conditionObject['field'];
            const conditionOperator = conditionObject['operator'] as TransactionExplorerConditionOperatorType;
            const conditionValue = conditionObject['value'];

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
                case TransactionExplorerConditionField.GeoLocation.value:
                    if (TransactionExplorerGeoLocationCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue)) {
                        condition = new TransactionExplorerGeoLocationCondition(conditionOperator as GeoLocationConditionOperator, conditionValue as [number, number, number, number]);
                    }
                    break;
                case TransactionExplorerConditionField.TransactionTag.value:
                    if (TransactionExplorerTransactionTagCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue)) {
                        condition = new TransactionExplorerTransactionTagCondition(conditionOperator as TransactionTagConditionOperator, conditionValue as string[]);
                    }
                    break;
                case TransactionExplorerConditionField.Pictures.value:
                    if (TransactionExplorerPicturesCondition.supportedOperators[conditionOperator] && Array.isArray(conditionValue)) {
                        condition = new TransactionExplorerPicturesCondition(conditionOperator as PicturesConditionOperator, conditionValue as string[]);
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
        } else if (relation === TransactionExplorerConditionRelation.SubEnd) {
            condition = new TransactionExplorerUndefinedCondition();
        } else {
            return null;
        }

        if (condition === null) {
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

export class TransactionExplorerUndefinedCondition implements TransactionExplorerCondition {
    public readonly field = TransactionExplorerConditionFieldType.Undefined;
    public readonly operator = TransactionExplorerConditionOperatorType.Equals;
    public value = '';

    public getValueForStore(): string {
        return this.value;
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        return !!transaction;
    }

    public toExpression(): string {
        return '';
    }
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

type GeoLocationConditionOperator = TransactionExplorerConditionOperatorType.IsEmpty |
    TransactionExplorerConditionOperatorType.IsNotEmpty |
    TransactionExplorerConditionOperatorType.LatitudeBetween |
    TransactionExplorerConditionOperatorType.LatitudeNotBetween |
    TransactionExplorerConditionOperatorType.LongitudeBetween |
    TransactionExplorerConditionOperatorType.LongitudeNotBetween;

export class TransactionExplorerGeoLocationCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.GeoLocation, [number, number, number, number]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.IsEmpty]: true,
        [TransactionExplorerConditionOperatorType.IsNotEmpty]: true,
        [TransactionExplorerConditionOperatorType.LatitudeBetween]: true,
        [TransactionExplorerConditionOperatorType.LatitudeNotBetween]: true,
        [TransactionExplorerConditionOperatorType.LongitudeBetween]: true,
        [TransactionExplorerConditionOperatorType.LongitudeNotBetween]: true
    };
    public static readonly MIN_LATITUDE: number = -90.0;
    public static readonly MAX_LATITUDE: number = 90.0;
    public static readonly MIN_LONGITUDE: number = -180.0;
    public static readonly MAX_LONGITUDE: number = 180.0;

    public readonly field = TransactionExplorerConditionFieldType.GeoLocation;
    public readonly operator: GeoLocationConditionOperator = TransactionExplorerConditionOperatorType.IsNotEmpty;
    public value: [number, number, number, number];

    constructor(operator: GeoLocationConditionOperator, value: [number, number, number, number]) {
        this.operator = operator;
        this.value = [
            value[0] ?? TransactionExplorerGeoLocationCondition.MIN_LATITUDE,
            value[1] ?? TransactionExplorerGeoLocationCondition.MAX_LATITUDE,
            value[2] ?? TransactionExplorerGeoLocationCondition.MIN_LONGITUDE,
            value[3] ?? TransactionExplorerGeoLocationCondition.MAX_LONGITUDE
        ];
    }

    public getValueForStore(): [number, number, number, number] {
        return [
            this.value[0] ?? TransactionExplorerGeoLocationCondition.MIN_LATITUDE,
            this.value[1] ?? TransactionExplorerGeoLocationCondition.MAX_LATITUDE,
            this.value[2] ?? TransactionExplorerGeoLocationCondition.MIN_LONGITUDE,
            this.value[3] ?? TransactionExplorerGeoLocationCondition.MAX_LONGITUDE
        ];
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return !transaction.geoLocation;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return !!transaction.geoLocation;
        } else if (this.operator === TransactionExplorerConditionOperatorType.LatitudeBetween) {
            if (!transaction.geoLocation) {
                return false;
            }

            const latitude = transaction.geoLocation.latitude;

            if (typeof(this.value[0]) === 'number' && latitude < this.value[0]) {
                return false;
            }

            if (typeof(this.value[1]) === 'number' && latitude > this.value[1]) {
                return false;
            }

            return true;
        } else if (this.operator === TransactionExplorerConditionOperatorType.LatitudeNotBetween) {
            if (!transaction.geoLocation) {
                return false;
            }

            const latitude = transaction.geoLocation.latitude;

            if (typeof(this.value[0]) === 'number' && typeof(this.value[1]) === 'number' && latitude >= this.value[0] && latitude <= this.value[1]) {
                return false;
            }

            return true;
        } else if (this.operator === TransactionExplorerConditionOperatorType.LongitudeBetween) {
            if (!transaction.geoLocation) {
                return false;
            }

            const longitude = transaction.geoLocation.longitude;

            if (typeof(this.value[2]) === 'number' && longitude < this.value[2]) {
                return false;
            }

            if (typeof(this.value[3]) === 'number' && longitude > this.value[3]) {
                return false;
            }

            return true;
        } else if (this.operator === TransactionExplorerConditionOperatorType.LongitudeNotBetween) {
            if (!transaction.geoLocation) {
                return false;
            }

            const longitude = transaction.geoLocation.longitude;

            if (typeof(this.value[2]) === 'number' && typeof(this.value[3]) === 'number' && longitude >= this.value[2] && longitude <= this.value[3]) {
                return false;
            }

            return true;
        }

        return false;
    }

    public toExpression(): string {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return `geo_location IS EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return `geo_location IS NOT EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.LatitudeBetween || this.operator === TransactionExplorerConditionOperatorType.LatitudeNotBetween) {
            const minLatitude: number = this.value[0] ?? TransactionExplorerGeoLocationCondition.MIN_LATITUDE;
            const maxLatitude: number = this.value[1] ?? TransactionExplorerGeoLocationCondition.MAX_LATITUDE;

            if (this.operator === TransactionExplorerConditionOperatorType.LatitudeBetween) {
                return `(geo_location.latitude >= ${minLatitude} AND geo_location.latitude <= ${maxLatitude})`;
            } else if (this.operator === TransactionExplorerConditionOperatorType.LatitudeNotBetween) {
                return `(geo_location.latitude < ${minLatitude} OR geo_location.latitude > ${maxLatitude})`;
            }
        } else if (this.operator === TransactionExplorerConditionOperatorType.LongitudeBetween || this.operator === TransactionExplorerConditionOperatorType.LongitudeNotBetween) {
            const minLongitude: number = this.value[2] ?? TransactionExplorerGeoLocationCondition.MIN_LONGITUDE;
            const maxLongitude: number = this.value[3] ?? TransactionExplorerGeoLocationCondition.MAX_LONGITUDE;

            if (this.operator === TransactionExplorerConditionOperatorType.LongitudeBetween) {
                return `(geo_location.longitude >= ${minLongitude} AND geo_location.longitude <= ${maxLongitude})`;
            } else if (this.operator === TransactionExplorerConditionOperatorType.LongitudeNotBetween) {
                return `(geo_location.longitude < ${minLongitude} OR geo_location.longitude > ${maxLongitude})`;
            }
        }

        return '';
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

type PicturesConditionOperator = TransactionExplorerConditionOperatorType.IsEmpty |
    TransactionExplorerConditionOperatorType.IsNotEmpty;

export class TransactionExplorerPicturesCondition implements TransactionExplorerCondition<TransactionExplorerConditionFieldType.Pictures, string[]> {
    public static readonly supportedOperators: PartialRecord<TransactionExplorerConditionOperatorType, true> = {
        [TransactionExplorerConditionOperatorType.IsEmpty]: true,
        [TransactionExplorerConditionOperatorType.IsNotEmpty]: true
    };

    public readonly field = TransactionExplorerConditionFieldType.Pictures;
    public readonly operator: PicturesConditionOperator = TransactionExplorerConditionOperatorType.IsNotEmpty;
    public value: string[];

    constructor(operator: PicturesConditionOperator, value: string[]) {
        this.operator = operator;
        this.value = value;
    }

    public getValueForStore(): string[] {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty || this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return [];
        }

        return [];
    }

    public match(transaction: TransactionInsightDataItem): boolean {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return !transaction.pictures || transaction.pictures.length < 1;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return !!transaction.pictures && transaction.pictures.length > 0;
        }

        return false;
    }

    public toExpression(): string {
        if (this.operator === TransactionExplorerConditionOperatorType.IsEmpty) {
            return `pictures IS EMPTY`;
        } else if (this.operator === TransactionExplorerConditionOperatorType.IsNotEmpty) {
            return `pictures IS NOT EMPTY`;
        }

        return '';
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
