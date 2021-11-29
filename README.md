# dhc-common
Contains code that is common across digital health cards (vc and eu dgc)

# Card Verifications

Verification covers:
1. Identity level
    1. For v1, SafetyPASS just show the person's picture does NOT name check
2. Vaccination Credential Level
    1. Cards signature verifications
        1. Get issuers public key - passed/failed
        2. Verify card signature with issuers public key - passed/failed/not-checked
    2. Card expired
       1. Verify card not expired using exp - passed/failed/one
    3. The issuer verifications
        1. The issuer is on the CommonTrust or EP3 networks whitelist - passed/failed
    4. The immunization requirements
        1. Is a trusted vaccine - may vary by country - passed/failed
        2. Vaccine specific
            1. The required number of shots have been had
            2. The time between doses was not exceeded, for example 17-92 days
            3. At least some number of days (typically 14) has elapsed since last dose
            4. future - Booster shots TDB how to handle

The card states are as follows, the order is ranked so check in that order

1. **UnKnown** no verifications have been performed
2. **Valid** (Green) 
   1. The card structure verifications have passed 
   2. The card has not expired
   3. The issuer is trusted 
   4. The immunization requirements have been met
3. **Corrupted Card** (Red)
   1. Fetched issuer key and the signature is bad no other checks made
      1. Note invalid cards cannot be loaded, but maybe something happened since loaded, or issuer key changed
4. **UnVerified** (Orange) if cannot check signature then all else is untrusted
   - get key failed so cannot check signature
5. **Safety Criteria Not Met** (Orange) if safety criteria are not met does not matter if issuer is unknown or expired
   - vaccine on whitelist: passed/failed
   - required number shots have been met: passed/failed
   - The time between doses was not exceeded, for example 17-92 days: passed/failed
   - At least some number of days (typically 14) has elapsed since last dose: passed/failed
6. **Issuer Unknown** -(Orange)  if issuer unknown then cannot trust
   - issuer trusted - failed
7. **Expired** - (Orange) if expired but trusted issuer and safety checks made it may be ok
   - card expired