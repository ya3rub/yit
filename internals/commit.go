package internals

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

type Commit struct {
	parentObjID []byte
	treeObjID   []byte
	ObjID       []byte
	msg         string
	author      Author
	comitter    Author
}

func ParseCommit(s Storer, commitObjID []byte) (*Commit, error) {
	commitData, err := s.ReadObj(commitObjID)
	if err != nil {
		return nil, err
	}
	dStr := string(commitData)
	i := strings.Index(dStr, " ")
	objType := dStr[0:i]
	if objType != "*internals.Commit" {
		return nil, errors.New("not a commit type")
	}

	dStartIdx := strings.Index(dStr, "\x00")
	lines := strings.Split(dStr[dStartIdx:], "\n")
	commit := &Commit{}
	commit.ObjID = commitObjID

	_, t, _ := strings.Cut(lines[0], " ")
	treeObjID, _ := hex.DecodeString(t)
	commit.treeObjID = treeObjID

	_, p, _ := strings.Cut(lines[1], " ")
	parentObjID, _ := hex.DecodeString(p)
	// fmt.Println("p", p, "hex:", parentObjID)
	commit.parentObjID = parentObjID

	_, a, _ := strings.Cut(lines[2], " ")
	// fmt.Println("A is :", a)
	author, _ := ParseAuthor(a)
	commit.author = *author

	// fmt.Println("line 3 is :", lines[3])
	_, c, _ := strings.Cut(lines[3], " ")
	// fmt.Println("C is :", c)
	comitter, _ := ParseAuthor(c)
	commit.comitter = *comitter

	commit.msg = lines[5]
	return commit, nil
}

func (c *Commit) New(
	parent, tree []byte,
	msg string,
	author, commiter Author,
) error {
	c.parentObjID = parent
	c.treeObjID = tree
	c.msg = msg
	c.author = author
	c.comitter = commiter
	return nil
}

func (c *Commit) ToString() string {
	parent := func() string {
		if c.parentObjID == nil {
			return ""
		}
		return fmt.Sprintf("%x\n", c.parentObjID)
	}
	return fmt.Sprintf(
		"tree %x\nparent %sauthor %s\ncommiter %s\n\n%s",
		c.treeObjID,
		parent(),
		c.author.ToString(),
		c.comitter.ToString(),
		c.msg,
	)
}

func (c *Commit) GetTree() []byte {
	return c.treeObjID
}

func (c *Commit) GetObjID() []byte {
	return c.ObjID
}

func (c *Commit) SetObjID(objID []byte) error {
	c.ObjID = objID
	return nil
}

func (c *Commit) GetParent() []byte {
	return c.parentObjID
}

func (c *Commit) GetMessage() string {
	return c.msg
}

func (c *Commit) GetType() string {
	return "commit"
}
