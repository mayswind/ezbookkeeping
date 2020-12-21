export default function ({i18n}, value, fieldName, translate) {
    let content = '';

    if (fieldName) {
        content = value[fieldName];
    } else {
        content = value;
    }

    if (translate) {
        content = i18n.t(content);
    }

    return content;
}
