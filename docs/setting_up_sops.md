# **Team SOPS Setup Guide**

## **Why SOPS?**

SOPS encrypts secrets (API keys, passwords) so we can safely store them in Git. Each team member gets a PGP key that allows them to decrypt the shared secrets file. Without Slack messages with passwords

**Time needed:** 10-20 minutes for initial setup, 3 minutes for daily use.

---

## **Quick Commands Reference**

```bash
# Decrypt secrets for development
sops -d .env.encrypted > .env

# Edit encrypted secrets
sops .env.encrypted

# Add new team member (after getting their public key)
gpg --import newmember_public_key.asc
# Add their fingerprint to .sops.yaml, then:
sops -e .env > .env.encrypted
```

---

## **Prerequisites**

- **Your OS** (Arch linux assumed: `sudo pacman -Syu`)
- **Git repository** set up
- **Secure communication** channel (Signal, encrypted Slack, etc.)
- **SOPS installed**: `sudo pacman -S sops`

---

## **Part A: First-Time Setup (Project Admin)**

### **1. Install GnuPG**

```bash
sudo pacman -S gnupg
gpg --version  # Should show GnuPG 2.x.x
```

### **2. Generate Your PGP Key**

```bash
gpg --full-generate-key
```

**Prompts:**

- **Key Type**: `RSA and RSA` (default, press Enter)
- **Key Size**: `2048` (fine for team use)
- **Expiration**: `0` (no expiration) or `1y` (1 year)
- **Name/Email**: Use your name + your email (e.g., `dd@proton.me`)
- **Passphrase**: Strong passphrase → store in password manager

### **3. Export and Share Your Public Key**

```bash
# Find your key ID
gpg --list-keys

# Export (replace ABCDEF1234567890 with your actual key ID)
gpg --armor --export ABCDEF1234567890 > dd_public_key.asc
```

**Share `dd_public_key.asc` via secure channel (Signal, etc.)**

### **4. Collect Team Keys**

Each team member sends you their `.asc` file:

```bash
gpg --import bob_public_key.asc
gpg --import charlie_public_key.asc
gpg --list-keys  # Verify all keys imported
```

### **5. Configure SOPS**

Create `.sops.yaml` in project root:

```bash
vim .sops.yaml
```

Add all team fingerprints:

```yaml
creation_rules:
  - path_regex: \.env$
    pgp: |
      ABCDEF1234567890ABCDEF1234567890ABCDEF12  # dd
      1234567890ABCDEF1234567890ABCDEF12345678  # Bob  
      567890ABCDEF1234567890ABCDEF1234567890AB  # Charlie
```

**Get fingerprints:**

```bash
gpg --list-keys --fingerprint
# Copy 40-character fingerprint (remove spaces)
```

### **6. Encrypt and Commit Secrets**

```bash
# Create your secrets file
vim .env
# Add: API_KEY=sk_1234567890abcdef
#      DB_PASSWORD=supersecretpassword

# Encrypt
sops -e .env > .env.encrypted

# Add to Git
git add .sops.yaml .env.encrypted
echo ".env" >> .gitignore
git add .gitignore
git commit -m "Add encrypted secrets and SOPS config"
git push
```

**Optional:** Create `.env.example` with placeholder values for new team members.

---

## **Part B: Team Member Setup**

### **1. Install Tools**

```bash
# Arch Linux
sudo pacman -S gnupg sops

# macOS: brew install gnupg sops
# Ubuntu: apt install gnupg sops
```

### **2. Generate Your Key**

```bash
gpg --full-generate-key
# Follow same prompts as admin above
```

### **3. Share Your Public Key**

```bash
gpg --list-keys  # Find your key ID
gpg --armor --export YOUR_KEY_ID > yourname_public_key.asc
# Send yourname_public_key.asc to admin via secure channel
```

### **4. Import Team Keys**

Admin will send you other team members' public keys:

```bash
gpg --import dd_public_key.asc
gpg --import bob_public_key.asc
```

### **5. Clone and Decrypt**

```bash
git clone <repo-url>
cd <project-directory>
sops -d .env.encrypted > .env
# Enter your passphrase when prompted
```

---

## **Daily Usage**

### **Developer Workflow**

```bash
# Get latest secrets
git pull
sops -d .env.encrypted > .env

# Your app now loads from .env as usual
```

### **Updating Secrets**

```bash
# Edit encrypted file directly
sops .env.encrypted
# SOPS opens editor, make changes, save

# Commit changes
git add .env.encrypted
git commit -m "Update API keys"
git push
```

### **Adding New Team Member**

1. Get their public key: `yourname_public_key.asc`
2. Import: `gpg --import yourname_public_key.asc`
3. Add fingerprint to `.sops.yaml`
4. Re-encrypt: `sops -e .env > .env.encrypted`
5. Commit: `git add .sops.yaml .env.encrypted && git commit -m "Add new team member"`

### **Team Member Leaving**

1. Remove their fingerprint from `.sops.yaml`
2. Rotate all secrets (generate new API keys, passwords)
3. Update `.env` with new secrets
4. Re-encrypt: `sops -e .env > .env.encrypted`
5. Commit and push

---

## **Troubleshooting**

| **Problem**       | **Solution**                                       |
| ----------------- | -------------------------------------------------- |
| `sops -d` fails   | Check: `gpg --list-secret-keys` (need private key) |
| "Key not found"   | Verify your fingerprint is in `.sops.yaml`         |
| Forgot passphrase | Generate new key pair, update `.sops.yaml`         |
| Can't import key  | Check file exists: `ls *.asc`                      |

**Need help?** Run `gpg --list-keys --fingerprint` and share output (fingerprints are safe to share).

---

## **Optional: Auto-loading with direnv**

Auto-decrypt secrets when entering project directory:

```bash
# Install direnv
sudo pacman -S direnv
echo 'eval "$(direnv hook bash)"' >> ~/.bashrc
source ~/.bashrc

# Create .envrc in project root
echo "source <(sops -d .env.encrypted)" > .envrc
direnv allow
```

Now `.env` variables load automatically when you `cd` into the project.

---

## **Security Notes**

- **Private keys**: Never share, store passphrases in password manager
- **Key exchange**: Use encrypted channels (Signal, not email)
- **Team changes**: Always rotate secrets when members leave
- **Backup**: Keep encrypted `.env.encrypted` in Git, never plain `.env`

