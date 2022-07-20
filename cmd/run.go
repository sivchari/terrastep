package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
	yml "gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatalf("Failed to get config path err = %s\n", err.Error())
		}
		if err := run(context.Background(), c); err != nil {
			log.Fatalf("Failed to run err = %s\n", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("config", "c", "", "config path")
}

func run(ctx context.Context, c string) error {
	cfg, err := parseyml(c)
	if err != nil {
		return err
	}
	if err := runtf(ctx, cfg); err != nil {
		return err
	}
	return nil
}

func parseyml(c string) (*Config, error) {
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	f, err := os.Open(p + "/" + c)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := make([]byte, 1024)
	l, err := f.Read(b)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yml.Unmarshal(b[:l], &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func runtf(ctx context.Context, cfg *Config) error {
	eg, egctx := errgroup.WithContext(ctx)
	p, err := os.Getwd()
	if err != nil {
		return err
	}
	for _, t := range cfg.Tasks {
		t := t
		eg.Go(func() error {
			for _, s := range t.Steps {
				if out, err := exec.CommandContext(egctx, "sh", "-c", fmt.Sprintf("cd %s && terraform init", p+"/"+s)).Output(); err != nil {
					return fmt.Errorf("falied to exec terraform init err = %w out = %s", err, string(out))
				}
				for _, tac := range t.Tactics {
					out, err := exec.CommandContext(egctx, "sh", "-c", fmt.Sprintf("cd %s && terraform %s", p+"/"+s, tac)).Output()
					if err != nil {
						return fmt.Errorf("failed to exec terraform %s err = %w out = %s", tac, err, string(out))
					}
					if tac == "plan" || strings.HasSuffix(tac, "apply") {
						fmt.Fprint(os.Stdout, string(out))
					}
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Tasks []*Task `yaml:"tasks"`
}

type Task struct {
	Name    string   `yaml:"name"`
	Tactics []string `yaml:"tactics"`
	Steps   []string `yaml:"steps"`
}
