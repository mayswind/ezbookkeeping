export function isNoAvailableTag(tags, showHidden) {
    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            return false;
        }
    }

    return true;
}

export function getAvailableTagCount(tags, showHidden) {
    let count = 0;

    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            count++;
        }
    }

    return count;
}

export function getFirstShowingId(tags, showHidden) {
    for (let i = 0; i < tags.length; i++) {
        if (showHidden || !tags[i].hidden) {
            return tags[i].id;
        }
    }

    return null;
}

export function getLastShowingId(tags, showHidden) {
    for (let i = tags.length - 1; i >= 0; i--) {
        if (showHidden || !tags[i].hidden) {
            return tags[i].id;
        }
    }

    return null;
}
