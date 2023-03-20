package commands

import (
	"io/ioutil"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/bogatyr285/sumcalc-testtask/config"
	"github.com/bogatyr285/sumcalc-testtask/internal/api"
	"github.com/bogatyr285/sumcalc-testtask/internal/http"
	jwtService "github.com/bogatyr285/sumcalc-testtask/internal/services/jwt"
	sha256Service "github.com/bogatyr285/sumcalc-testtask/internal/services/sha256"
	sumService "github.com/bogatyr285/sumcalc-testtask/internal/services/sum"
)

func NewServeCmd() *cobra.Command {
	var configPath string
	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// todo we can use DI system like wire/dig/fx if our app become bigger, for now & simplicity reasons assembly manually
			logger, _ := zap.NewProduction()
			defer logger.Sync() // flushes buffer, if any
			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()

			logger.Info("starting")
			var conf config.Config
			err := config.ReadYaml(configPath, &conf)
			if err != nil {
				return err
			}
			logger.Info("config parsed", zap.Any("conf", conf))
			if conf.Env == config.EnvProduction {
				gin.SetMode(gin.ReleaseMode)
			}

			// todo we can do special `config file parser`(embeded struct) and read files there
			pk, err := ioutil.ReadFile(conf.JWT.PrivateKey)
			if err != nil {
				return err
			}
			pubk, err := ioutil.ReadFile(conf.JWT.PublicKey)
			if err != nil {
				return err
			}

			jwtManager, err := jwtService.NewJWTManager(conf.JWT.Issuer, conf.JWT.ExpiresIn, pubk, pk, jwa.ES256)
			if err != nil {
				return err
			}

			sumCalculator := sumService.NewSumLogger(sumService.NewSumCalculator(), logger)
			sha256Hasher := sha256Service.NewSHA256Hasher()

			ginEngine := api.NewRouter(jwtManager, jwtManager, sumCalculator, sha256Hasher, logger)
			httpServer, err := http.NewHTTPServer(conf.Listen.HTTP, ginEngine, logger)
			if err != nil {
				return err
			}

			httpStop, err := httpServer.Start()
			if err != nil {
				return err
			}

			logger.Info("started")
			<-ctx.Done()
			if err := httpStop(); err != nil {
				logger.Warn("closing http", zap.Error(err))
			}

			logger.Info("finished")

			return nil
		},
	}
	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
