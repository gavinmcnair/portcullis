package main

import (
	"crypto/rand"
	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/big"
	"net"
	"time"
)

type config struct {
	StartDelay   time.Duration `env:"DELAY_SECONDS" envDefault:"0"`
	Hostname     string        `env:"HOSTNAME,required"`
	InitialDelay time.Duration `env:"INITIAL_DELAY" envDefault:"0s"`
	Port         string        `env:"PORT,required"`
	RandomWindow int64         `env:"RANDOM_WINDOW" envDefault:"0"`
	SleepCount   time.Duration `env:"SLEEP_COUNT" envDefault:"5s"`
	Timeout      time.Duration `env:"TIMEOUT" envDefault:"1s"`
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
	log.Info().Msg("ending portcullis - gracefully exiting")

}

func run() error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return err
	}

	log.Info().
		Dur("start_delay", cfg.StartDelay).
		Str("hostname", cfg.Hostname).
		Str("ports", cfg.Port).
		Dur("initial_delay", cfg.InitialDelay).
		Int64("random_window", cfg.RandomWindow).
		Dur("sleep_count", cfg.SleepCount).
		Dur("timeout", cfg.Timeout).
		Msg("starting portcullis")

	if cfg.InitialDelay.Seconds() > 0 {
		log.Info().
			Dur("initial_delay", cfg.InitialDelay).
			Msg("Waiting for initial_delay")
		time.Sleep(cfg.StartDelay)
	}

	failureCount := 0

	for {

		success, err := raw_connect(cfg.Hostname, cfg.Timeout, cfg.Port)

		if !success {
			failureCount++
			log.Info().
				Int("failure_count", failureCount).
				Dur("sleep_count", cfg.SleepCount).
				Err(err).
				Msg("port not responding. Will retry after delay")
			time.Sleep(cfg.SleepCount)
			continue
		}

		log.Info().Msg("target port is responding")

		if cfg.RandomWindow != 0 && failureCount > 1 {
			err := waitFor(cfg.RandomWindow)
			if err != nil {
				return err
			}
		}

		if cfg.StartDelay.Seconds() > 0 {
			log.Info().
				Dur("start_delay", cfg.StartDelay).
				Msg("delaying start of containers")
			time.Sleep(cfg.StartDelay)
		}

		return nil
	}

}

func raw_connect(host string, timeout time.Duration, port string) (bool, error) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)

	if conn != nil {
		defer conn.Close()
		return true, err
	}

	return false, err
}

func waitFor(randomWindow int64) error {
	randomSeconds, err := genRandNum(randomWindow)
	if err != nil {
		return err
	}

	log.Info().
		Int64("wait_duration_seconds", randomSeconds).
		Int64("max_duration_seconds", randomWindow).
		Msg("sleeping for random period")

	time.Sleep(time.Duration(randomSeconds) * time.Second)
	return nil
}

func genRandNum(max int64) (int64, error) {
	bg := big.NewInt(max)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		return 0, err
	}

	return n.Int64(), err
}
