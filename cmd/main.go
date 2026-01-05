package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"mexa/internal/config"
	mexadomain "mexa/internal/domains/mexa"
	tg "mexa/internal/infra/bot/telegram"
	"mexa/internal/infra/db/sqlite"
	mexasqlite "mexa/internal/infra/db/sqlite/mexa"
	"mexa/internal/infra/fsm/memory"
	mexaservice "mexa/internal/services/mexa"
	"mexa/internal/worker"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	configFPathFlag = flag.String("config_file", filepath.Join(
		"config.json",
	), "path to config file")
)

const (
	exitCodeOk = iota
	exitCodeErrRun
	exitCodeErrParseConfigFile
	exitCodeErrInvalidConfig
	exitCodeErrSetupDb
	exitCodeErrParseCasesDir
	exitCodeErrNoCases
	exitCodeErrClearCases
	exitCodeErrAddCase
	exitCodeErrCreateBot
	exitCodeErrCreateService
)

func main() {
	flag.Parse()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error
	c, err := parseConfigFile(*configFPathFlag)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		os.Exit(exitCodeErrParseConfigFile)
	}

	err = c.Validate()
	if err != nil {
		fmt.Println("Invalid config:", err)
		os.Exit(exitCodeErrInvalidConfig)
	}

	fmt.Println("Setting up database at:", c.DbPath)
	db, err := setupDb(ctx, *c)
	if err != nil {
		fmt.Println("Failed to create database:", err)
		os.Exit(exitCodeErrSetupDb)
	}
	defer db.Close()

	cases, err := parseCasesDir(*c.CasesDir)
	if err != nil {
		fmt.Println("Failed to parse cases dir:", err)
		os.Exit(exitCodeErrParseCasesDir)
	} else if len(cases) == 0 {
		fmt.Println("No cases found in dir:", *c.CasesDir)
		os.Exit(exitCodeErrNoCases)
	}

	{
		base := sqlite.NewBaseRepo(db)
		repo := mexasqlite.NewCasesRepo(&base)

		err = repo.ClearCases(ctx, c.Exercise.Id)
		if err != nil {
			fmt.Println("Failed to clear cases:", err)
			os.Exit(exitCodeErrClearCases)
		}

		for _, cs := range cases {
			_, err = repo.AddCase(ctx, c.Exercise.Id, cs)
			if err != nil {
				fmt.Println("Failed to insert case:", err)
				os.Exit(exitCodeErrAddCase)
			}
		}
	}

	var bot *tg.Bot
	{
		bot, err = tg.New(ctx, *c.Token, *c.Batch, *c.Exercise, c.AdminIds)
		if err != nil {
			fmt.Println("Failed to create bot:", err)
			os.Exit(exitCodeErrCreateBot)
		}
	}

	var ser *mexaservice.Service
	{
		base := sqlite.NewBaseRepo(db)
		repos := mexaservice.Repos{
			Transactional: sqlite.NewTransactional(db.Db()),
			Users:         mexasqlite.NewUsersRepo(&base),
			Cases:         mexasqlite.NewCasesRepo(&base),
			Casualties:    mexasqlite.NewCasualtiesRepo(&base),
			Exercises:     mexasqlite.NewExercisesRepo(&base),
			Deterioration: mexasqlite.NewCasualtiesDeteriorationRepo(&base),
			ExLogs:        mexasqlite.NewExLogsRepo(&base),
			CCLogs:        mexasqlite.NewCadetCaseLogsRepo(&base),
		}

		ser, err = mexaservice.NewService(ctx, mexaservice.ServiceConfig{
			Bot:      bot,
			Repos:    repos,
			Exercise: *c.Exercise,
			Batch:    *c.Batch,
			Admins:   c.AdminIds,
			Fsm:      memory.NewFsm(),
		})
		if err != nil {
			fmt.Println("Failed to create service:", err)
			os.Exit(exitCodeErrCreateService)
		}
	}

	errCh := make(chan error, 1)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		w := worker.New(bot, *ser)
		errCh <- w.Start(ctx)
	})

	select {
	case sig := <-sigCh:
		fmt.Println("Received os signal to terminate, shutting down services, signal:", sig)
	case err = <-errCh:
		if err != nil {
			fmt.Println("Received error from a service, shutting down all services, error", err)
		}
		cancel()
		wg.Wait()

		if err != nil {
			os.Exit(exitCodeErrRun)
		}

		os.Exit(exitCodeOk)
	}
}

func parseConfigFile(fp string) (c *config.Config, err error) {
	b, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	c = &config.Config{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	fmt.Println("Loaded config:", c)

	return c, nil
}

func parseCasesDir(dir string) (res []mexadomain.CaseValue, err error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fmt.Println("Parsing cases dir:", abs)
	fps, err := filepath.Glob(filepath.Join(abs, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, fp := range fps {
		c, err := parseCaseFile(fp)
		if err != nil {
			return nil, err
		}

		res = append(res, *c)
	}

	return res, nil
}

func setupDb(ctx context.Context, c config.Config) (db *sqlite.DB, err error) {
	db, err = sqlite.New(*c.DbPath)
	if err != nil {
		return nil, err
	}

	err = db.Init(ctx)
	if err != nil {
		return nil, err
	}
	{
		base := sqlite.NewBaseRepo(db)
		{
			ex := c.Exercise

			repo := mexasqlite.NewExercisesRepo(&base)
			id, err := repo.AddExercise(ctx, ex.Code, ex.Name)
			if err != nil {
				return nil, err
			}

			c.Exercise.Id = *id
		}
	}

	return db, nil
}

func parseCaseFile(fp string) (res *mexadomain.CaseValue, err error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res = &mexadomain.CaseValue{}
	err = json.NewDecoder(f).Decode(&res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
