import { createVault, unlockVault, type Vault, type VaultParams } from '@nicodaimus/crypto';
import { toBase64, fromBase64, toBytes, fromBytes } from '@nicodaimus/crypto';

let _vault: Vault | null = null;

export interface SensitiveTransactionFields {
    amount: number;
    relatedAccountAmount: number;
    comment: string;
    geoLongitude: number;
    geoLatitude: number;
    hideAmount: boolean;
}

export interface VaultParamsDTO {
    vaultVersion: number;
    vaultSalt: string;
    argon2Params: string;
    encryptedDek: string;
    encryptedX25519: string;
    x25519Public: string;
}

export function isVaultUnlocked(): boolean {
    return _vault !== null;
}

export function getVault(): Vault {
    if (!_vault) {
        throw new Error('Vault is locked');
    }
    return _vault;
}

export async function setupVault(passphrase: string): Promise<VaultParamsDTO> {
    const { vault, params } = await createVault(passphrase);
    _vault = vault;
    return vaultParamsToDTO(params);
}

export async function unlock(passphrase: string, dto: VaultParamsDTO): Promise<void> {
    const params = dtoToVaultParams(dto);
    _vault = await unlockVault(passphrase, params);
}

export function lock(): void {
    if (_vault) {
        _vault.shred();
        _vault = null;
    }
}

export function encryptTransactionData(encryptionId: string, data: SensitiveTransactionFields): string {
    const vault = getVault();
    const json = JSON.stringify(data);
    const plaintext = toBytes(json);
    const ciphertext = vault.encryptRecord('tx', encryptionId, plaintext);
    return toBase64(ciphertext);
}

export function decryptTransactionData(encryptionId: string, encryptedDataBase64: string): SensitiveTransactionFields {
    const vault = getVault();
    const ciphertext = fromBase64(encryptedDataBase64);
    const plaintext = vault.decryptRecord('tx', encryptionId, ciphertext);
    const json = fromBytes(plaintext);
    return JSON.parse(json) as SensitiveTransactionFields;
}

export function encryptPicture(pictureId: string, data: Uint8Array): Uint8Array {
    const vault = getVault();
    return vault.encryptRecord('pic', pictureId, data);
}

export function decryptPicture(pictureId: string, data: Uint8Array): Uint8Array {
    const vault = getVault();
    return vault.decryptRecord('pic', pictureId, data);
}

function vaultParamsToDTO(params: VaultParams): VaultParamsDTO {
    return {
        vaultVersion: 1,
        vaultSalt: toBase64(params.argon2Salt),
        argon2Params: JSON.stringify(params.argon2Params),
        encryptedDek: toBase64(params.encryptedDek),
        encryptedX25519: toBase64(params.encryptedX25519Private),
        x25519Public: toBase64(params.x25519Public),
    };
}

function dtoToVaultParams(dto: VaultParamsDTO): VaultParams {
    return {
        argon2Salt: fromBase64(dto.vaultSalt),
        argon2Params: JSON.parse(dto.argon2Params),
        encryptedDek: fromBase64(dto.encryptedDek),
        encryptedX25519Private: fromBase64(dto.encryptedX25519),
        x25519Public: fromBase64(dto.x25519Public),
    };
}
