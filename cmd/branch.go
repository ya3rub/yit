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

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
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
		oid := getCommitToBranchOID(s)
		fmt.Printf("OID OF BRANCH %x",oid)
		mustTrue(len(branchName) != 0, "no branch name found")
		err = s.CreateBranch(branchName, oid)
		mustNil(err, "couldn't create branch")
	},
}

//AFTER BRANCHING
var (
	commitToBranch string
	branchName     string
)

func getCommitToBranchOID(s internals.Storer) []byte {
	if len(commitToBranch) == 0 {
		fmt.Println("getting head...")
		head, err := s.LookupRef("HEAD")
		fmt.Printf("found head... %x",head.GetValue())
		mustNil(err, "couldn't get head")
		return head.GetValue()
	}
	dec, err := hex.DecodeString(commitToBranch)
	mustNil(err, "couldn't decode commit id")
	return dec
}

func init() {
	rootCmd.AddCommand(branchCmd)

	branchCmd.Flags().
		StringVarP(&commitToBranch, "startCommit", "s", "", "Use the given <msg> as the commit message.")

	branchCmd.Flags().
		StringVarP(&branchName, "name", "n", "", "Use the given <msg> as the commit message.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// branchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// branchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
