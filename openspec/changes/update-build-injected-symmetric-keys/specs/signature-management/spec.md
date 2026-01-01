# Signature Management - Spec Delta

## ADDED Requirements

### Requirement: Build-injected symmetric signature keys

Normative: The system MUST support overriding the signature symmetric keys at build time via Go `-ldflags -X` injection, while keeping open-source builds fully functional without any private key files.

- Keys covered by this requirement:
  - KeyA: encrypt signature IDs and derive per-signature dynamic keys
  - KeyB: encrypt/decrypt `.ktsign` import/export payloads

#### Scenario: Open-source build without injected keys

- **GIVEN** no build-time injection is applied
- **WHEN** the backend performs signature ID encryption, dynamic key derivation, or `.ktsign` import/export
- **THEN** the system MUST use the default in-source keys

#### Scenario: Private build with injected keys

- **GIVEN** build-time injection is applied for KeyA and/or KeyB
- **WHEN** the backend performs signature ID encryption, dynamic key derivation, or `.ktsign` import/export
- **THEN** the system MUST use the injected keys (after runtime deobfuscation) instead of the defaults

#### Scenario: Injection value format

- **GIVEN** the injected value is an XOR-obfuscated hex string produced by the key-obfuscator tool
- **WHEN** the runtime reads the injected value
- **THEN** the system MUST deobfuscate it back to plaintext before use
