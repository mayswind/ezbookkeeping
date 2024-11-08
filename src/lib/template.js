export function isNoAvailableTemplate(templates, showHidden) {
    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableTemplateCount(templates, showHidden) {
    let count = 0;

    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(templates, showHidden) {
    for (let i = 0; i < templates.length; i++) {
        if (showHidden || !templates[i].hidden) {
            return templates[i].id;
        }
    }

    return null;
}

export function getLastShowingId(templates, showHidden) {
    for (let i = templates.length - 1; i >= 0; i--) {
        if (showHidden || !templates[i].hidden) {
            return templates[i].id;
        }
    }

    return null;
}
