package devicecmd

import (
	"fmt"
	"github.com/nais/cli/pkg/naisdevice"
	"github.com/urfave/cli/v2"
	"k8s.io/utils/strings/slices"
)

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "Shows the status of your naisdevice",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Action: func(context *cli.Context, flag string) error {
					if !slices.Contains([]string{"yaml", "json"}, flag) {
						return fmt.Errorf("%v is not a implemented format\n", flag)
					}

					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
			},
		},
		Action: func(context *cli.Context) error {
			outputFormat := context.String("output")
			verbose := context.Bool("verbose")

			status, err := naisdevice.GetStatus(context.Context)
			if err != nil {
				return err
			}

			if outputFormat != "" {
				return naisdevice.PrintFormattedStatus(outputFormat, status)
			}

			if verbose {
				naisdevice.PrintVerboseStatus(status)
				return nil
			}

			fmt.Println(status.ConnectionState.String())

			return nil
		},
	}
}