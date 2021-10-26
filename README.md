# dhc-common
Contains code that is common across digital health cards (vc and eu dgc)

# Card Verifications

Verification covers:
1. Identity level
    1. For v1, SafetyPASS just show the person's picture does NOT name check
2. Vaccination Credential Level
    1. Cards structure verifications
        1. Get issuers public key - passed/failed
        2. Verify card signature with issuers public key - passed/failed/not-checked
        3. Verify card not expired using exp - passed/failed/one
    2. The issuer verifications
        1. The issuer is on the CommonTrust or EP3 networks whitelist - passed/failed
    3. The immunization requirements
        1. Is a trusted vaccine - may vary by country - passed/failed
        2. Vaccine specific
            1. The required number of shots have been had
            2. The time between doses was not exceeded, for example 17-92 days
            3. At least some number of days (typically 14) has elapsed since last dose
            4. future - Booster shots TDB how to handle

The card states are
1. **UnKnown** no verifications have been performed
2. **Verified** (Green) 
   1. The card structure verifications have passed 
   2. The issuer is trusted 
   3. The immunization requirements have been met
3. **Corrupted Card** (Red)
   1. Fetched issuer key and the signature is bad no other checks made
      1. Note invalid cards cannot be loaded, but maybe something happened since loaded, or issuer key changed
4. **Partly-Verified** (Yellow) one or more verification failed
   1. card structure
       - get key - passed
       - check signature - passed
       - expired - none
   2. Issuer verifications
       - issuer trusted - failed
   3. immunization verifications
       - vaccine on whitelist: passed/failed
       - required number shots have been met: passed/failed
       - The time between doses was not exceeded, for example 17-92 days: passed/failed
       -  At least some number of days (typically 14) has elapsed since last dose: passed/failed