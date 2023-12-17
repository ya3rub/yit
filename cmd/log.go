/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"fmt"
	"os"

	"github.com/ya3rub/yit/internals"
	"github.com/spf13/cobra"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
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

		commitOID := getCommitOID(s)
		// fmt.Printf("commitOID: %v\n", commitOID)
		commit, err := internals.ParseCommit(s, commitOID)
		mustNil(err, "couldn't parse commit")
		fmt.Printf(
			"(HEAD)  %x\n\tMessage:%s\n",
			commit.ObjID, commit.GetMessage(),
		)

		for len(commit.GetParent()) != 0 {
			commit, err = internals.ParseCommit(s, commit.GetParent())
			fmt.Printf(
				"\t%x\n\tMessage:%s\n",
				commit.ObjID, commit.GetMessage(),
			)
			mustNil(err, "cannot parse commit")
		}

		// fmt.Printf("%+v", commit)
	},
}

// func iterCommits(OIDs []string) []byte {
// 	visited := map[string]bool{}
// }

func getCommitOID(s internals.Storer) []byte {
	if len(tagToLog) == 0 {
		head, err := s.LookupRef("HEAD")
		fmt.Printf("head: %v\n", head)
		mustNil(err, "couldn't get head")
		return head.GetValue()
	}
	oid, err := s.GetTag(tagToLog)
	mustNil(err, "couldn't get tag oid")
	return oid
}

var tagToLog string

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().
		StringVarP(&tagToLog, "tag", "t", "", "Use the given <msg> as the commit message.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
