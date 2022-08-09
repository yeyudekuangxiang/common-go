/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package certificate

import (
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"mio/internal/pkg/core/app"
	"mio/internal/pkg/core/initialize"
	"mio/internal/pkg/model/entity"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Flags().Parse(args)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		conPath := cmd.Flag("config").Value.String()
		initialize.Initialize(conPath)

		file, err := excelize.OpenFile("cert.xlsx")
		if err != nil {
			log.Panicln(err)
		}
		defer file.Close()

		if file.SheetCount == 0 {
			log.Panicln("没有数据")
		}

		rows, err := file.GetRows(file.GetSheetList()[0])
		if err != nil {
			log.Panicln(err)
		}

		app.DB.Transaction(func(tx *gorm.DB) error {
			certIds := make([]string, 0)
			for i, row := range rows {
				if i <= 1 {
					continue
				}
				certificateId := row[0]
				code := row[1]

				certIds = append(certIds, certificateId)
				stock := entity.CertificateStock{}
				app.DB.Where("certificate_id = ? and code = ?", certificateId, code).First(&stock)
				if stock.ID != 0 {
					continue
				}

				err := app.DB.Create(&entity.CertificateStock{
					CertificateId: certificateId,
					Code:          code,
					Used:          false,
				}).Error
				if err != nil {
					return err
				}

			}

			type CertStock struct {
				ProductItemId string
				Count         int64
			}
			for _, certId := range certIds {
				sql := "select certificate.product_item_id,COUNT (*) COUNT FROM certificate_stock INNER JOIN certificate ON certificate.certificate_id = certificate_stock.certificate_id WHERE certificate_stock.used = FALSE and certificate.certificate_id = ? group by certificate.product_item_id"
				stock := CertStock{}
				tx.Raw(sql, certId).First(&stock)
				if stock.ProductItemId != "" {
					tx.Table("product_item").Where("product_item_id = ?", stock.ProductItemId).UpdateColumn("remaining_count", stock.Count)
				}
			}
			return nil
		})

	},
}

func init() {
	CertificateCmd.AddCommand(importCmd)
	importCmd.Flags().StringP("config", "c", "./config.ini", "config file")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
