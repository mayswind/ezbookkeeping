import type {
    TransactionCreateRequest,
    TransactionModifyRequest,
    TransactionInfoResponse
} from '@/models/transaction.ts';

import {
    isVaultUnlocked,
    encryptTransactionData,
    decryptTransactionData,
    type SensitiveTransactionFields
} from './vault-service.ts';

import { generateRandomUUID } from './misc.ts';

/**
 * Encrypt sensitive fields in a create request before sending to the server.
 * Replaces amount/comment/geo with dummy values and adds encryptedData blob.
 */
export function encryptCreateRequest(req: TransactionCreateRequest): TransactionCreateRequest {
    if (!isVaultUnlocked()) {
        return req;
    }

    const encryptionId = generateRandomUUID();
    const sensitive: SensitiveTransactionFields = {
        amount: req.sourceAmount,
        relatedAccountAmount: req.destinationAmount,
        comment: req.comment,
        geoLongitude: req.geoLocation?.longitude ?? 0,
        geoLatitude: req.geoLocation?.latitude ?? 0,
        hideAmount: req.hideAmount,
    };

    const encryptedData = encryptTransactionData(encryptionId, sensitive);

    return {
        ...req,
        sourceAmount: 0,
        destinationAmount: 0,
        comment: '',
        hideAmount: false,
        geoLocation: undefined,
        encryptionId,
        encryptedData,
    };
}

/**
 * Encrypt sensitive fields in a modify request before sending to the server.
 */
export function encryptModifyRequest(req: TransactionModifyRequest): TransactionModifyRequest {
    if (!isVaultUnlocked()) {
        return req;
    }

    const encryptionId = req.encryptionId || generateRandomUUID();
    const sensitive: SensitiveTransactionFields = {
        amount: req.sourceAmount,
        relatedAccountAmount: req.destinationAmount,
        comment: req.comment,
        geoLongitude: req.geoLocation?.longitude ?? 0,
        geoLatitude: req.geoLocation?.latitude ?? 0,
        hideAmount: req.hideAmount,
    };

    const encryptedData = encryptTransactionData(encryptionId, sensitive);

    return {
        ...req,
        sourceAmount: 0,
        destinationAmount: 0,
        comment: '',
        hideAmount: false,
        geoLocation: undefined,
        encryptionId,
        encryptedData,
    };
}

/**
 * Decrypt a transaction response, merging real values from the encrypted blob
 * over the dummy plaintext values.
 */
export function decryptTransactionResponse(resp: TransactionInfoResponse): TransactionInfoResponse {
    if (!resp.encryptionId || !resp.encryptedData || !isVaultUnlocked()) {
        return resp;
    }

    try {
        const sensitive = decryptTransactionData(resp.encryptionId, resp.encryptedData);

        return {
            ...resp,
            sourceAmount: sensitive.amount,
            destinationAmount: sensitive.relatedAccountAmount,
            comment: sensitive.comment,
            hideAmount: sensitive.hideAmount,
            geoLocation: (sensitive.geoLongitude || sensitive.geoLatitude)
                ? { longitude: sensitive.geoLongitude, latitude: sensitive.geoLatitude }
                : resp.geoLocation,
        };
    } catch {
        // Decryption failed - return response as-is (wrong key or corrupted data)
        return resp;
    }
}

/**
 * Decrypt an array of transaction responses.
 */
export function decryptTransactionResponses(responses: TransactionInfoResponse[]): TransactionInfoResponse[] {
    if (!isVaultUnlocked()) {
        return responses;
    }

    return responses.map(decryptTransactionResponse);
}
