# Login Flow
1. Load a GPG key into the keychain
   1. Use GOpenGPG library.
   2. Pull the GPG key from storage.
   3. Load the GPG keys from file.
2. Store data encrypted with the webapp GPG key into a cookie.
   1. When the user has authenticated, but before you return a response:
      1. Encrypt some info about their login status with the GPG key.
      2. Save this encrypted value in a secure cookie.
3. Decrypt data with the webapp GPG key.
   1. Check if a user is logged in:
      1. Look for a secure cookie.
      2. Take the value and try to decrypt it with the webapp GPG key.
      3. If there is a valid info, then:
         1. pull the account based on user validated info.
         2. Get the device ID.
         3. Use the device ID and search for it in the clients account,
            if there is a match, then:
            1. Pull the provider the logged in with on the device.
            2. Pull the provider login info from storage, if found, then
               1. Check to see if authentication token has expired:
                  1. If not, then restore it.
                  2. If yes, then re-authenticate or get a fresh token.
