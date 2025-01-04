import CBOR from 'cbor-js';
import logger from './logger.ts';
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

const publicKeyCredentialCreationOptionsBaseTemplate = {
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

const publicKeyCredentialRequestOptionsBaseTemplate = {
    allowCredentials: [{
        type: 'public-key'
    }],
    userVerification: "discouraged",
    timeout: 120000
};

function isSupported() {
    return !!window.PublicKeyCredential
        && !!navigator.credentials
        && isFunction(window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable);
}

function isCompletelySupported() {
    if (!isSupported()) {
        return Promise.resolve(false);
    }

    return window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
}

function registerCredential({ username, secret }, { nickname }) {
    if (!window.location || !window.location.hostname) {
        return Promise.reject({
            notSupported: true
        });
    }

    if (!isSupported() || !navigator.credentials.create) {
        return Promise.reject({
            notSupported: true
        });
    }

    const challenge = generateRandomString();
    const userId = `${username}|${secret}`; // username 32bytes(max) + secret 24bytes = 56bytes(max)

    const publicKeyCredentialCreationOptions = Object.assign({}, publicKeyCredentialCreationOptionsBaseTemplate, {
        challenge: stringToArrayBuffer(challenge),
        rp: {
            name: window.location.hostname,
            id: window.location.hostname
        },
        user: {
            id: stringToArrayBuffer(userId),
            name: username,
            displayName: nickname
        }
    });

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
            const ret = {
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

function parseClientData(credential) {
    const utf8Decoder = new TextDecoder('utf-8');
    const decodedClientData = utf8Decoder.decode(credential.response.clientDataJSON);
    return JSON.parse(decodedClientData);
}

function parsePublicKeyFromAttestationData(credential) {
    const decodedAttestationData = CBOR.decode(credential.response.attestationObject);
    const authData = decodedAttestationData.authData;

    const dataView = new DataView(new ArrayBuffer(2));
    const idLenBytes = authData.slice(53, 55);
    idLenBytes.forEach((value, index) => dataView.setUint8(index, value));

    const credentialIdLength = dataView.getUint16();
    const publicKeyBytes = authData.slice(55 + credentialIdLength);

    return publicKeyBytes;
}

function verifyCredential({ username }, credentialId) {
    if (!window.location || !window.location.hostname) {
        return Promise.reject({
            notSupported: true
        });
    }

    if (!isSupported() || !navigator.credentials.get) {
        return Promise.reject({
            notSupported: true
        });
    }

    const challenge = generateRandomString();
    const publicKeyCredentialRequestOptions = Object.assign({}, publicKeyCredentialRequestOptionsBaseTemplate, {
        challenge: stringToArrayBuffer(challenge),
        rpId: window.location.hostname
    });
    publicKeyCredentialRequestOptions.allowCredentials[0].id = stringToArrayBuffer(base64decode(credentialId));

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
            userIdParts && userIdParts.length === 2 && userIdParts[0] === username) {
            const ret = {
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

export default {
    isSupported,
    isCompletelySupported,
    registerCredential,
    verifyCredential
}
