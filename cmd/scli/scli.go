/*
Copyright © 2022

Licensed under the  License, Version 20221227 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/shammishailaj/scli/internal/app/scli/cmd"
)

var (
	BuildDate  string
	GitBranch  string
	GitCommit  string
	GitState   string
	GitSummary string
	Version    string
)

func main() {
	cmd.BuildDate = BuildDate
	cmd.GitBranch = GitBranch
	cmd.GitCommit = GitCommit
	cmd.GitState = GitState
	cmd.GitSummary = GitSummary
	cmd.Version = Version
	cmd.Execute()
}
