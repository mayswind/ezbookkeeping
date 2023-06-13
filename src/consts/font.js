const allFontSize = {
    Small: {
        type: 0,
        className: 'font-size-small'
    },
    Default: {
        type: 1,
        className: 'font-size-default'
    },
    Large: {
        type: 2,
        className: 'font-size-large'
    },
    XLarge: {
        type: 3,
        className: 'font-size-x-large'
    },
    XXLarge: {
        type: 4,
        className: 'font-size-xx-large'
    },
    XXXLarge: {
        type: 5,
        className: 'font-size-xxx-large'
    },
    XXXXLarge: {
        type: 6,
        className: 'font-size-xxxx-large'
    }
}

const allFontSizeArray = [
    allFontSize.Small,
    allFontSize.Default,
    allFontSize.Large,
    allFontSize.XLarge,
    allFontSize.XXLarge,
    allFontSize.XXXLarge,
    allFontSize.XXXXLarge
];

const defaultFontSize = allFontSize.Default;
const fontSizePreviewClassNamePrefix = 'preview-';

export default {
    allFontSize: allFontSize,
    allFontSizeArray: allFontSizeArray,
    defaultFontSize: defaultFontSize,
    fontSizePreviewClassNamePrefix: fontSizePreviewClassNamePrefix
};
