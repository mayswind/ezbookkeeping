import type { BigDecimal } from '@/core/numeral.ts';
import { AMOUNT_FACTOR } from '@/consts/numeral.ts';
import { TRANSACTION_MIN_AMOUNT, TRANSACTION_MAX_AMOUNT } from '../consts/transaction.ts';

import { replaceAll } from './common.ts';
import { parseBigDecimal } from './numeral.ts';

import logger from './logger.ts';

type Operator = '+' | '-' | '*' | '/';
type OperatorAndParenthesis = Operator | '(' | ')';

const minAllowedAmount: BigDecimal = parseBigDecimal(TRANSACTION_MIN_AMOUNT).divide(AMOUNT_FACTOR);
const maxAllowedAmount: BigDecimal = parseBigDecimal(TRANSACTION_MAX_AMOUNT).divide(AMOUNT_FACTOR);

const operatorPriority: Record<Operator, number> = {
    '+': 1,
    '-': 1,
    '*': 2,
    '/': 2,
};

function parseNumber(textualNumber: string): BigDecimal {
    const result: BigDecimal = parseBigDecimal(textualNumber);

    if (!result.between(minAllowedAmount, maxAllowedAmount)) {
        throw new Error('Numeric Overflow');
    }

    return result;
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

function evaluatePostfixExpr(tokens: string[]): BigDecimal | null {
    const stack: BigDecimal[] = [];

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
                const b = stack.pop() as BigDecimal;
                const a = stack.pop() as BigDecimal;

                // evaluate the operation
                let result: BigDecimal;
                switch (token) {
                    case '+':
                        result = a.add(b);
                        break;
                    case '-':
                        result = a.subtract(b);
                        break;
                    case '*':
                        result = a.multiply(b);
                        break;
                    case '/':
                        if (b.isZero()) {
                            logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because division by zero`);
                            return null;
                        }
                        result = a.divide(b);
                        break;
                    default:
                        return null;
                }

                // push the result back to the stack
                stack.push(result);
                break;
            default: // operands
                const num = parseNumber(token);
                stack.push(num);
                break;
        }
    }

    if (stack.length !== 1) {
        logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because missing operator`);
        return null;
    }

    return stack[0] as BigDecimal;
}
export function evaluateExpressionToAmount(expr: string): BigDecimal | undefined {
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

    if (!result.between(minAllowedAmount, maxAllowedAmount)) {
        throw new Error('Numeric Overflow');
    }

    return result.multiply(AMOUNT_FACTOR).truncate();
}
