/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ya3rub/yit/internals"
	"github.com/spf13/cobra"
)

var (
	commitToTag string
	tag         string
	// tagCmd represents the tag command
	tagCmd = &cobra.Command{
		Use:   "tag",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(_ *cobra.Command, _ []string) {
			rootDir, err := os.Getwd()
			mustNil(err, "Error: falied to get current directory")
			db_dir := dbDir(rootDir)
			fmt.Println(db_dir)
			s := &internals.Store{}
			err = s.New(db_dir)
			mustNil(err, "couldn't create new store")
			commitToTagOID := getCommitToTagOID(s)
			mustTrue(len(tag) != 0, "no tag found")
			err = s.CreateTag(tag, commitToTagOID)
			mustNil(err, "couldn't create Tag")
			fmt.Printf("commitToTag: %x\n%s\n", commitToTagOID, tag)
		},
	}
)

func getCommitToTagOID(s internals.Storer) []byte {
	if len(commitToTag) == 0 {
		head, err := s.LookupRef("HEAD")
		mustNil(err, "couldn't get head")
		return head.GetValue()
	}
	dec, err := hex.DecodeString(commitToTag)
	mustNil(err, "couldn't decode commit id")
	return dec
}

func init() {
	rootCmd.AddCommand(tagCmd)

	tagCmd.Flags().
		StringVarP(&commitToTag, "commit", "c", "", "Use the given <msg> as the commit message.")

	tagCmd.Flags().
		StringVarP(&tag, "tag", "t", "", "Use the given <msg> as the commit message.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
