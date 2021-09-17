package cmd

import (
	"fmt"
	"github.com/nais/nais-cli/cmd/helpers"
	"github.com/nais/nais-cli/pkg/aiven"
	aivenclient "github.com/nais/nais-cli/pkg/client"
	"github.com/spf13/cobra"
)

const (
	KafkaNavDev             = "nav-dev"
	KafkaNavProd            = "nav-prod"
	KafkaNavIntegrationTest = "nav-integration-test"
)

var aivenCommand = &cobra.Command{
	Use:   "aiven [command] [args] [flags]",
	Short: "Create a protected & time-limited aivenApplication",
	Long:  `This command will apply a aivenApplication based on information given and aivenator creates a set of credentials`,
	Example: `nais-cli aiven username namespace | nais-cli aiven username namespace -p nav-dev |
nais-cli aiven username namespace -e 10 | nais-cli aiven username namespace -s some-secret-name`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 2 {
			return fmt.Errorf("%s %s %s : reqired arguments", cmd.CommandPath(), UsernameFlag, TeamFlag)
		}
		username := args[0]
		team := args[1]

		pool, _ := helpers.GetString(cmd, PoolFlag, false)
		if pool != KafkaNavDev && pool != KafkaNavProd && pool != KafkaNavIntegrationTest {
			return fmt.Errorf("valid values for '--%s': %s | %s | %s", PoolFlag, KafkaNavDev, KafkaNavProd, KafkaNavIntegrationTest)
		}

		expiry, err := cmd.Flags().GetInt(ExpireFlag)
		if err != nil {
			return fmt.Errorf("getting flag %s", err)
		}

		secretName, err := helpers.GetString(cmd, SecretNameFlag, false)
		if err != nil {
			return fmt.Errorf("getting flag %s", err)
		}

		aivenConfig := aiven.SetupAiven(aivenclient.SetupClient(), username, team, pool, secretName, expiry)
		if _, err := aivenConfig.GenerateApplication(); err != nil {
			return fmt.Errorf("an error occurred generating aivenApplication %s", err)
		}
		return nil
	},
}
