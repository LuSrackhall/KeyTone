# Contract: Album Signing

## List Album Signatures
GET sdk://albums/{albumId}/signatures
Response: 200 [{ signerId, signedAt: [timestamp...], allowReexport }]

## Sign Album (First or Additional)
POST sdk://albums/{albumId}/sign
Body: { signatureId, allowReexport?: boolean }
Response: 201 { updated: true }
Errors: 400 (invalid signatureId)

## Export Album With Signatures
POST sdk://albums/{albumId}/export
Body: { signatureId?: string }
Response: 200 { fileNameSuggested, blobBase64 }

