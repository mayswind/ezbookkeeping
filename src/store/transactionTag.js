import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    LOAD_TRANSACTION_TAG_LIST,
    ADD_TAG_TO_TRANSACTION_TAG_LIST,
    SAVE_TAG_IN_TRANSACTION_TAG_LIST,
    CHANGE_TAG_DISPLAY_ORDER_IN_TRANSACTION_TAG_LIST,
    UPDATE_TAG_VISIBILITY_IN_TRANSACTION_TAG_LIST,
    REMOVE_TAG_FROM_TRANSACTION_TAG_LIST,
    UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE,
} from './mutations.js';

export function loadAllTags(context, { force }) {
    if (!force && !context.state.transactionTagListStateInvalid) {
        return new Promise((resolve) => {
            resolve(context.state.allTransactionTags);
        });
    }

    return new Promise((resolve, reject) => {
        services.getAllTransactionTags().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get tag list' });
                return;
            }

            context.commit(LOAD_TRANSACTION_TAG_LIST, data.result);

            if (context.state.transactionTagListStateInvalid) {
                context.commit(UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            if (force) {
                logger.error('failed to force load tag list', error);
            } else {
                logger.error('failed to load tag list', error);
            }

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get tag list' });
            } else {
                reject(error);
            }
        });
    });
}

export function saveTag(context, { tag }) {
    return new Promise((resolve, reject) => {
        let promise = null;

        if (!tag.id) {
            promise = services.addTransactionTag(tag);
        } else {
            promise = services.modifyTransactionTag(tag);
        }

        promise.then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                if (!tag.id) {
                    reject({ message: 'Unable to add tag' });
                } else {
                    reject({ message: 'Unable to save tag' });
                }
                return;
            }

            if (!tag.id) {
                context.commit(ADD_TAG_TO_TRANSACTION_TAG_LIST, data.result);
            } else {
                context.commit(SAVE_TAG_IN_TRANSACTION_TAG_LIST, data.result);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to save tag', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                if (!tag.id) {
                    reject({ message: 'Unable to add tag' });
                } else {
                    reject({ message: 'Unable to save tag' });
                }
            } else {
                reject(error);
            }
        });
    });
}

export function changeTagDisplayOrder(context, { tagId, from, to }) {
    return new Promise((resolve, reject) => {
        let tag = null;

        for (let i = 0; i < context.state.allTransactionTags.length; i++) {
            if (context.state.allTransactionTags[i].id === tagId) {
                tag = context.state.allTransactionTags[i];
                break;
            }
        }

        if (!tag || !context.state.allTransactionTags[to]) {
            reject({ message: 'Unable to move tag' });
            return;
        }

        if (!context.state.transactionTagListStateInvalid) {
            context.commit(UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE, true);
        }

        context.commit(CHANGE_TAG_DISPLAY_ORDER_IN_TRANSACTION_TAG_LIST, {
            tag: tag,
            from: from,
            to: to
        });

        resolve();
    });
}

export function updateTagDisplayOrders(context) {
    const newDisplayOrders = [];

    for (let i = 0; i < context.state.allTransactionTags.length; i++) {
        newDisplayOrders.push({
            id: context.state.allTransactionTags[i].id,
            displayOrder: i + 1
        });
    }

    return new Promise((resolve, reject) => {
        services.moveTransactionTag({
            newDisplayOrders: newDisplayOrders
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to move tag' });
                return;
            }

            if (context.state.transactionTagListStateInvalid) {
                context.commit(UPDATE_TRANSACTION_TAG_LIST_INVALID_STATE, false);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to save tags display order', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to move tag' });
            } else {
                reject(error);
            }
        });
    });
}

export function hideTag(context, { tag, hidden }) {
    return new Promise((resolve, reject) => {
        services.hideTransactionTag({
            id: tag.id,
            hidden: hidden
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                if (hidden) {
                    reject({ message: 'Unable to hide this tag' });
                } else {
                    reject({ message: 'Unable to unhide this tag' });
                }

                return;
            }

            context.commit(UPDATE_TAG_VISIBILITY_IN_TRANSACTION_TAG_LIST, {
                tag: tag,
                hidden: hidden
            });

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to change tag visibility', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                if (hidden) {
                    reject({ message: 'Unable to hide this tag' });
                } else {
                    reject({ message: 'Unable to unhide this tag' });
                }
            } else {
                reject(error);
            }
        });
    });
}

export function deleteTag(context, { tag, beforeResolve }) {
    return new Promise((resolve, reject) => {
        services.deleteTransactionTag({
            id: tag.id
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to delete this tag' });
                return;
            }

            if (beforeResolve) {
                beforeResolve(() => {
                    context.commit(REMOVE_TAG_FROM_TRANSACTION_TAG_LIST, tag);
                });
            } else {
                context.commit(REMOVE_TAG_FROM_TRANSACTION_TAG_LIST, tag);
            }

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to delete tag', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to delete this tag' });
            } else {
                reject(error);
            }
        });
    });
}
