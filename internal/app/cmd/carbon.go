/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var carbonCmd = &cobra.Command{
	Use:   "carbon",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//service2.DefaultCarbonTransactionService.Bank(service2.BankCarbonTransactionParam{})

		/*a, err := service2.DefaultCarbonTransactionService.Create(service2.CreateCarbonTransactionParam{
			OpenId:  "1",
			UserId:  1,
			Type:    entity.CARBON_COFFEE_CUP,
			Value:   1,
			Info:    fmt.Sprintf("{imageUrl=%s}", 1),
			AdminId: 1,
		})
		if err != nil {
			log.Fatal(err)
		}*/
		//println(a)
	},
}

func init() {
	rootCmd.AddCommand(carbonCmd)
}
