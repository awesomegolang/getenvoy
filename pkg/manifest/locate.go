// Copyright 2019 Tetrate
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

package manifest

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// NewKey creates a manifest key based on the reference it is given
func NewKey(reference string) (*Key, error) {
	r := regexp.MustCompile(`^(.+):(.+)/(.+)$`)
	matches := r.FindStringSubmatch(reference)
	if len(matches) != 4 {
		return nil, fmt.Errorf("reference %v is not of valid format <flavor>:<version>/<platform>", reference)
	}

	return &Key{strings.ToLower(matches[1]), strings.ToLower(matches[2]), platformToEnum(matches[3])}, nil
}

func platformToEnum(s string) string {
	s = strings.ToUpper(s)
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

// Key is the primary key used to locate Envoy builds in the manifest
type Key struct {
	Flavor   string
	Version  string
	Platform string
}

// Locate returns the location of the binary for the passed parameters from the passed manifest
// The build version is searched for as a prefix of the OperatingSystemVersion.
// If the OperatingSystemVersion is empty it returns the first build listed for that operating system
func Locate(key *Key, manifestLocation string) (string, error) {
	if u, err := url.Parse(manifestLocation); err != nil || u.Host == "" || u.Scheme == "" {
		return "", errors.New("only URL manifest locations are supported")
	}
	manifest, err := fetch(manifestLocation)
	if err != nil {
		return "", err
	}

	// This is pretty horrible... Not sure there is a nicer way though.
	if manifest.Flavors[key.Flavor] != nil && manifest.Flavors[key.Flavor].Versions[key.Version] != nil {
		for _, build := range manifest.Flavors[key.Flavor].Versions[key.Version].Builds {
			if strings.EqualFold(build.Platform.String(), key.Platform) {
				return build.DownloadLocationUrl, nil
			}
		}
	}
	return "", fmt.Errorf("unable to find matching build for %v", key)
}
