package cmd

import (
	"time"
)

const seedTimeout = 30 * time.Second

//var seedCmd = &cobra.Command{
//	Use:   "seed",
//	Short: "run seeders",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		return runSeeder()
//	},
//}
