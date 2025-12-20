import Clipboard from 'clipboard';

import { ThemeType } from '@/core/theme.ts';

import { type AmountColor, PresetAmountColor } from '@/core/color.ts';
import { KnownFileType } from '@/core/file.ts';

import logger from '../logger.ts';

export function scrollToSelectedItem(parentEl: Element | null | undefined, containerSelector: string | null, scrollableListSelector: string | null, selectedItemSelector: string): void {
    if (!parentEl) {
        return;
    }

    const container = containerSelector ? parentEl.querySelector<HTMLElement>(containerSelector) : parentEl;

    if (!container) {
        return;
    }

    const scrollableList = scrollableListSelector ? parentEl.querySelector<HTMLElement>(scrollableListSelector) : parentEl;

    if (!scrollableList) {
        return;
    }

    const selectedItems = scrollableList.querySelectorAll<HTMLElement>(selectedItemSelector);

    if (!selectedItems.length) {
        return;
    }

    const containerHeight: number = container.clientHeight;

    const firstSelectedItem: HTMLElement = selectedItems[0] as HTMLElement;
    const lastSelectedItem: HTMLElement = selectedItems[selectedItems.length - 1] as HTMLElement;

    const firstSelectedItemHeight: number = firstSelectedItem.offsetHeight;
    const firstSelectedItemTop: number = firstSelectedItem.offsetTop;
    const lastSelectedItemTop: number = lastSelectedItem.offsetTop;
    const lastSelectedItemBottom: number = lastSelectedItem.offsetTop + lastSelectedItem.offsetHeight;

    const middle: number = (firstSelectedItemTop + lastSelectedItemBottom) / 2;
    let targetScrollTop: number = middle - containerHeight / 2;

    if (containerSelector !== scrollableListSelector) {
        const scrollableListStyle = window.getComputedStyle(scrollableList);
        const paddingTop: number = parseFloat(scrollableListStyle.paddingTop) || 0;
        targetScrollTop += paddingTop / 3 * 2;

        if (selectedItems.length > 1 && lastSelectedItemTop - firstSelectedItemTop > containerHeight - firstSelectedItemHeight - paddingTop) {
            targetScrollTop = firstSelectedItemTop;
        }
    } else {
        if (selectedItems.length > 1 && lastSelectedItemTop - firstSelectedItemTop > containerHeight - firstSelectedItemHeight) {
            targetScrollTop = firstSelectedItemTop;
        }
    }

    scrollableList.scrollTop = Math.max(0, targetScrollTop);
}

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

export function compressJpgImage(file: File, maxWidth: number, maxHeight: number, quality: number): Promise<Blob> {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = (event) => {
            const img = new Image();

            img.onload = () => {
                let width = img.width;
                let height = img.height;

                if (width > maxWidth || height > maxHeight) {
                    const scale = Math.min(maxWidth / width, maxHeight / height);
                    width = Math.floor(width * scale);
                    height = Math.floor(height * scale);
                }

                const canvas = document.createElement('canvas');
                const ctx = canvas.getContext('2d');

                if (!ctx) {
                    reject(new Error('failed to get canvas context'));
                    return;
                }

                canvas.width = width;
                canvas.height = height;

                ctx.drawImage(img, 0, 0, width, height);

                canvas.toBlob((blob) => {
                    if (blob) {
                        resolve(blob);
                    } else {
                        reject(new Error('failed to compress image'));
                    }
                }, KnownFileType.JPG.contentType, quality);
            };

            img.onerror = (error) => {
                reject(error);
            };

            if (event.target && event.target.result) {
                img.src = event.target.result as string;
            } else {
                reject(new Error('failed to read file'));
            }
        };

        reader.onerror = (error) => {
            reject(error);
        };

        reader.readAsDataURL(file);
    });
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
