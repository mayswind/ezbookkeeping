import CBOR from 'cbor-js';

import type { ApplicationLockState } from '@/core/setting.ts';
import type { UserBasicInfo } from '@/models/user.ts';

import {
    isFunction,
    stringToArrayBuffer,
    arrayBufferToString,
    base64encode,
    base64decode
} from './common.ts';
import {
    generateRandomString
} from './misc.ts';
import logger from './logger.ts';

interface ClientData {
    challenge: string;
    crossOrigin: boolean;
    origin: string;
    type: string;
}

interface AttestationData {
    authData: Uint8Array;
    fmt: string;
}

interface WebAuthnRegisterResponse {
    id: string | null;
    clientData: ClientData;
    publicKey: Uint8Array | null;
    rawCredential: Credential;
}

interface WebAuthnVerifyResponse {
    id: string | null;
    userName: string;
    userSecret: string;
    clientData: ClientData;
    rawCredential: Credential;
}

const PUBLIC_KEY_CREDENTIAL_CREATION_OPTIONS_BASE_TEMPLATE = {
    attestation: "none",
    authenticatorSelection: {
        authenticatorAttachment: 'platform',
        requireResidentKey: false,
        userVerification: "discouraged"
    },
    pubKeyCredParams: [
        // https://www.iana.org/assignments/cose/cose.xhtml#algorithms
        {type: "public-key", alg: -7},   // ECDSA w/ SHA-256
        {type: "public-key", alg: -257}, // RSASSA-PKCS1-v1_5 using SHA-256
    ],
    timeout: 120000
};

const PUBLIC_KEY_CREDENTIAL_REQUEST_OPTIONS_BASE_TEMPLATE = {
    allowCredentials: [{
        type: 'public-key'
    }],
    userVerification: "discouraged",
    timeout: 120000
};

function parseClientData(credential: Credential): ClientData | null {
    const utf8Decoder = new TextDecoder('utf-8');
    const decodedClientData = utf8Decoder.decode(credential.response.clientDataJSON);
    return JSON.parse(decodedClientData) as ClientData;
}

function parsePublicKeyFromAttestationData(credential: Credential): Uint8Array {
    const decodedAttestationData = CBOR.decode(credential.response.attestationObject) as AttestationData;
    const authData = decodedAttestationData.authData;

    const dataView = new DataView(new ArrayBuffer(2));
    const idLenBytes = authData.slice(53, 55);
    idLenBytes.forEach((value, index) => dataView.setUint8(index, value));

    const credentialIdLength = dataView.getUint16(0);
    const publicKeyBytes = authData.slice(55 + credentialIdLength);

    return publicKeyBytes;
}

export function isWebAuthnSupported(): boolean {
    return !!window.PublicKeyCredential
        && !!navigator.credentials
        && isFunction(window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable);
}

export function isWebAuthnCompletelySupported(): Promise<boolean> {
    if (!isWebAuthnSupported()) {
        return Promise.resolve(false);
    }

    return window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
}

export function registerWebAuthnCredential(lockState: ApplicationLockState, userInfo: UserBasicInfo): Promise<WebAuthnRegisterResponse> {
    if (!window.location || !window.location.hostname) {
        return Promise.reject({
            notSupported: true
        });
    }

    if (!isWebAuthnSupported() || !navigator.credentials.create) {
        return Promise.reject({
            notSupported: true
        });
    }

    const challenge = generateRandomString();
    const userId = `${lockState.username}|${lockState.secret}`; // username 32bytes(max) + secret 24bytes = 56bytes(max)

    const publicKeyCredentialCreationOptions: PublicKeyCredentialCreationOptions = Object.assign({}, PUBLIC_KEY_CREDENTIAL_CREATION_OPTIONS_BASE_TEMPLATE, {
        challenge: stringToArrayBuffer(challenge),
        rp: {
            name: window.location.hostname,
            id: window.location.hostname
        },
        user: {
            id: stringToArrayBuffer(userId),
            name: lockState.username,
            displayName: userInfo.nickname
        }
    }) as PublicKeyCredentialCreationOptions;

    logger.debug('webauthn create options', publicKeyCredentialCreationOptions);

    return navigator.credentials.create({
        publicKey: publicKeyCredentialCreationOptions
    }).then(rawCredential => {
        const clientData = rawCredential ? parseClientData(rawCredential) : null;
        const publicKey = rawCredential ? parsePublicKeyFromAttestationData(rawCredential) : null;

        const challengeFromClientData = clientData && clientData.challenge ? base64decode(clientData.challenge) : null;

        logger.debug('webauthn create raw response', rawCredential);

        if (rawCredential && rawCredential.rawId &&
            clientData && clientData.type === 'webauthn.create' && challengeFromClientData === challenge) {
            const ret: WebAuthnRegisterResponse = {
                id: base64encode(rawCredential.rawId),
                clientData: clientData,
                publicKey: publicKey,
                rawCredential: rawCredential
            };

            logger.debug('webauthn create response', ret);

            return ret;
        } else {
            return Promise.reject({
                invalid: true
            });
        }
    });
}

export function verifyWebAuthnCredential(userInfo: UserBasicInfo, credentialId: string): Promise<WebAuthnVerifyResponse> {
    if (!window.location || !window.location.hostname) {
        return Promise.reject({
            notSupported: true
        });
    }

    if (!isWebAuthnSupported() || !navigator.credentials.get) {
        return Promise.reject({
            notSupported: true
        });
    }

    const challenge = generateRandomString();
    const publicKeyCredentialRequestOptions: PublicKeyCredentialRequestOptions = Object.assign({}, PUBLIC_KEY_CREDENTIAL_REQUEST_OPTIONS_BASE_TEMPLATE, {
        challenge: stringToArrayBuffer(challenge),
        rpId: window.location.hostname
    }) as PublicKeyCredentialRequestOptions;

    if (publicKeyCredentialRequestOptions.allowCredentials && publicKeyCredentialRequestOptions.allowCredentials.length > 0) {
        publicKeyCredentialRequestOptions.allowCredentials[0].id = stringToArrayBuffer(base64decode(credentialId));
    }

    logger.debug('webauthn get options', publicKeyCredentialRequestOptions);

    return navigator.credentials.get({
        publicKey: publicKeyCredentialRequestOptions
    }).then(rawCredential => {
        const clientData = rawCredential ? parseClientData(rawCredential) : null;
        const challengeFromClientData = clientData && clientData.challenge ? base64decode(clientData.challenge) : null;
        const userIdParts = rawCredential && rawCredential.response && rawCredential.response.userHandle ? arrayBufferToString(rawCredential.response.userHandle).split('|') : null;

        logger.debug('webauthn get raw response', rawCredential);

        if (rawCredential && rawCredential.rawId &&
            clientData && clientData.type === 'webauthn.get' && challengeFromClientData === challenge &&
            userIdParts && userIdParts.length === 2 && userIdParts[0] === userInfo.username) {
            const ret: WebAuthnVerifyResponse = {
                id: base64encode(rawCredential.rawId),
                userName: userIdParts[0],
                userSecret: userIdParts[1],
                clientData: clientData,
                rawCredential: rawCredential
            };

            logger.debug('webauthn get response', ret);

            return ret;
        } else {
            return Promise.reject({
                invalid: true
            });
        }
    });
}
