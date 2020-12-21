export default function ({i18n}, value, fieldName, defaultValue, translate) {
    let content = defaultValue;

    if (fieldName) {
        content = value[fieldName];
    }

    if (translate && content) {
        content = i18n.t(content);
    }

    return content;
}
