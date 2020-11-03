// This file is part of arduino-cli.
//
// Copyright 2020 ARDUINO SA (http://www.arduino.cc/)
//
// This software is released under the GNU General Public License version 3,
// which covers the main part of arduino-cli.
// The terms of this license can be found at:
// https://www.gnu.org/licenses/gpl-3.0.en.html
//
// You can be released from the requirements of the above licenses by purchasing
// a commercial license. Buying such a license is mandatory if you want to
// modify or otherwise use the software for commercial activities involving the
// Arduino software without disclosing the source code of your own applications.
// To purchase a commercial license, send an email to license@arduino.cc.

package builder

import (
	"os"
	"sort"
	"strings"

	"github.com/arduino/arduino-cli/legacy/builder/builder_utils"
	"github.com/arduino/arduino-cli/legacy/builder/constants"
	"github.com/arduino/arduino-cli/legacy/builder/types"
	"github.com/arduino/arduino-cli/legacy/builder/utils"
	properties "github.com/arduino/go-properties-orderedmap"
	"github.com/pkg/errors"
)

type RecipeByPrefixSuffixRunner struct {
	Prefix string
	Suffix string
}

func (s *RecipeByPrefixSuffixRunner) Run(ctx *types.Context) error {
	logger := ctx.GetLogger()
	if ctx.DebugLevel >= 10 {
		logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_LOOKING_FOR_RECIPES, s.Prefix, s.Suffix)
	}

	buildProperties := ctx.BuildProperties.Clone()
	recipes := findRecipes(buildProperties, s.Prefix, s.Suffix)

	properties := buildProperties.Clone()
	for _, recipe := range recipes {
		if ctx.DebugLevel >= 10 {
			logger.Fprintln(os.Stdout, constants.LOG_LEVEL_DEBUG, constants.MSG_RUNNING_RECIPE, recipe)
		}
		_, _, _, err := builder_utils.ExecRecipe(ctx, properties, recipe, false /* stdout */, utils.ShowIfVerbose /* stderr */, utils.Show)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil

}

func findRecipes(buildProperties *properties.Map, patternPrefix string, patternSuffix string) []string {
	var recipes []string
	for _, key := range buildProperties.Keys() {
		if strings.HasPrefix(key, patternPrefix) && strings.HasSuffix(key, patternSuffix) && buildProperties.Get(key) != "" {
			recipes = append(recipes, key)
		}
	}

	sort.Strings(recipes)

	return recipes
}
