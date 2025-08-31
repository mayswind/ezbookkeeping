package beancount

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mayswind/ezbookkeeping/pkg/core"
	"github.com/mayswind/ezbookkeeping/pkg/errs"
)

func TestToPostfixExprTokens_ValidExpression(t *testing.T) {
	context := core.NewNullContext()

	result, err := toPostfixExprTokens(context, "1+2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "+"}, result)

	result, err = toPostfixExprTokens(context, "3-4")
	assert.Nil(t, err)
	assert.Equal(t, []string{"3", "4", "-"}, result)

	result, err = toPostfixExprTokens(context, "5*6")
	assert.Nil(t, err)
	assert.Equal(t, []string{"5", "6", "*"}, result)

	result, err = toPostfixExprTokens(context, "8/2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"8", "2", "/"}, result)

	result, err = toPostfixExprTokens(context, "1+2*3-(4/2)")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "3", "*", "+", "4", "2", "/", "-"}, result)

	result, err = toPostfixExprTokens(context, "1 + 2 * 3")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "3", "*", "+"}, result)

	result, err = toPostfixExprTokens(context, "-1+2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"-1", "2", "+"}, result)

	result, err = toPostfixExprTokens(context, "1.5+2.3")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1.5", "2.3", "+"}, result)

	result, err = toPostfixExprTokens(context, "(1+2)-3")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "+", "3", "-"}, result)

	result, err = toPostfixExprTokens(context, "2*-3-3/-2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"2", "-3", "*", "3", "-2", "/", "-"}, result)

	result, err = toPostfixExprTokens(context, "-1.2-3.4*(-5.6/7.8*(9.0-1.2))")
	assert.Nil(t, err)
	assert.Equal(t, []string{"-1.2", "3.4", "-5.6", "7.8", "/", "9.0", "1.2", "-", "*", "*", "-"}, result)

	result, err = toPostfixExprTokens(context, "((((((1+2)*(3+4))))))")
	assert.Nil(t, err)
	assert.Equal(t, []string{"1", "2", "+", "3", "4", "+", "*"}, result)

	result, err = toPostfixExprTokens(context, "(((())))")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, result)

	result, err = toPostfixExprTokens(context, "+-*/")
	assert.Nil(t, err)
	assert.Equal(t, []string{"-", "*", "/", "+"}, result)

	result, err = toPostfixExprTokens(context, "")
	assert.Nil(t, err)
	assert.Equal(t, []string{}, result)
}

func TestToPostfixExprTokens_InvalidExpression(t *testing.T) {
	context := core.NewNullContext()

	_, err := toPostfixExprTokens(context, "1=2")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = toPostfixExprTokens(context, "(1")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = toPostfixExprTokens(context, "2)")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = toPostfixExprTokens(context, "((((1+2)))")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = toPostfixExprTokens(context, ")(")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)
}

func TestEvaluatePostfixExpr_ValidExpression(t *testing.T) {
	context := core.NewNullContext()

	result, err := evaluatePostfixExpr(context, []string{"1", "2", "+"})
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(3000000), result)

	result, err = evaluatePostfixExpr(context, []string{"5", "3", "-"})
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(2000000), result)

	result, err = evaluatePostfixExpr(context, []string{"4", "3", "*"})
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(12000000), result)

	result, err = evaluatePostfixExpr(context, []string{"6", "2", "/"})
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(3000000), result)

	result, err = evaluatePostfixExpr(context, []string{"1", "2", "3", "*", "+", "4", "2", "/", "-"})
	assert.Nil(t, err)
	assert.Equal(t, big.NewInt(5000000), result)
}

func TestEvaluatePostfixExpr_InvalidExpression(t *testing.T) {
	context := core.NewNullContext()

	_, err := evaluatePostfixExpr(context, []string{"1", "0", "/"})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"1", "+"})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"1", "="})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"1", "("})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"1", ")"})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"1", "2", "+", "3"})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluatePostfixExpr(context, []string{"abc"})
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)
}

func TestEvaluateBeancountAmountExpression_ValidExpression(t *testing.T) {
	context := core.NewNullContext()

	result, err := evaluateBeancountAmountExpression(context, "")
	assert.Nil(t, err)
	assert.Equal(t, "", result)

	result, err = evaluateBeancountAmountExpression(context, "1+2")
	assert.Nil(t, err)
	assert.Equal(t, "3.00", result)

	result, err = evaluateBeancountAmountExpression(context, "(1+2)*3")
	assert.Nil(t, err)
	assert.Equal(t, "9.00", result)

	result, err = evaluateBeancountAmountExpression(context, "-1+2")
	assert.Nil(t, err)
	assert.Equal(t, "1.00", result)

	result, err = evaluateBeancountAmountExpression(context, "1.5+2.5")
	assert.Nil(t, err)
	assert.Equal(t, "4.00", result)

	result, err = evaluateBeancountAmountExpression(context, "1+2*3-(4/2)")
	assert.Nil(t, err)
	assert.Equal(t, "5.00", result)

	result, err = evaluateBeancountAmountExpression(context, "2*-3-3/-2")
	assert.Nil(t, err)
	assert.Equal(t, "-4.50", result)

	result, err = evaluateBeancountAmountExpression(context, "-1.2-3.4*(-5.6/7.8*(9.0-1.2))")
	assert.Nil(t, err)
	assert.Equal(t, "17.84", result)

	result, err = evaluateBeancountAmountExpression(context, "(((2+3)))*(((((-5+7)))))")
	assert.Nil(t, err)
	assert.Equal(t, "10.00", result)

	result, err = evaluateBeancountAmountExpression(context, "3.5+0.1")
	assert.Nil(t, err)
	assert.Equal(t, "3.60", result)

	result, err = evaluateBeancountAmountExpression(context, "3.55+0.11")
	assert.Nil(t, err)
	assert.Equal(t, "3.66", result)

	result, err = evaluateBeancountAmountExpression(context, "3.555+0.111")
	assert.Nil(t, err)
	assert.Equal(t, "3.66", result)
}

func TestEvaluateBeancountAmountExpression_InvalidExpression(t *testing.T) {
	context := core.NewNullContext()

	_, err := evaluateBeancountAmountExpression(context, "1++2")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1^2")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "+-*/")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "a+b")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1/0")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1+(2*3")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1+2*3)")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1+((((2*3)))")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1+2(3)")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "1)*(2")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "0.abcd+1")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)

	_, err = evaluateBeancountAmountExpression(context, "0.1234567+1")
	assert.Equal(t, errs.ErrInvalidAmountExpression, err)
}
