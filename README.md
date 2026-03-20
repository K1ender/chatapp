# ChatApp (Experimental)

**Status:** 🚧 Project is currently on hold, may resume in the future

## Overview

ChatApp is an experimental end-to-end encrypted chat application with a focus on **passwordless login via magic links**.  

The main design goals were:

- Users register and log in **solely via magic links**, no password required.  
- Each user has a **signing key** (Ed25519) and an **encryption key** (X25519) for end-to-end encryption.  
- Messages are encrypted client-side before being stored in the database.  
- Conversations and messages are stored in PostgreSQL.  

### Key Idea

The experimental idea was to allow users to **recover their chat history without a password**, only by using the magic link. This posed major challenges around:

- Securely storing private keys so they can be restored.
- Allowing multi-device access without compromising end-to-end encryption.
- Ensuring magic links cannot be misused if intercepted.

At the moment, this is **not fully implemented**, as a secure and practical solution proved very difficult.

## Architectural Decisions

- **Database:** PostgreSQL stores users, magic links, conversations, and encrypted messages.
- **Repositories:** Repository layer abstracts database operations, with separate interfaces for users, messages, conversations, and magic links.
- **Service Layer:**  
  - `AuthService` handles magic link generation and verification.  
  - `EmailService` handles sending magic links.  
  - Key generation occurs when a new user is created.  
- **Keys:**  
  - Signing key (Ed25519) for message integrity.  
  - Encryption key (X25519) for encrypting messages.  
  - Initially considered client-side generation and server-side encryption with a derived key from the magic link token.

## Current Limitations

- **Multi-device support:** Private keys tied to a single device make it impossible to log in from multiple devices without losing access.  
- **Magic link recovery:** Fully secure passwordless recovery of private keys is currently unimplemented.  
- **Key storage:** Storing private keys securely while allowing recovery is non-trivial without user passwords or a secure KMS.

## Possible Future Directions

- Explore secure storage solutions like a **Key Management Service (KMS)** or encrypted vaults for multi-device key recovery.  
- Add optional password or passphrase to allow private key recovery across devices.  
- Improve magic link workflow to support temporary, revocable access.  
- Implement proper **end-to-end encryption in multi-device scenarios** with key rotation and sharing.

## Status

This project is **experimental and on hold**. The main reason is the difficulty of balancing **passwordless login** with **secure, recoverable key management**. Future work may revisit these concepts with more secure architectural patterns.