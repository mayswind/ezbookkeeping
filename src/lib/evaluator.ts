import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '../consts/transaction.ts';

import { replaceAll } from './common.ts';

import logger from './logger.ts';

type Operator = '+' | '-' | '*' | '/';
type OperatorAndParenthesis = Operator | '(' | ')';

const maxAllowedDecimalCount = 6;
const normalizeFactor: number = 1000000;
const normalizedDecimalsMaxZeroString: string = '000000';
const normalizedNumberToAmountFactor: number = 10000; // 1000000 / 100

const operatorPriority: Record<Operator, number> = {
    '+': 1,
    '-': 1,
    '*': 2,
    '/': 2,
};

function normalizeNumber(textualNumber: string): number {
    const decimalSeparatorPos = textualNumber.indexOf('.');

    if (decimalSeparatorPos < 0) {
        return parseInt(textualNumber + normalizedDecimalsMaxZeroString);
    }

    const integer = textualNumber.substring(0, decimalSeparatorPos);
    const decimals = textualNumber.substring(decimalSeparatorPos + 1);

    if (decimals.length > maxAllowedDecimalCount) {
        throw new Error('Numeric Overflow');
    }

    const paddedDecimals = (decimals + normalizedDecimalsMaxZeroString).substring(0, maxAllowedDecimalCount);
    return parseInt(integer + paddedDecimals);
}

function denormalizeNumberToAmount(num: number): number {
    return Math.trunc(num / normalizedNumberToAmountFactor);
}

function checkNumberRange(num: number): void {
    const amount = denormalizeNumberToAmount(num);

    if (amount > TRANSACTION_MAX_AMOUNT || amount < TRANSACTION_MIN_AMOUNT) {
        throw new Error('Numeric Overflow');
    }
}

function toPostfixExprTokens(expr: string): string[] | null {
    const finalTokens: string[] = [];
    const operatorStack: OperatorAndParenthesis[] = [];
    let currentNumberBuilder = '';
    let isLastTokenOperator = true;

    expr = replaceAll(expr, ' ', '');

    for (let i = 0; i < expr.length; i++) {
        const ch = expr[i] as string;

        // number
        if ('0' <= ch && ch <= '9' || ch === '.') {
            currentNumberBuilder += ch;
            continue
        } else if (ch === '-' && i + 1 < expr.length && '0' <= (expr[i + 1] as string) && (expr[i + 1] as string) <= '9' && currentNumberBuilder.length === 0 && isLastTokenOperator) {
            currentNumberBuilder += ch;
            continue
        }

        // operator or parenthesis
        if (currentNumberBuilder.length > 0) {
            finalTokens.push(currentNumberBuilder);
            currentNumberBuilder = '';
            isLastTokenOperator = false;
        }

        switch (ch) {
            case '+':
            case '-':
            case '*':
            case '/':
                if (ch === '-' && isLastTokenOperator) {
                    currentNumberBuilder += ch;
                    continue;
                }

                while (operatorStack.length > 0) {
                    const topOperator = operatorStack[operatorStack.length - 1] as OperatorAndParenthesis;

                    if (topOperator === '(') {
                        break;
                    }

                    const isCurrentOperator = topOperator === '+' || topOperator === '-' || topOperator === '*' || topOperator === '/';

                    if (isCurrentOperator && operatorPriority[topOperator] >= operatorPriority[ch]) {
                        finalTokens.push(topOperator);
                        operatorStack.pop();
                    } else {
                        break;
                    }
                }

                operatorStack.push(ch);
                isLastTokenOperator = true;
                break;
            case '(':
                operatorStack.push(ch);
                isLastTokenOperator = true;
                break;
            case ')':
                let hasLeftParenthesis = false;

                while (operatorStack.length > 0) {
                    const topOperator = operatorStack.pop() as string;

                    if (topOperator === '(') {
                        hasLeftParenthesis = true;
                        break;
                    }

                    finalTokens.push(topOperator);
                }

                if (!hasLeftParenthesis) {
                    logger.warn(`cannot parse expression "${expr}", because missing left parenthesis`);
                    return null;
                }

                isLastTokenOperator = false;
                break;
            default:
                logger.warn(`cannot parse expression "${expr}", because containing unknown token "${ch}"`);
                return null;
        }
    }

    if (currentNumberBuilder.length > 0) {
        finalTokens.push(currentNumberBuilder);
    }

    while (operatorStack.length > 0) {
        const topOperator = operatorStack.pop() as string;

        if (topOperator === '(') {
            logger.warn(`cannot parse expression "${expr}", because missing right parenthesis`);
            return null;
        }

        finalTokens.push(topOperator);
    }

    return finalTokens;
}

function evaluatePostfixExpr(tokens: string[]): number | null {
    const stack: number[] = [];

    for (let i = 0; i < tokens.length; i++) {
        const token = tokens[i] as string;

        switch (token) {
            case '+':
            case '-':
            case '*':
            case '/': // operators
                if (stack.length < 2) {
                    logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because not enough operands`);
                    return null;
                }

                // pop the top two operands
                const b = stack.pop() as number;
                const a = stack.pop() as number;

                // evaluate the operation
                let result: number;
                switch (token) {
                    case '+':
                        result = a + b;
                        break;
                    case '-':
                        result = a - b;
                        break;
                    case '*':
                        result = Math.trunc(a * b / normalizeFactor);
                        break;
                    case '/':
                        if (b === 0) {
                            logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because division by zero`);
                            return null;
                        }
                        result = Math.trunc(a * normalizeFactor / b);
                        break;
                    default:
                        return null;
                }

                checkNumberRange(result);

                // push the result back to the stack
                stack.push(result);
                break;
            default: // operands
                const normalizedNum = normalizeNumber(token);
                checkNumberRange(normalizedNum);
                stack.push(normalizedNum);
                break;
        }
    }

    if (stack.length !== 1) {
        logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because missing operator`);
        return null;
    }

    return stack[0] as number;
}
export function evaluateExpressionToAmount(expr: string): number | undefined {
    if (!expr) {
        return undefined;
    }

    const postfixExprTokens = toPostfixExprTokens(expr);

    if (!postfixExprTokens) {
        return undefined;
    }

    const result = evaluatePostfixExpr(postfixExprTokens);

    if (result === null) {
        return undefined;
    }

    return denormalizeNumberToAmount(result);
}
