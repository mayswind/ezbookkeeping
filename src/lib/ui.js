import colorConstants from '@/consts/color.js';

export function getSystemTheme() {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return 'dark';
    } else {
        return 'light';
    }
}

export function getExpenseAndIncomeAmountColor(expenseAmountColorType, incomeAmountColorType, isDarkMode) {
    let expenseAmountColor = expenseAmountColorType ? colorConstants.allAmountColorTypesMap[expenseAmountColorType] : null;
    let incomeAmountColor = incomeAmountColorType ? colorConstants.allAmountColorTypesMap[incomeAmountColorType] : null;

    if (!expenseAmountColor) {
        expenseAmountColor = colorConstants.defaultExpenseAmountColor;
    }

    if (!incomeAmountColor) {
        incomeAmountColor = colorConstants.defaultIncomeAmountColor;
    }

    if (isDarkMode) {
        return {
            expenseAmountColor: expenseAmountColor.darkThemeColor,
            incomeAmountColor: incomeAmountColor.darkThemeColor
        }
    } else {
        return {
            expenseAmountColor: expenseAmountColor.lightThemeColor,
            incomeAmountColor: incomeAmountColor.lightThemeColor
        }
    }
}

export function setExpenseAndIncomeAmountColor(expenseAmountColorType, incomeAmountColorType) {
    let expenseAmountColor = expenseAmountColorType ? colorConstants.allAmountColorTypesMap[expenseAmountColorType] : null;
    let incomeAmountColor = incomeAmountColorType ? colorConstants.allAmountColorTypesMap[incomeAmountColorType] : null;

    if (!expenseAmountColor) {
        expenseAmountColor = colorConstants.defaultExpenseAmountColor;
    }

    if (!incomeAmountColor) {
        incomeAmountColor = colorConstants.defaultIncomeAmountColor;
    }

    const htmlElement = document.querySelector('html');

    for (let i = 0; i < colorConstants.allAmountColorsArray.length; i++) {
        const amountColor = colorConstants.allAmountColorsArray[i];

        if (amountColor.type === expenseAmountColor.type) {
            if (!htmlElement.classList.contains(amountColor.expenseClassName)) {
                htmlElement.classList.add(amountColor.expenseClassName);
            }
        } else {
            htmlElement.classList.remove(amountColor.expenseClassName);
        }

        if (amountColor.type === incomeAmountColor.type) {
            if (!htmlElement.classList.contains(amountColor.incomeClassName)) {
                htmlElement.classList.add(amountColor.incomeClassName);
            }
        } else {
            htmlElement.classList.remove(amountColor.incomeClassName);
        }
    }
}

export function startDownloadFile(fileName, fileData) {
    const dataObjectUrl = URL.createObjectURL(fileData);
    const dataLink = document.createElement('a');

    dataLink.style.display = 'none';
    dataLink.href = dataObjectUrl;
    dataLink.setAttribute('download', fileName);

    document.body.appendChild(dataLink);

    dataLink.click();
}
