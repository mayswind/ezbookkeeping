import CBOR from 'cbor-js';
import logger from './logger.js';
import utils from './utils.js';

const publicKeyCredentialCreationOptionsBaseTemplate = {
    attestation: "none",
    authenticatorSelection: {
        authenticatorAttachment: 'platform',
        requireResidentKey: false,
        userVerification: "required"
    },
    pubKeyCredParams: [
        // https://www.iana.org/assignments/cose/cose.xhtml#algorithms
        {type: "public-key", alg: -7},   // ECDSA w/ SHA-256
        {type: "public-key", alg: -257}, // RSASSA-PKCS1-v1_5 using SHA-256
    ],
    timeout: 1800000
};

const publicKeyCredentialRequestOptionsBaseTemplate = {
    allowCredentials: [{
        type: 'public-key'
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
    const userId = `${username}|${userSecret}`; // username 32bytes(max) + userSecret 24bytes = 56bytes(max)

    const publicKeyCredentialCreationOptions = Object.assign({}, publicKeyCredentialCreationOptionsBaseTemplate, {
        challenge: utils.stringToArrayBuffer(challenge),
        rp: {
            name: window.location.hostname,
            id: window.location.hostname
        },
        user: {
            id: utils.stringToArrayBuffer(userId),
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

    const challenge = utils.generateRandomString();
    const publicKeyCredentialRequestOptions = Object.assign({}, publicKeyCredentialRequestOptionsBaseTemplate, {
        challenge: utils.stringToArrayBuffer(challenge),
        rpId: window.location.hostname
    });
    publicKeyCredentialRequestOptions.allowCredentials[0].id = utils.stringToArrayBuffer(atob(credentialId));

    logger.debug('webauthn get options', publicKeyCredentialRequestOptions);

    return navigator.credentials.get({
        publicKey: publicKeyCredentialRequestOptions
    }).then(rawCredential => {
        const clientData = rawCredential ? parseClientData(rawCredential) : null;
        const challengeFromClientData = clientData && clientData.challenge ? atob(clientData.challenge) : null;
        const userIdParts = rawCredential && rawCredential.response && rawCredential.response.userHandle ? utils.arrayBufferToString(rawCredential.response.userHandle).split('|') : null;

        logger.debug('webauthn get raw response', rawCredential);

        if (rawCredential && rawCredential.rawId &&
            clientData && clientData.type === 'webauthn.get' && challengeFromClientData === challenge &&
            userIdParts && userIdParts.length === 2 && userIdParts[0] === username) {
            const ret = {
                id: utils.base64encode(rawCredential.rawId),
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
