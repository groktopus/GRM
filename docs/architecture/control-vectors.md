# Synthetic control protocol v1 byte vectors

Public synthetic HMAC key hex: `000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f`. Lineage `fixtureA`, key ID `fixtureKey1`. These are deliberately fictional test identifiers and key; production accepts only authenticated server-issued operation identities and independently generated source/grant references. This fixture tests the persistent encoder, not issuance. No CRM/model/runtime execution or production signing occurred.

Genesis head: `7fee20f0219f61745fb70fa74edaa2502d492b269cdea910e031dc0aabf3f88b`.

## 1 operation_reserved

Exact payload P (ASCII, no trailing newline):

```json
{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":1,"previous":"7fee20f0219f61745fb70fa74edaa2502d492b269cdea910e031dc0aabf3f88b","kind":"operation_reserved","payload":{"operation_id":"op_fixtureRandom1","recorded_at":"2026-10-20T12:00:00Z"}}
```

P byte length: 260. Length-prefix hex: `0000000000000104`.

Tag hex: `6b382331d579132f24afb70135d007139334913628acdec2d793992d406f9b5f`.

Entire frame hex:

```text
00000000000001047b2276657273696f6e223a312c226c696e65616765223a226669787475726541222c226b65795f6964223a22666978747572654b657931222c2273657175656e6365223a312c2270726576696f7573223a2237666565323066303231396636313734356662373066613734656461613235303264343932623236396364656139313065303331646330616162663366383862222c226b696e64223a226f7065726174696f6e5f7265736572766564222c227061796c6f6164223a7b226f7065726174696f6e5f6964223a226f705f6669787475726552616e646f6d31222c227265636f726465645f6174223a22323032362d31302d32305431323a30303a30305a227d7d6b382331d579132f24afb70135d007139334913628acdec2d793992d406f9b5f
```

New chain head: `60f733ca7a91e79b5a0510f8f37a9fc00402d47363d90e2047157a9e0e5d11e7`. Exact head/witness file (no trailing newline):

```json
{"body":{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":1,"head":"60f733ca7a91e79b5a0510f8f37a9fc00402d47363d90e2047157a9e0e5d11e7"},"mac":"e451227db722f1718a10124fbe1b18087eb22d1bb8fc4205278788616bb3a916"}
```

## 2 source_deleted

Exact payload P (ASCII, no trailing newline):

```json
{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":2,"previous":"60f733ca7a91e79b5a0510f8f37a9fc00402d47363d90e2047157a9e0e5d11e7","kind":"source_deleted","payload":{"ref":"src_fixtureS7","state":"deleted","recorded_at":"2026-10-02T12:00:00Z","deleted_at":"2026-10-20T12:00:00Z"}}
```

P byte length: 297. Length-prefix hex: `0000000000000129`.

Tag hex: `0d89e56929a3224647943d43d0325caee2c9c34e53dd4dcfa5e3e1c9fd105813`.

Entire frame hex:

```text
00000000000001297b2276657273696f6e223a312c226c696e65616765223a226669787475726541222c226b65795f6964223a22666978747572654b657931222c2273657175656e6365223a322c2270726576696f7573223a2236306637333363613761393165373962356130353130663866333761396663303034303264343733363364393065323034373135376139653065356431316537222c226b696e64223a22736f757263655f64656c65746564222c227061796c6f6164223a7b22726566223a227372635f666978747572655337222c227374617465223a2264656c65746564222c227265636f726465645f6174223a22323032362d31302d30325431323a30303a30305a222c2264656c657465645f6174223a22323032362d31302d32305431323a30303a30305a227d7d0d89e56929a3224647943d43d0325caee2c9c34e53dd4dcfa5e3e1c9fd105813
```

New chain head: `df62f52903ee9a116b4573ff89a09fa3603e66c6253d3de04042b012b2d478e1`. Exact head/witness file (no trailing newline):

```json
{"body":{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":2,"head":"df62f52903ee9a116b4573ff89a09fa3603e66c6253d3de04042b012b2d478e1"},"mac":"48cd97657bece801933669d7edfcfa33920d881cf0ae9b94edc80374cdefd903"}
```

## 3 grant_revoked

Exact payload P (ASCII, no trailing newline):

```json
{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":3,"previous":"df62f52903ee9a116b4573ff89a09fa3603e66c6253d3de04042b012b2d478e1","kind":"grant_revoked","payload":{"grant":"fixtureGrant9","state":"revoked","recorded_at":"2026-10-20T12:00:00Z"}}
```

