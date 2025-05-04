import { replaceAll } from './common.ts';

import logger from './logger.ts';

const operatorPriority: Record<string, number> = {
    '+': 1,
    '-': 1,
    '*': 2,
    '/': 2,
};

function toPostfixExprTokens(expr: string): string[] | null {
    const finalTokens: string[] = [];
    const operatorStack: string[] = [];
    let currentNumberBuilder = '';
    let isLastTokenOperator = true;

    expr = replaceAll(expr, ' ', '');

    for (let i = 0; i < expr.length; i++) {
        const ch = expr[i];

        // number
        if ('0' <= ch && ch <= '9' || ch === '.') {
            currentNumberBuilder += ch;
            continue
        } else if (ch === '-' && i + 1 < expr.length && '0' <= expr[i + 1] && expr[i + 1] <= '9' && currentNumberBuilder.length === 0 && isLastTokenOperator) {
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
                    const topOperator = operatorStack[operatorStack.length - 1];

                    if (topOperator === '(') {
                        break;
                    }

                    if (operatorPriority[topOperator] >= operatorPriority[ch]) {
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
        const token = tokens[i];

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
                        result = a * b;
                        break;
                    case '/':
                        if (b === 0) {
                            logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because division by zero`);
                            return null;
                        }
                        result = a / b;
                        break;
                    default:
                        return null;
                }

                // push the result back to the stack
                stack.push(result);
                break;
            default: // operands
                const num = parseFloat(token);

                if (isNaN(num)) {
                    logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because containing invalid number`);
                    return null;
                }

                stack.push(num);
                break;
        }
    }

    if (stack.length !== 1) {
        logger.warn(`cannot evaluate expression "${tokens.join(' ')}", because missing operator`);
        return null;
    }

    return stack[0];
}
export function evaluateExpression(expr: string): number | undefined {
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

    return result;
}
