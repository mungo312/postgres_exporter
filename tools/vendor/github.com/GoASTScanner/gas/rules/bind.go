// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
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

package rules

import (
	"go/ast"
	"regexp"

	"github.com/GoASTScanner/gas"
)

// Looks for net.Listen("0.0.0.0") or net.Listen(":8080")
type bindsToAllNetworkInterfaces struct {
	gas.MetaData
	calls   gas.CallList
	pattern *regexp.Regexp
}

func (r *bindsToAllNetworkInterfaces) Match(n ast.Node, c *gas.Context) (*gas.Issue, error) {
	callExpr := r.calls.ContainsCallExpr(n, c)
	if callExpr == nil {
		return nil, nil
	}
	if arg, err := gas.GetString(callExpr.Args[1]); err == nil {
		if r.pattern.MatchString(arg) {
			return gas.NewIssue(c, n, r.What, r.Severity, r.Confidence), nil
		}
	}
	return nil, nil
}

// NewBindsToAllNetworkInterfaces detects socket connections that are setup to
// listen on all network interfaces.
func NewBindsToAllNetworkInterfaces(conf gas.Config) (gas.Rule, []ast.Node) {
	calls := gas.NewCallList()
	calls.Add("net", "Listen")
	calls.Add("crypto/tls", "Listen")
	return &bindsToAllNetworkInterfaces{
		calls:   calls,
		pattern: regexp.MustCompile(`^(0.0.0.0|:).*$`),
		MetaData: gas.MetaData{
			Severity:   gas.Medium,
			Confidence: gas.High,
			What:       "Binds to all network interfaces",
		},
	}, []ast.Node{(*ast.CallExpr)(nil)}
}
