/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"fmt"
	"os"

	"github.com/ya3rub/yit/internals"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		rootDir, err := os.Getwd()
		mustNil(err, "Error: falied to get current directory")
		db_dir := dbDir(rootDir)
		fmt.Println(db_dir)
		s := &internals.Store{}
		err = s.New(db_dir)
		ref := internals.NewSymRef(
			"HEAD",
			internals.OsPathJoin("refs", "heads", "main"),
		)
		err = s.UpdateRef(ref)
		mustNil(err, "couldn't update Ref")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
