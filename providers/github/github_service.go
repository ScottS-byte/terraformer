// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"context"
	"net/http"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v35/github"
	"golang.org/x/oauth2"
)

const githubDefaultURL = "https://api.github.com/"

type GithubService struct { //nolint
	terraformutils.Service
}

func (g *GithubService) createClient() (*github.Client, error) {
	if g.GetArgs()["base_url"].(string) == githubDefaultURL {
		return g.createRegularClient(), nil
	}
	return g.createEnterpriseClient()
}

func (g *GithubService) createRegularClient() *github.Client {
	ctx := context.Background()
	if g.Args["app_id"].(int64) != 0 && g.Args["installation_id"].(int64) != 0 && g.Args["pem"].(string) != "" {
		itr, err := ghinstallation.New(http.DefaultTransport, g.Args["app_id"].(int64), g.Args["installation_id"].(int64), []byte(g.Args["pem"].(string)))
		if err != nil {
			return nil
		}
		return github.NewClient(&http.Client{Transport: itr})
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Args["token"].(string)},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func (g *GithubService) createEnterpriseClient() (*github.Client, error) {
	ctx := context.Background()
	baseURL := g.GetArgs()["base_url"].(string)
	if g.Args["app_id"].(int64) != 0 && g.Args["installation_id"].(int64) != 0 && g.Args["pem"].(string) != "" {
		itr, err := ghinstallation.New(http.DefaultTransport, g.Args["app_id"].(int64), g.Args["installation_id"].(int64), []byte(g.Args["pem"].(string)))
		if err != nil {
			return nil, err
		}
		return github.NewEnterpriseClient(baseURL, baseURL, &http.Client{Transport: itr})
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Args["token"].(string)},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewEnterpriseClient(baseURL, baseURL, tc)
}
