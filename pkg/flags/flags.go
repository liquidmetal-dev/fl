package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		viper.BindEnv(f.Name) //nolint: errcheck

		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val)) //nolint: errcheck
		}
	})
}
