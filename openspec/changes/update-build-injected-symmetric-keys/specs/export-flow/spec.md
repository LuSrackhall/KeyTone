# Export Flow - Spec Delta

## ADDED Requirements

### Requirement: Build-injected album export XOR keys

Normative: The system MUST support overriding the `.ktalbum` XOR encryption keys at build time via Go `-ldflags -X` injection, while keeping open-source builds fully functional without any private key files.

- Keys covered by this requirement:
  - v1 key: decrypt legacy `.ktalbum` artifacts
  - v2 key: encrypt/decrypt current `.ktalbum` artifacts

#### Scenario: Export uses current key version

- **GIVEN** the exporter generates a `.ktalbum`
- **WHEN** the system encrypts the album zip payload
- **THEN** it MUST use the current key (v2 by default) and write the key version into the file header

#### Scenario: Import selects key by header version

- **GIVEN** a `.ktalbum` file with header version `V`
- **WHEN** the system decrypts the album payload
- **THEN** it MUST use the key corresponding to `V`

#### Scenario: Import falls back to v1 on checksum failure

- **GIVEN** a `.ktalbum` file with header version not equal to v1
- **WHEN** decryption with the version key fails checksum validation
- **THEN** the system MUST attempt v1 key decryption as a fallback

### Requirement: Build-injected album config encryption seed and inner signature key

Normative: The system MUST support overriding the album-config encryption seed and the album signature inner-field key at build time via Go `-ldflags -X` injection, while keeping open-source builds fully functional without any private key files.

- Keys covered by this requirement:
  - `FixedSecret`: seed used to derive the album config AES key
  - `KeyToneAlbumSignatureEncryptionKey`: inner AES-GCM key for the `signature` field

#### Scenario: No injection keeps default behavior

- **GIVEN** no build-time injection is applied
- **WHEN** the system derives config keys or encrypts/decrypts the `signature` field
- **THEN** it MUST use the default in-source seed/key

#### Scenario: Injection uses deobfuscated plaintext

- **GIVEN** build-time injection is applied for `FixedSecret` and/or `KeyToneAlbumSignatureEncryptionKey`
- **WHEN** the runtime reads the injected value
- **THEN** it MUST deobfuscate it back to plaintext before use
