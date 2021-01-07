import services from '../lib/services.js';
import logger from '../lib/logger.js';

import {
    LOAD_TRANSACTION_LIST, UPDATE_ACCOUNT_LIST_INVALID_STATE,
    UPDATE_TRANSACTION_LIST_INVALID_STATE
} from './mutations.js';

function getTransaction(context, { transactionId }) {
    return new Promise((resolve, reject) => {
        services.getTransaction({
            id: transactionId
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                reject({ message: 'Unable to get transaction' });
                return;
            }

            context.commit(LOAD_TRANSACTION_LIST, data.result);
            context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, false);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to load transaction info', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                reject({ message: 'Unable to get transaction' });
            } else {
                reject(error);
            }
        });
    });
}

function saveTransaction(context, { transaction }) {
    return new Promise((resolve, reject) => {
        let promise = null;

        if (!transaction.id) {
            promise = services.addTransaction(transaction);
        } else {
            promise = services.modifyTransaction(transaction);
        }

        promise.then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                if (!transaction.id) {
                    reject({ message: 'Unable to add transaction' });
                } else {
                    reject({ message: 'Unable to save transaction' });
                }
                return;
            }

            context.commit(UPDATE_TRANSACTION_LIST_INVALID_STATE, true);
            context.commit(UPDATE_ACCOUNT_LIST_INVALID_STATE, true);

            resolve(data.result);
        }).catch(error => {
            logger.error('failed to save transaction', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                reject({ error: error.response.data });
            } else if (!error.processed) {
                if (!transaction.id) {
                    reject({ message: 'Unable to add transaction' });
                } else {
                    reject({ message: 'Unable to save transaction' });
                }
            } else {
                reject(error);
            }
        });
    });
}

export default {
    getTransaction,
    saveTransaction
}