P byte length: 262. Length-prefix hex: `0000000000000106`.

Tag hex: `2aac3c437cfc656fc7565f2cf6544978a5c236b05269b8c4d264c219fcf74197`.

Entire frame hex:

```text
00000000000001067b2276657273696f6e223a312c226c696e65616765223a226669787475726541222c226b65795f6964223a22666978747572654b657931222c2273657175656e6365223a332c2270726576696f7573223a2264663632663532393033656539613131366234353733666638396130396661333630336536366336323533643364653034303432623031326232643437386531222c226b696e64223a226772616e745f7265766f6b6564222c227061796c6f6164223a7b226772616e74223a22666978747572654772616e7439222c227374617465223a227265766f6b6564222c227265636f726465645f6174223a22323032362d31302d32305431323a30303a30305a227d7d2aac3c437cfc656fc7565f2cf6544978a5c236b05269b8c4d264c219fcf74197
```

New chain head: `f0f0eb9da884051321f01789e9f1609da5d2193a0641cf12774c39b303f81d53`. Exact head/witness file (no trailing newline):

```json
{"body":{"version":1,"lineage":"fixtureA","key_id":"fixtureKey1","sequence":3,"head":"f0f0eb9da884051321f01789e9f1609da5d2193a0641cf12774c39b303f81d53"},"mac":"0f16c7ccc01f4149fe88d01e5f59bb276867ee6c07251c6070180f7bedc17322"}
```

## Grant signing-message encoding

Domain bytes are ASCII `GRM-OWNER-GRANT-1` followed by one LF byte `0a`, not a literal backslash and n. Then fixed ordered compact ASCII JSON keys: version, installation, client, operation_id, plan_id, plan_digest, action, target, authority_epoch, key_id, nonce, issued_at, expires_at, session, base_revision. No escapes/whitespace/extra fields; all non-time strings use the registry alphabets. The example is a canonicalization fixture only; no authentic Ed25519 signature is supplied or asserted verified.

Exact JSON:

```json
{"version":"1","installation":"fixtureA","client":"fixtureClient","operation_id":"op_fixtureRandom1","plan_id":"plan_fixture1","plan_digest":"fixtureDigest","action":"delete_source","target":"src_fixtureS7","authority_epoch":"epoch1","key_id":"fixtureKey1","nonce":"fixtureNonce","issued_at":"2026-10-20T12:00:00Z","expires_at":"2026-10-20T12:05:00Z","session":"fixtureSession","base_revision":"r20"}
```

Message byte length: 418. Exact signing-message hex:

```text
47524d2d4f574e45522d4752414e542d310a7b2276657273696f6e223a2231222c22696e7374616c6c6174696f6e223a226669787475726541222c22636c69656e74223a2266697874757265436c69656e74222c226f7065726174696f6e5f6964223a226f705f6669787475726552616e646f6d31222c22706c616e5f6964223a22706c616e5f6669787475726531222c22706c616e5f646967657374223a2266697874757265446967657374222c22616374696f6e223a2264656c6574655f736f75726365222c22746172676574223a227372635f666978747572655337222c22617574686f726974795f65706f6368223a2265706f636831222c226b65795f6964223a22666978747572654b657931222c226e6f6e6365223a22666978747572654e6f6e6365222c226973737565645f6174223a22323032362d31302d32305431323a30303a30305a222c22657870697265735f6174223a22323032362d31302d32305431323a30353a30305a222c2273657373696f6e223a226669787475726553657373696f6e222c22626173655f7265766973696f6e223a22723230227d
```

Message SHA256 (encoding check, not signature): `90065c5527120e0798a6bd0303cd484f95d653dec640fdb00071f57a73b1be7c`.

Reject every single-field mutation against its original signature, including appended session/base. Persistent log negative vectors derive by mutating exactly one indicated byte/property of the fixed frames above: first length byte to ff (reject before allocation), final tag byte XOR01 (complete invalid MAC, no truncation), remove final tag byte (truncated suffix only beyond all witnessed heads), sequence=previous sequence, lineage=fixtureB, kind=unknown, previous=64 zeros, extra field, duplicated version, whitespace/newline suffix. Recomputing a valid MAC for altered sequence/lineage/type still fails semantic chain/type validation. These latter fixtures use the public test key; an attacker cannot recompute production MAC.
