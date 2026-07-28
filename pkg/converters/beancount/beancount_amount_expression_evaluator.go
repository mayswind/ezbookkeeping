package beancount

import (
	"errors"
	"math/big"
	"strings"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
	"github.com/mayswind/ezbookkeeping/pkg/log"
	"github.com/mayswind/ezbookkeeping/pkg/models"
	"github.com/mayswind/ezbookkeeping/pkg/utils"
)

var (
	minAllowedAmount = big.NewRat(models.MinimumTransactionAmount, 100)
	maxAllowedAmount = big.NewRat(models.MaximumTransactionAmount, 100)

	operatorPriority = map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}
)

func parseNumber(textualNumber string) (*big.Rat, error) {
	result := new(big.Rat)
	_, ok := result.SetString(textualNumber)

	if !ok {
		return nil, errs.ErrAmountInvalid
	}

	if result.Cmp(minAllowedAmount) < 0 || result.Cmp(maxAllowedAmount) > 0 {
		return nil, errs.ErrNumericOverflow
	}

	return result, nil
}

func numberToTextualAmount(num *big.Rat) (string, error) {
	if num.Cmp(minAllowedAmount) < 0 || num.Cmp(maxAllowedAmount) > 0 {
		return "", errs.ErrNumericOverflow
	}

	amount := new(big.Rat).Mul(num, big.NewRat(100, 1))
	truncatedAmount := new(big.Int).Quo(amount.Num(), amount.Denom())
	return utils.FormatAmount(truncatedAmount.Int64()), nil
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

func evaluatePostfixExpr(ctx core.Context, tokens []string) (*big.Rat, error) {
	stack := make([]*big.Rat, 0)

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
			result := new(big.Rat)
			switch token {
			case "+":
				result.Add(a, b)
			case "-":
				result.Sub(a, b)
			case "*":
				result.Mul(a, b)
			case "/":
				if b.Sign() == 0 {
					log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because division by zero", strings.Join(tokens, " "))
					return nil, errs.ErrInvalidAmountExpression
				}
				result.Quo(a, b)
			}

			// push the result back to the stack
			stack = append(stack, result)
		default: // operands
			num, err := parseNumber(token)

			if errors.Is(err, errs.ErrNumericOverflow) {
				log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because numeric overflow", strings.Join(tokens, " "))
				return nil, err
			} else if err != nil {
				log.Warnf(ctx, "[beancount_amount_expression_evaluator.evaluatePostfixExpr] cannot evaluate expression \"%s\", because containing invalid number", strings.Join(tokens, " "))
				return nil, errs.ErrInvalidAmountExpression
			}

			stack = append(stack, num)
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

	return numberToTextualAmount(result)
}
