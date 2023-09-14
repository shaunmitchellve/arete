// Copyright 2023 Shaun Mitchell

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
	"strings"
	"path/filepath"
)

// init the config directory and file
func Init() {
	usrH, _ := os.UserHomeDir();

	configPath := filepath.Join(usrH, ".arete")
	configFile := "config.yaml"
	solutionFile := "solutions.yaml"
	repoURL := "https://github.com/shaunmitchellve/arete"
	repoBranch := "main"
	repoSubFolder := "solutions"

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	viper.SetConfigName(configFile[:strings.Index(configFile, ".")])
	viper.SetConfigType(configFile[strings.Index(configFile, ".")+1:])
	viper.AddConfigPath(configPath)
	viper.SetDefault("cache", configPath)
	viper.SetDefault("repoUrl", repoURL)
	viper.SetDefault("repoBranch", repoBranch)
	viper.SetDefault("repoSubFolder", repoSubFolder)

	if _, cpErr := os.Stat(configPath); os.IsNotExist(cpErr) {
		if mkErr := os.Mkdir(configPath, 0744); mkErr != nil {
			log.Fatal().Err(mkErr).Msg("")
		}
	}

	if _, fErr := os.Stat(filepath.Join(configPath, configFile)); os.IsNotExist(fErr) {
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal().Err(err).Msg("Unable to create config file")
		}
	}

	solutionFilePath := filepath.Join(configPath, solutionFile)
	if _, fErr := os.Stat(solutionFilePath); os.IsNotExist(fErr) {
		if _, err := os.Create(solutionFilePath); err != nil {
			log.Fatal().Err(err).Msg("Unable to create cached solution file")
		}
	}

	if cErr := viper.ReadInConfig(); cErr != nil {
		log.Error().Err(cErr).Msg("Unable to read config file")
	}
}