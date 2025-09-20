package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"

	"github.com/VieiraGabrielAlexandre/aliasgen/internal/learn"
	"github.com/VieiraGabrielAlexandre/aliasgen/internal/shell"
	"github.com/VieiraGabrielAlexandre/aliasgen/internal/store"
	"github.com/VieiraGabrielAlexandre/aliasgen/pkg/aliasfmt"
)

func main() {
	root := &cobra.Command{
		Use:   "aliasgen",
		Short: "Gerador automático de aliases para Bash/Zsh/Fish",
		Long:  "Observa seu histórico, sugere aliases e aplica ao seu shell.",
	}

	root.AddCommand(cmdLearn(), cmdSuggest(), cmdApply(), cmdList(), cmdExplain())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func cmdLearn() *cobra.Command {
	return &cobra.Command{
		Use:   "learn",
		Short: "Ingere o histórico e atualiza o banco",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := learn.Run()
			if err != nil {
				return err
			}
			fmt.Printf("ingestão concluída: %d registros\n", n)
			return nil
		},
	}
}

func cmdSuggest() *cobra.Command {
	var limit int
	c := &cobra.Command{
		Use:   "suggest",
		Short: "Mostra sugestões de aliases com score",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := store.Open()
			if err != nil {
				return err
			}
			defer db.Close()

			stats, err := db.TopCommands(limit)
			if err != nil {
				return err
			}
			sugs := learn.GenerateSuggestions(stats, time.Now())

			sort.Slice(sugs, func(i, j int) bool { return sugs[i].Score > sugs[j].Score })
			for _, s := range sugs {
				fmt.Printf("%-12s -> %s   (score=%.2f)  [%s]\n", s.Alias, s.Command, s.Score, s.Reason)
			}
			if len(sugs) == 0 {
				fmt.Println("sem sugestões por enquanto. rode `aliasgen learn` após usar mais o terminal.")
			}
			return nil
		},
	}
	c.Flags().IntVar(&limit, "limit", 200, "limite de comandos a considerar")
	return c
}

func cmdApply() *cobra.Command {
	var sh string
	var top int
	c := &cobra.Command{
		Use:   "apply",
		Short: "Gera/escreve arquivo de aliases e dá instruções para fazer source",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := store.Open()
			if err != nil {
				return err
			}
			defer db.Close()

			if sh == "auto" || sh == "" {
				sh = shell.Detect()
				if sh == "unknown" {
					fmt.Println("não consegui detectar seu shell. Use: --shell bash|zsh|fish")
					return nil
				}
			}

			stats, err := db.TopCommands(500)
			if err != nil {
				return err
			}
			sugs := learn.GenerateSuggestions(stats, time.Now())
			sort.Slice(sugs, func(i, j int) bool { return sugs[i].Score > sugs[j].Score })
			if top > 0 && top < len(sugs) {
				sugs = sugs[:top]
			}

			var lines []string
			for _, s := range sugs {
				var line string
				switch sh {
				case "fish":
					line = aliasfmt.Fish(s.Alias, s.Command)
				default:
					line = aliasfmt.Bash(s.Alias, s.Command) // serve para bash e zsh
				}
				lines = append(lines, line)
				_ = db.SaveAlias(sh, s.Alias, s.Command, s.Score)
			}

			path, err := shell.WriteAliases(sh, lines)
			if err != nil {
				return err
			}
			fmt.Printf("aliases gravados em: %s\n\n%s\n", path, shell.SourceHint(sh))
			return nil
		},
	}
	c.Flags().StringVar(&sh, "shell", "auto", "bash|zsh|fish|auto")
	c.Flags().IntVar(&top, "top", 30, "quantidade de aliases a aplicar (0 = todos)")
	return c
}

func cmdList() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lista aliases atuais do banco",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := store.Open()
			if err != nil {
				return err
			}
			defer db.Close()
			rows, err := db.ListAliases()
			if err != nil {
				return err
			}
			for _, r := range rows {
				fmt.Printf("%-12s -> %-40s  [shell=%s score=%.2f]\n", r.Alias, r.ExpandsTo, r.Shell, r.Score)
			}
			if len(rows) == 0 {
				fmt.Println("sem aliases aplicados ainda. rode `aliasgen apply`.")
			}
			return nil
		},
	}
}

func cmdExplain() *cobra.Command {
	return &cobra.Command{
		Use:   "explain <alias>",
		Short: "Explica por que uma sugestão existe (MVP: mostra linha atual do banco)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			db, err := store.Open()
			if err != nil {
				return err
			}
			defer db.Close()

			rows, err := db.ListAliases()
			if err != nil {
				return err
			}
			for _, r := range rows {
				if r.Alias == target {
					fmt.Printf("%s -> %s (score=%.2f, shell=%s, criado=%s)\n", r.Alias, r.ExpandsTo, r.Score, r.Shell, r.CreatedAt.Format(time.RFC3339))
					return nil
				}
			}
			fmt.Println("alias não encontrado. rode `aliasgen suggest` e `aliasgen apply` antes.")
			return nil
		},
	}
}
