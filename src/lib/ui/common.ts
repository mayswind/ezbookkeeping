import Clipboard from 'clipboard';

import { ThemeType } from '@/core/theme.ts';

import { type AmountColor, PresetAmountColor } from '@/core/color.ts';

import logger from '../logger.ts';

export function getSystemTheme(): ThemeType {
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        return ThemeType.Dark;
    } else {
        return ThemeType.Light;
    }
}

export function getExpenseAndIncomeAmountColor(expenseAmountColorType: number, incomeAmountColorType: number, isDarkMode?: boolean): AmountColor {
    let expenseAmountColor = expenseAmountColorType ? PresetAmountColor.valueOf(expenseAmountColorType) : null;
    let incomeAmountColor = incomeAmountColorType ? PresetAmountColor.valueOf(incomeAmountColorType) : null;

    if (!expenseAmountColor) {
        expenseAmountColor = PresetAmountColor.DefaultExpenseColor;
    }

    if (!incomeAmountColor) {
        incomeAmountColor = PresetAmountColor.DefaultIncomeColor;
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

export function setExpenseAndIncomeAmountColor(expenseAmountColorType: number, incomeAmountColorType: number): void {
    let expenseAmountColor = expenseAmountColorType ? PresetAmountColor.valueOf(expenseAmountColorType) : null;
    let incomeAmountColor = incomeAmountColorType ? PresetAmountColor.valueOf(incomeAmountColorType) : null;

    if (!expenseAmountColor) {
        expenseAmountColor = PresetAmountColor.DefaultExpenseColor;
    }

    if (!incomeAmountColor) {
        incomeAmountColor = PresetAmountColor.DefaultIncomeColor;
    }

    const htmlElement = document.querySelector('html');

    if (!htmlElement) {
        return;
    }

    const allPresetAmountColors = PresetAmountColor.values();

    for (const amountColor of allPresetAmountColors) {
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

export function copyTextToClipboard(text: string, container?: Element | null): void {
    Clipboard.copy(text, {
        container: container || document.body
    });
}

export function openTextFileContent({ allowedExtensions }: { allowedExtensions: string }): Promise<string> {
    return new Promise((resolve, reject) => {
        const fileInput = document.createElement('input');

        fileInput.style.display = 'none';
        fileInput.type = 'file';
        fileInput.accept = allowedExtensions;

        fileInput.onchange = (event) => {
            const el = event.target as HTMLInputElement;

            if (el.files && el.files.length > 0 && el.files[0]) {
                const file = el.files[0] as File;
                const reader = new FileReader();

                reader.onload = (e) => {
                    const content = e.target?.result;

                    if (typeof content === 'string') {
                        resolve(content);
                    } else {
                        reject(new Error('failed to load file, because file reader result is not string'));
                    }
                };

                reader.onerror = (error) => {
                    logger.error('failed to load file', error);
                };

                reader.readAsText(file);
            }
        };

        fileInput.click();
    });
}

export function startDownloadFile(fileName: string, fileData: Blob): void {
    const dataObjectUrl = URL.createObjectURL(fileData);
    const dataLink = document.createElement('a');

    dataLink.style.display = 'none';
    dataLink.href = dataObjectUrl;
    dataLink.setAttribute('download', fileName);

    document.body.appendChild(dataLink);

    dataLink.click();
}

export function clearBrowserCaches(): Promise<void> {
    if (!window.caches) {
        logger.error('caches API is not supported in this browser');
        return Promise.reject();
    }

    return new Promise((resolve, reject) => {
        window.caches.keys().then(cacheNames => {
            const promises = [];

            for (const cacheName of cacheNames) {
                promises.push(window.caches.delete(cacheName).then(success => {
                    if (success) {
                        logger.info(`cache "${cacheName}" cleared successfully`);
                        return Promise.resolve(cacheName);
                    } else {
                        logger.warn(`failed to clear cache "${cacheName}"`);
                        return Promise.reject(cacheName);
                    }
                }));
            }

            Promise.all(promises).then(() => {
                logger.info("all caches cleared successfully");
                resolve();
            }).catch(() => {
                resolve();
            });
        }).catch(error => {
            logger.warn("failed to clear cache", error);
            reject(error);
        });
    });
}
