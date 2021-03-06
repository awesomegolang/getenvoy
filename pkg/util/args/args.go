// Copyright 2020 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package args

import (
	"github.com/pkg/errors"

	"github.com/mattn/go-shellwords"
)

// SplitCommandLine splits fragments of a command line down to individual arguments.
func SplitCommandLine(fragments ...string) ([]string, error) {
	args := make([]string, 0, len(fragments))
	for _, fragment := range fragments {
		words, err := shellwords.Parse(fragment)
		if err != nil {
			return nil, errors.Errorf("%q is not a valid command line string", fragment)
		}
		args = append(args, words...)
	}
	return args, nil
}
