package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

var initConfig = &cobra.Command{
	Use:   "init",
	Short: "init a config file",
	Run: func(cmd *cobra.Command, args []string) {
		err := viper.SafeWriteConfig()
		if err != nil {
			panic(err)
		}
	},
}

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "patient service graphQL server",
}

var nodes = map[string]string{}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {

		name := viper.GetString("NAME")

		if name == "" {
			panic("instance name must be set")
		}
		clusterURL := viper.GetString("CLUSTER")
		if clusterURL == "" {
			panic("cluster url must be set")
		}

		err := register(name, clusterURL)
		if err != nil {
			log.Error().Str("clusterURL", clusterURL).Str("name", name).Msg("failed to connect to existing cluster")
			panic("Failed to connect to existing cluster, shutting down")
		}

		port := viper.GetString("PORT")
		log.Info().Str("Port", port).Msg("starting server...")
		r := gin.Default()
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		r.GET("/brokers", func(c *gin.Context) {
			c.JSON(200, nodes)
		})
		r.POST("/brokers", func(c *gin.Context) {
			var request registerRequest
			err := c.ShouldBindJSON(&request)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			nodes[request.Name] = request.URL
			c.JSON(200, nodes)
		})
		r.Run(fmt.Sprintf(":%s", port))
	},
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "5000")
	viper.SetDefault("DEBUG", true)

	err := viper.ReadInConfig()
	if err != nil {
		// if config file is not found, it's ok because we can use env variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Panic().Msgf("error reading config %v", err)
			panic("error reading config")
		}
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if viper.GetBool("DEBUG") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	rootCmd.AddCommand(initConfig)
	rootCmd.AddCommand(serverCmd)
	rootCmd.Execute()

}
