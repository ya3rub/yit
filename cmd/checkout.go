/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ya3rub/yit/internals"
	"github.com/spf13/cobra"
)

// checkoutCmd represents the checkout command
var (
	commitStr   string
	dir         string
	checkoutCmd = &cobra.Command{
		Use:   "checkout",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(_ *cobra.Command, _ []string) {
			rootDir, err := os.Getwd()
			mustNil(err, "Error: falied to get current directory")
			dst := strings.Join(
				[]string{rootDir, "someDir"},
				string(os.PathSeparator),
			)

			CheckDst(dst)
			db_dir := dbDir(rootDir)
			fmt.Println(db_dir)
			s := &internals.Store{}
			err = s.New(db_dir)

			b, err := s.GetBranch(branchStr)
			commit, err := internals.ParseCommit(s, b.OID)
			mustNil(err, "couldn't parse commit")
			fmt.Printf("oid: %+v\n", commit)
			err = internals.RestoreTree(s , dst , commit.GetTree())
			fmt.Printf("err: %v\n", err)
			ref := internals.NewSymRef("HEAD", b.Name)

			err = s.UpdateRef(ref)
			mustNil(err, "couldn't update Ref")
		},
	}
)
var branchStr string

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.Flags().
		StringVarP(&branchStr, "branch", "b", "", "Use the given <msg> as the commit message.")
	checkoutCmd.Flags().
		StringVarP(&commitStr, "commit", "c", "", "Use the given <msg> as the commit message.")

	checkoutCmd.Flags().
		StringVarP(&commitStr, "dir", "d", "", "Use the given <msg> as the commit message.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func CheckDst(dst string) {
	// fmt.Printf("treeID %x\n", commit.GetTree())
	st, err := os.Stat(dst)
	if os.IsNotExist(err) {
		mustNil(err, fmt.Sprintf("%s doesn't exist", dst))
	}
	if !st.IsDir() {
		mustNil(
			errors.New("NOT_DIR"),
			fmt.Sprintf("%s is not a dir", dst),
		)
	}
	if ents, _ := os.ReadDir(dst); len(ents) != 0 {
		mustNil(
			errors.New("NOT_EMPTY"),
			fmt.Sprintf("%s is not a empty", dst),
		)
	}
}
