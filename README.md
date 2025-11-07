# hyoso
> hyōsō (表層) translates to "outer layer" in japanese

extremely minimal bastion/session monitor
## configuration
```toml
[core]
# generate one with: ssh-keygen -t ed25519 -f ~/.hyoso/keys/master_key
master_key = "~/.hyoso/keys/master_key"

# port that hyoso listens on for ssh connections
listen_port = 2223

# where session logs/recordings are stored
log_dir = "~/.hyoso/logs"

# authentication method: "pubkey" | "password" | "custom"
auth_method = "password"

password_type = "sha256" # sha256 | plaintext
password_file = "~/.hyoso/auth/password"

# pubkey auth
# authkey_file = "~/.hyoso/authorized_keys"

# custom auth
# auth_command = "/usr/local/bin/validate_creds.sh"
```
