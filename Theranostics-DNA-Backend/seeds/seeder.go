// seeds/seeder.go
package seeds

import (
    "log"

    "theransticslabs/m/utils"
)

// SeedAll runs all seeding functions.
func SeedAll() {
    SeedRoles()
    SeedUsers()

    log.Println(utils.MsgSeedingCompleted)
}
