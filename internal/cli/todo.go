package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nios-V/wl/internal/render"
	"github.com/Nios-V/wl/internal/store"
	"github.com/spf13/cobra"
)

func newTodoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "todo [texto]",
		Short: "Agrega un pendiente; sin texto, lista los abiertos",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			st, err := openStore()
			if err != nil {
				return err
			}
			defer st.Close()

			if len(args) == 0 {
				todos, err := st.OpenTodos()
				if err != nil {
					return err
				}
				if len(todos) == 0 {
					fmt.Fprintln(out, "Sin pendientes.")
				}
				render.Todos(out, todos)
				return nil
			}

			text := strings.Join(strings.Fields(strings.Join(args, " ")), " ")
			id, err := st.AddTodo(text, time.Now())
			if err != nil {
				return err
			}
			fmt.Fprintf(out, "Pendiente agregado (id %d)\n", id)
			return nil
		},
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "done <id>",
			Short: "Marca un pendiente como completado",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return closeTodo(cmd, args[0], true)
			},
		},
		&cobra.Command{
			Use:   "drop <id>",
			Short: "Elimina un pendiente",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return closeTodo(cmd, args[0], false)
			},
		},
	)
	return cmd
}

func closeTodo(cmd *cobra.Command, arg string, done bool) error {
	id, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return fmt.Errorf("id inválido: %q", arg)
	}
	st, err := openStore()
	if err != nil {
		return err
	}
	defer st.Close()

	msg := "Completado."
	if done {
		err = st.CompleteTodo(id, time.Now())
	} else {
		err = st.DeleteTodo(id)
		msg = "Eliminado."
	}
	if errors.Is(err, store.ErrNotFound) {
		return fmt.Errorf("no hay un pendiente abierto con id %d", id)
	}
	if err != nil {
		return err
	}
	fmt.Fprintln(cmd.OutOrStdout(), msg)
	return nil
}
