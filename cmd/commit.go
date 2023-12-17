/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ya3rub/yit/internals"
	"github.com/spf13/cobra"
)

type ConfigVars struct {
	// TODO: read from config file
	YIT_AUTHOR_NAME  string
	YIT_AUTHOR_EMAIL string
}

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		author := genAuthor()
		rootDir, err := os.Getwd()
		mustNil(err, "Error: falied to get current directory")
		db_dir := dbDir(rootDir)
		s := &internals.Store{}
		err = s.New(db_dir)
		mustNil(err, "Error: falied to create a new db")
		rootTree, err := storeTree(s, rootDir)
		parsedTree, err := internals.ParseTree(s, rootTree.GetObjID())
		mustNil(err, "cannot parse tree")
		fmt.Printf("parsedTree: %x\n", parsedTree.GetObjID())
		parent, err := s.LookupRef("HEAD")
		commit := StoreCommit(
			s,
			parent.GetValue(),
			rootTree,
			message,
			author,
			author,
		)

		ref := internals.NewRef("HEAD", commit.GetObjID())
		err = s.UpdateRef(ref)
		mustNil(err, "Cannot set Head")
		fmt.Printf(
			"[(ROOT_COMMIT) %x] %s\n",
			commit.GetObjID(),
			commit.GetMessage(),
		)
	},
}

func storeTree(s internals.Storer, dir string) (*internals.Tree, error) {
	// fmt.Println(dir)
	files, err := os.ReadDir(dir)
	mustNil(err, "Error: falied to read the current directory")
	var blobs []internals.Blob
	var entries []internals.Entry
	for _, dirEntry := range files {
		currPath := strings.Join(
			[]string{dir, dirEntry.Name()},
			string(os.PathSeparator),
		)
		if dirEntry.Name() == "." || dirEntry.Name() == ".." {
			continue
		}
		if dirEntry.IsDir() {
			if dirEntry.Name() == ".git" || dirEntry.Name() == ".yit" {
				continue
			}
			// fmt.Println("scanning...", currPath)
			tree, err := storeTree(s, currPath)
			// fmt.Printf("%+v", tree)
			mustNil(err, "tree gen error")
			entries = append(
				entries,
				genEntry(tree.GetObjID(), dirEntry.Name()),
			)

			continue
		}

		// fmt.Println("openning...", currPath)
		file, err := os.Open(currPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		blob := genBlob(file, dirEntry.Name())
		s.Store(&blob)
		// fmt.Printf("%+v \n", blob.GetObjID())
		blobs = append(blobs, blob)
	}
	// fmt.Printf("gen entries for blobs: %+v", blobs)
	entries = append(entries, genEntries(blobs)...)
	tree := internals.Tree{}
	treeObjID, err := internals.GenObjID(&tree)
	mustNil(err, "couldn't genereate obj for tree")
	err = tree.New(entries)
	mustNil(err, "Error: falied to create a new tree")
	tree.SetObjID(treeObjID)
	err = s.Store(&tree)
	mustNil(
		err,
		fmt.Sprintf(
			"Error: falied to store tree with obj ID: %x\n",
			tree.GetObjID(),
		),
	)
	return &tree, nil
}

func genBlob(file *os.File, name string) internals.Blob {
	var cnt bytes.Buffer
	_, err := io.Copy(&cnt, file)
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: falied to copy the bytes. %v\n",
			err,
		)
		os.Exit(1)
	}
	blob := internals.Blob{}
	if err := blob.New(cnt.Bytes(), name); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: falied to create the blob %v\n",
			err,
		)
		os.Exit(1)
	}
	objID, err := internals.GenObjID(&blob)
	blob.SetObjID(objID)
	return blob
}

func genEntry(
	objID []byte,
	name string,
) internals.Entry {
	entry := internals.Entry{}

	if err := entry.New(objID, name); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: falied to create entry - %v\n",
			err,
		)
		os.Exit(1)
	}
	return entry
}

func StoreBlobs(s internals.Storer, files []os.DirEntry) []internals.Blob {
	var blobs []internals.Blob
	for _, dirEntry := range files {
		// . is the cwd and .. is the prev dir
		if dirEntry.Name() == "." || dirEntry.Name() == ".." {
			continue
		}
		if dirEntry.IsDir() {
			continue
		}
		file, err := os.Open(dirEntry.Name())
		if err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error: falied to read the file. %v\n",
				err,
			)
			os.Exit(1)
		}
		defer file.Close()

		blob := genBlob(file, dirEntry.Name())
		if err := s.Store(&blob); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error: falied to store the blob. %v\n",
				err,
			)
			os.Exit(1)
		}
		blobs = append(blobs, blob)
	}
	return blobs
}

func genEntries(
	blobs []internals.Blob,
) []internals.Entry {
	var entries []internals.Entry
	for _, blob := range blobs {
		entry := genEntry(blob.GetObjID(), blob.GetName())
		entries = append(entries, entry)
	}
	return entries
}

func StoreTree(s internals.Storer, entries []internals.Entry) internals.Tree {
	tree := internals.Tree{}
	if err := tree.New(entries); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: falied to create a new tree - %v\n",
			err,
		)
		os.Exit(1)
	}
	if err := s.Store(&tree); err != nil {
		fmt.Fprintf(
			os.Stderr,
			"Error: falied to store the tree - %v\n",
			err,
		)
		os.Exit(1)
	}
	return tree
}

func genAuthor() internals.Author {
	config := ConfigVars{}
	config.YIT_AUTHOR_EMAIL = os.Getenv("YIT_AUTHOR_EMAIL")
	config.YIT_AUTHOR_NAME = os.Getenv("YIT_AUTHOR_NAME")
	author := internals.Author{}
	if err := author.New(config.YIT_AUTHOR_NAME, config.YIT_AUTHOR_EMAIL, time.Now()); err != nil {

		fmt.Fprintf(
			os.Stderr,
			"Error: falied to create a new author %v\n",
			err,
		)
		os.Exit(1)
	}
	return author
}

func dbDir(rootDir string) string {
	yit_dir := strings.Join(
		[]string{rootDir, internals.YitMetadataDir},
		string(os.PathSeparator),
	)
	return yit_dir
}

var message string

func mustTrue(cond bool, msg string) {
	if !cond {
		fmt.Fprintf(
			os.Stderr,
			"%v\n",
			msg,
		)
		os.Exit(1)
	}
}

func mustNil(err error, msg string) {
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
			"%s - %v\n",
			msg, err,
		)
		os.Exit(1)
	}
}

func StoreCommit(
	s internals.Storer,
	parent []byte,
	tree *internals.Tree,
	msg string,
	author internals.Author,
	commiter internals.Author,
) internals.Commit {
	commit := internals.Commit{}
	err := commit.New(parent, tree.GetObjID(), msg, author, commiter)
	mustNil(err, "Error: falied to create new commit")
	err = s.Store(&commit)
	mustNil(err, "Error: falied to store the commit")
	return commit
}

func init() {
	rootCmd.AddCommand(commitCmd)
	commitCmd.Flags().
		StringVarP(&message, "message", "m", "", "Use the given <msg> as the commit message.")
}
