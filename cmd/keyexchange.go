package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nikolatw/sshkeyexchange-go/pkg/keygen"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pass := strings.TrimSpace(os.Args[1])

	keys, err := keygen.NewWithPasscode(pass)
	handleError(err)

	fmt.Println(string(keys.Private))
	fmt.Println(string(keys.Public))

	shell, err := interp.New(interp.StdIO(os.Stdin, os.Stdout, os.Stderr))
	handleError(err)

	sshkey, err := keys.SSHPublicKey()
	handleError(err)

	cmd := fmt.Sprintf("echo \"%s\" >> .ssh/authorized_keys", sshkey)
	err = run(shell, cmd, "add to .ssh/authorized_keys")
	handleError(err)
}

func run(r *interp.Runner, command string, name string) error {
	prog, err := syntax.NewParser().Parse(strings.NewReader(command), name)
	if err != nil {
		return err
	}
	r.Reset()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	return r.Run(ctx, prog)
}
