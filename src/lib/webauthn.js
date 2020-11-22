import CBOR from 'cbor-js';
import logger from './logger.js';
import utils from './utils.js';

const PUBLIC_KEY_CREDENTIAL_CREATION_OPTIONS_TEMPLATE = {
    attestation: "none",
    authenticatorSelection: {
        authenticatorAttachment: 'platform',
        requireResidentKey: false,
        userVerification: "required"
    },
    pubKeyCredParams: [
        // https://www.iana.org/assignments/cose/cose.xhtml#algorithms
        {type: "public-key", alg: -7},   // ECDSA w/ SHA-256
        {type: "public-key", alg: -35},  // ECDSA w/ SHA-384
        {type: "public-key", alg: -36},  // ECDSA w/ SHA-512
        {type: "public-key", alg: -257}, // RSASSA-PKCS1-v1_5 using SHA-256
        {type: "public-key", alg: -258}, // RSASSA-PKCS1-v1_5 using SHA-384
        {type: "public-key", alg: -259}, // RSASSA-PKCS1-v1_5 using SHA-512
        {type: "public-key", alg: -37},  // RSASSA-PSS w/ SHA-256
        {type: "public-key", alg: -38},  // RSASSA-PSS w/ SHA-384
        {type: "public-key", alg: -39},  // RSASSA-PSS w/ SHA-512
        {type: "public-key", alg: -8}    // EdDSA
    ],
    timeout: 1800000
};

const PUBLIC_KEY_CREDENTIAL_REQUEST_OPTIONS_TEMPLATE = {
    allowCredentials: [{
        type: 'public-key',
        transports: ['internal']
    }],
    userVerification: "required",
    timeout: 1800000
}

function isSupported() {
    return !!window.PublicKeyCredential && !!navigator.credentials;
}

function isCompletelySupported() {
    if (!isSupported()) {
        return Promise.resolve(false);
    }

    return window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
}

function registerCredential({ username, nickname }, userSecret) {
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

    const challenge = utils.generateRandomString();
    const publicKeyCredentialCreationOptions = Object.assign({}, PUBLIC_KEY_CREDENTIAL_CREATION_OPTIONS_TEMPLATE, {
        challenge: Uint8Array.from(challenge, c => c.charCodeAt(0)),
        rp: {
            name: window.location.hostname,
            id: window.location.hostname
        },
        user: {
            id: Uint8Array.from(userSecret, c => c.charCodeAt(0)),
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

        const challengeFromClientData = clientData && clientData.challenge ? atob(clientData.challenge) : null;

        logger.debug('webauthn create raw response', rawCredential);

        if (rawCredential && rawCredential.rawId &&
            clientData && clientData.type === 'webauthn.create' && challengeFromClientData === challenge) {
            const ret = {
                id: utils.base64encode(rawCredential.rawId),
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

function verifyCredential(credentialId) {
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

    const challenge = utils.generateRandomString();
    const publicKeyCredentialRequestOptions = Object.assign({}, PUBLIC_KEY_CREDENTIAL_REQUEST_OPTIONS_TEMPLATE, {
        challenge: Uint8Array.from(challenge, c => c.charCodeAt(0)),
        rpId: window.location.hostname
    });
    publicKeyCredentialRequestOptions.allowCredentials[0].id = Uint8Array.from(atob(credentialId), c=>c.charCodeAt(0)).buffer;

    logger.debug('webauthn get options', publicKeyCredentialRequestOptions);

    return navigator.credentials.get({
        publicKey: publicKeyCredentialRequestOptions
    }).then(rawCredential => {
        const clientData = rawCredential ? parseClientData(rawCredential) : null;
        const challengeFromClientData = clientData && clientData.challenge ? atob(clientData.challenge) : null;

        logger.debug('webauthn get raw response', rawCredential);

        if (rawCredential && rawCredential.rawId &&
            rawCredential.response && rawCredential.response.userHandle &&
            clientData && clientData.type === 'webauthn.get' && challengeFromClientData === challenge) {

            const ret = {
                id: utils.base64encode(rawCredential.rawId),
                userSecret: utils.arrayBufferToString(rawCredential.response.userHandle),
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
