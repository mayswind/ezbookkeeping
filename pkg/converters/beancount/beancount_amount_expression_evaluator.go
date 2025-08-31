package beancount

import (
	"math/big"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

const maxAllowedDecimalCount = 6
const normalizeFactor = int64(1000000)
const normalizedDecimalsMaxZeroString = "000000"
const normalizedNumberToAmountFactor = int64(10000) // 1000000 / 100

var operatorPriority = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func normalizeNumber(textualNumber string) (*big.Int, error) {
	decimalSeparatorPos := strings.Index(textualNumber, ".")

	if decimalSeparatorPos < 0 {
		result := big.NewInt(0)
		_, ok := result.SetString(textualNumber+normalizedDecimalsMaxZeroString, 10)

		if !ok {
			return nil, errs.ErrAmountInvalid
		}

		return result, nil
	}

	integer := utils.SubString(textualNumber, 0, decimalSeparatorPos)
	decimals := utils.SubString(textualNumber, decimalSeparatorPos+1, len(textualNumber))

	if len(decimals) > maxAllowedDecimalCount {
		return nil, errs.ErrAmountInvalid
	}

	paddedDecimals := utils.SubString(decimals+normalizedDecimalsMaxZeroString, 0, maxAllowedDecimalCount)
	result := big.NewInt(0)
	_, ok := result.SetString(integer+paddedDecimals, 10)

	if !ok {
		return nil, errs.ErrAmountInvalid
	}

	return result, nil
}

func denormalizeNumberToTextualAmount(num *big.Int) string {
	result := big.NewInt(0).Add(num, big.NewInt(0)) // make a copy of num
	result = result.Div(result, big.NewInt(normalizedNumberToAmountFactor))
	return utils.FormatAmount(result.Int64())
}

func toPostfixExprTokens(ctx core.Context, expr string) ([]string, error) {
	finalTokens := make([]string, 0)
	operatorStack := make([]rune, 0)
	currentNumberBuilder := strings.Builder{}
	isLastTokenOperator := true

	expr = strings.ReplaceAll(expr, " ", "")

	for i := 0; i < len(expr); i++ {
		ch := rune(expr[i])

		// number
		if '0' <= ch && ch <= '9' || ch == '.' {
			currentNumberBuilder.WriteRune(ch)
			continue
		} else if ch == '-' && i+1 < len(expr) && '0' <= expr[i+1] && expr[i+1] <= '9' && currentNumberBuilder.Len() == 0 && isLastTokenOperator {
			currentNumberBuilder.WriteRune(ch)
			continue
		}

		// operator or parenthesis
		if currentNumberBuilder.Len() > 0 {
			finalTokens = append(finalTokens, currentNumberBuilder.String())
			currentNumberBuilder.Reset()
			isLastTokenOperator = false
		}

		switch ch {
		case '+', '-', '*', '/':
			if ch == '-' && isLastTokenOperator {
				currentNumberBuilder.WriteRune(ch)
				continue
			}

			for len(operatorStack) > 0 {
				topOperator := operatorStack[len(operatorStack)-1]

				if topOperator == '(' {
					break
				}

				if operatorPriority[topOperator] >= operatorPriority[ch] {
					finalTokens = append(finalTokens, string(topOperator))
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}

			operatorStack = append(operatorStack, ch)
			isLastTokenOperator = true
		case '(':
			operatorStack = append(operatorStack, ch)
			isLastTokenOperator = true
		case ')':
			hasLeftParenthesis := false

			for len(operatorStack) > 0 {
				topOperator := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]

				if topOperator == '(' {
					hasLeftParenthesis = true
					break
				}

				finalTokens = append(finalTokens, string(topOperator))
			}

			if !hasLeftParenthesis {
				log.Warnf(ctx, "[beancount_amount_expression_evaluator.toPostfixExprTokens] cannot parse expression \"%s\", because missing left parenthesis", expr)
				return nil, errs.ErrInvalidAmountExpression
			}

			isLastTokenOperator = false
		default:
			log.Warnf(ctx, "[beancount_amount_expression_evaluator.toPostfixExprTokens] cannot parse expression \"%s\", because containing unknown token \"%c\"", expr, ch)
			return nil, errs.ErrInvalidAmountExpression
		}
	}

	if currentNumberBuilder.Len() > 0 {
		finalTokens = append(finalTokens, currentNumberBuilder.String())
	}

	for len(operatorStack) > 0 {
		topOperator := operatorStack[len(operatorStack)-1]
		operatorStack = operatorStack[:len(operatorStack)-1]

		if topOperator == '(' {
			log.Warnf(ctx, "[beancount_amount_expression_evaluator.toPostfixExprTokens] cannot parse expression \"%s\", because missing right parenthesis", expr)
			return nil, errs.ErrInvalidAmountExpression
		}

		finalTokens = append(finalTokens, string(topOperator))
	}

	return finalTokens, nil
}

func evaluatePostfixExpr(ctx core.Context, tokens []string) (*big.Int, error) {
	stack := make([]*big.Int, 0)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch token {
		case "+", "-", "*", "/": // operators
			if len(stack) < 2 {
				log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because not enough operands", strings.Join(tokens, " "))
				return nil, errs.ErrInvalidAmountExpression
			}

			// pop the top two operands
			b := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			// evaluate the operation
			result := big.NewInt(0)
			switch token {
			case "+":
				result.Add(a, b)
			case "-":
				result.Sub(a, b)
			case "*":
				result.Mul(a, b)
				result.Div(result, big.NewInt(normalizeFactor))
			case "/":
				if b.Int64() == 0 {
					log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because division by zero", strings.Join(tokens, " "))
					return nil, errs.ErrInvalidAmountExpression
				}
				result.Mul(a, big.NewInt(normalizeFactor))
				result.Div(result, b)
			}

			// push the result back to the stack
			stack = append(stack, result)
		default: // operands
			normalizedNum, err := normalizeNumber(token)

			if err != nil {
				log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because containing invalid number", strings.Join(tokens, " "))
				return nil, errs.ErrInvalidAmountExpression
			}

			stack = append(stack, normalizedNum)
		}
	}

	if len(stack) != 1 {
		log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because missing operator", strings.Join(tokens, " "))
		return nil, errs.ErrInvalidAmountExpression
	}

	return stack[0], nil
}

func evaluateBeancountAmountExpression(ctx core.Context, expr string) (string, error) {
	if expr == "" {
		return "", nil
	}

	postfixExprTokens, err := toPostfixExprTokens(ctx, expr)

	if err != nil {
		return "", err
	}

	result, err := evaluatePostfixExpr(ctx, postfixExprTokens)

	if err != nil {
		return "", err
	}

	return denormalizeNumberToTextualAmount(result), nil
}
