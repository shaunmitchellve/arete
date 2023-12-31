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

package cmdsolution

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"errors"
	"regexp"

	"github.com/shaunmitchellve/arete/pkg/utils"
	solutionFilev1 "github.com/shaunmitchellve/arete/pkg/api/solution/v1"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// solutions List struct stores the YAML list of solutions from either
// the GitHub core solutions, local cached copy and / or the merged version of both
type SolutionsList struct {
	Solutions []Solution `yaml:"solutions"`
}

type Solution struct {
	Solution string `yaml:"solution"`
	Description string `yaml:"description"`
	Url string `yaml:"url"`
}

// String return a simple string representation of the Solutions maps
func (s SolutionsList) String() string {
	return fmt.Sprintf("%s", s.Solutions)
}

// Compare 2 solutionlists and return a combined list of unique solutions
func (firstSL *SolutionsList) compareSolutions (secondSL *SolutionsList) error {

	var found bool

	if len(firstSL.Solutions) == 0 && len(secondSL.Solutions) > 0 {
		firstSL.Solutions = secondSL.Solutions

		return nil
	}

	for _, sSolution := range secondSL.Solutions {
		found = false
		for _, fSolution := range firstSL.Solutions {
			if reflect.DeepEqual(fSolution, sSolution) {
				found = true
			}
		}

		if !found {
			firstSL.Solutions = append(firstSL.Solutions, sSolution)
		}
	}

	return nil
}

// Download the solution from it's repo and update the local solution.yaml file
func (sl *SolutionsList) GetSolution(url string, branch string, subFolder string) error {
	 if err := sl.GetRemoteSolution(url, branch, subFolder); err != nil {
		return err
	 }

	 cacheDir := filepath.Join(viper.GetString("cache"),  sl.Solutions[0].Solution)
	 _, statErr := os.Stat(cacheDir)
	 if statErr == nil {
		os.RemoveAll(cacheDir)
	}

	log.Info().Msg("Pulling package from repo...")

	// Use KPT to pull down the package
	resp, err := utils.CallCommand(utils.Kpt, []string{"pkg", "get", sl.Solutions[0].Url, cacheDir}, false)

	if err != nil {
		log.Error().Err(err).Msg(string(resp))
		return err
	}

	if viper.GetBool("verbose") {
		log.Debug().Msg(string(resp))
	}

	log.Info().Msg("Solution downloaded to " + cacheDir)

	return nil
}

// Get a GitHub raw file from the url, branch and subFolder provided
func getGitHubRaw(url string, branch string, subFolder string, file string) (string, error) {
	var lines []string
	ret := ""

	if branch == "" {
		branch = "main"
	}

	// Remove prefix and suffix forward slashes
	if subFolder == "/" {
		subFolder = "base"
	} else if strings.Index(subFolder, "/") == 0 {
		subFolder = strings.Replace(subFolder, "/", "", 1)
	}

	subFolder = strings.TrimSuffix(subFolder, "/")

	reg := regexp.MustCompile(`^https://github\.com/([a-zA-Z0-9-/]*)`)
	res := reg.FindStringSubmatch(url)

	// If the repo is private then a token can be passed in the URL to get access to the solutions file
	gitToken := viper.GetString("git_token")

	if len(res) == 2 {
		url = "https://"

		if  gitToken != "" {
			url = url + gitToken + "@"
		}

		url = url + "raw.githubusercontent.com/" + res[1] + "/" + branch + "/" + subFolder + "/" + file
	} else {
		return ret, errors.New("malformed URL("+ url +"), unable to process")
	}

	if viper.GetBool("verbose") {
		log.Debug().Msgf("Getting %s from url: %s", file, url)
	}

	resp, err := http.Get(url)

	if err != nil || resp.StatusCode == 404 {
		return ret, err
	}

	if viper.GetBool("verbose") {
		log.Debug().Msg(resp.Status)
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Err(err).Msg("")
	}

	ret = strings.Join(lines, "\n")

	return ret, nil
}

// Get the solutions.yaml file from GitHub and store in the local cache.
func (sl *SolutionsList) GetCoreSolutions() error {
	var (
		err error
		ret string
		cachedSolutions SolutionsList
		yamlout []byte
	)

	if ret, err = getGitHubRaw(viper.GetString("repoUrl"), viper.GetString("repoBranch"), viper.GetString("repoSubFolder"), "solutions.yaml"); err != nil {
		return err
	}

	if err = yaml.Unmarshal([]byte(ret), &sl); err != nil {
		return err
	}

	if err = cachedSolutions.GetCacheSolutions(); err != nil {
		return err
	}

	if err = sl.compareSolutions(&cachedSolutions); err != nil {
		return err
	}

	if yamlout, err = yaml.Marshal(&sl); err != nil {
		return err
	}

	ret = string(yamlout)

	if err = utils.WriteToCache(&ret, "solutions.yaml", false); err != nil {
		return err
	}

	return nil
}

// Get the solution.yaml file from a GitHub repo.
func (sl *SolutionsList) GetRemoteSolution(url string, branch string, subFolder string) error {
	var (
		err error
		ret string
		cachedSolutions SolutionsList
		yamlout []byte
	)

	if ret, err = getGitHubRaw(url, branch, subFolder, "solution.yaml"); err != nil {
		return err
	}

	solutionFile := solutionFilev1.SolutionFile{}

	if err = yaml.Unmarshal([]byte(ret), &solutionFile); err != nil {
		return err
	}

	if !solutionFile.Spec.IsEmpty() {
		sol := make([]Solution, 1)
		sol[0].Url = solutionFile.Spec.Url
		sol[0].Description = solutionFile.Spec.Description
		sol[0].Solution = solutionFile.Name

		sl.Solutions = append(sl.Solutions, sol[0])

		if err := cachedSolutions.GetCacheSolutions(); err != nil {
			return err
		}

		if err = sl.compareSolutions(&cachedSolutions); err != nil {
			return err
		}

		if yamlout, err = yaml.Marshal(&sl); err != nil {
			return err
		}

		ret = string(yamlout)

	 	if err = utils.WriteToCache(&ret, "solutions.yaml", false); err != nil {
			return err
		}
	}

	return nil
}

// Get the cached solutions.yaml file
func (sl *SolutionsList) GetCacheSolutions() error {
	solutionsFile := filepath.Join(viper.GetString("cache"), "solutions.yaml")
	cacheSl, err := os.ReadFile(solutionsFile)

	if err != nil {
		if _, err = os.Create(solutionsFile); err != nil {
			return err
		}
	}

	if err = yaml.Unmarshal(cacheSl, &sl); err != nil {
		return err
	}

	return nil
}

// GetUrl will search the solution list for the passed in solution and return
// the URL or error if the solution is not found
func (sl *SolutionsList) GetUrl(solutionName string) (string, error) {
	for _, solution := range sl.Solutions {
		if solution.Solution == solutionName {
			return solution.Url, nil
		}
	}

	return "", errors.New("solution not found")
}