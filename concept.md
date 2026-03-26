# Concept

## Login

There are a few scenarios we need to review in order to properly handle login.

**Scenario**

1. Receive a request to load the application.
   1. Look for the ClientApp in the client's application storage.
      1. Found
         1. Decrypt the cookie value and unmarshal the JSON and return a
            ClientApp.
            1. The time in the ClientApp will be updated.
            2. Lookup the ClientApp ID in the ProfileMap.
               1. Found
                  1. Lookup the profile.
                  2. Lookup the account.
                  3. Store these profile ID and the Account ID in the session.
                  4. Pull the method they logged in with and verify it's still
                     good.
                     1. Good
                        1. If the ClientApp is marked as trusted, then
                           1. Mark them as logged in.
                           2. Add the login method in the session.
                     2. Not Good
                        1. Then log them out and remove the ClientApp from the
                           Profile map and the Profile.
               2. Not Found, the return﹁
   2. Not Found <-----------------------|
      1. Generate a new ClientApp and store it in the client's application
         storage.
      2. Wait for them to login.

2. Receive a request to login.
   1. Look for the ClientApp in the client's application storage.
      1. Not Found
      2. Found
         1. Lookup the ClientApp ID in the ProfileMap.
            1. Found
               1. Lookup the profile.
               2. Get the Account ID from the Profile and lookup the account.
               3. Store these Profile ID and the Account ID in the session.
               4. Pull the method they logged in with and verify it's still
                  good.
                  1. Good
                     1. If the ClientApp is marked as trusted, then update the
                        session to:
                        1. Set them as logged in.
                        2. Set the login method.
                  2. Not Good
                     1. Then log them out and remove the ClientApp from the
                        Profile map and the Profile.