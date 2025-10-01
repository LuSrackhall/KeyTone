# Contract: Signature Operations

## List Signatures
GET /sdk/signatures
Response: 200 [{ id, name, intro?, cardImagePath?, createdAt }]

## Create Signature
POST /sdk/signatures
Body: { name, intro?, cardImageBase64? }
Response: 201 { id }
Errors: 400 (duplicate name)

## Delete Signature
DELETE /sdk/signatures/{id}
Response: 204

## Export Signature
POST /sdk/signatures/{id}/export
Body: { }
Response: 200 { fileNameSuggested: string, blobBase64: string }

## Import Signature
POST /sdk/signatures/import
Body: { fileName: string, blobBase64: string, overwrite?: boolean }
Response: 201 { id, overwritten: boolean }
Errors: 409 (exists without overwrite)
