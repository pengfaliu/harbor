// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package latestk

import (
	"sort"

	"github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/lib/log"
	"github.com/goharbor/harbor/src/lib/selector"
	"github.com/goharbor/harbor/src/pkg/retention/policy/action"
	"github.com/goharbor/harbor/src/pkg/retention/policy/rule"
)

const (
	// TemplateID of latest active k rule
	TemplateID = "latestActiveK"
	// ParameterK ...
	ParameterK = TemplateID
	// DefaultK defines the default K
	DefaultK = 10
)

// evaluator for evaluating latest active k images
type evaluator struct {
	// latest k
	k int
}

// Process the candidates based on the rule definition
func (e *evaluator) Process(artifacts []*selector.Candidate) ([]*selector.Candidate, error) {
	// Sort artifacts by their "active time"
	//
	// Active time is defined as the selection of c.PulledTime or c.PushedTime,
	// whichever is bigger, aka more recent.
	sort.Slice(artifacts, func(i, j int) bool {
		return activeTime(artifacts[i]) > activeTime(artifacts[j])
	})

	i := min(e.k, len(artifacts))
	return artifacts[:i], nil
}

// Specify what action is performed to the candidates processed by this evaluator
func (e *evaluator) Action() string {
	return action.Retain
}

// New a Evaluator
func New(params rule.Parameters) rule.Evaluator {
	if params != nil {
		if p, ok := params[ParameterK]; ok {
			if v, ok := utils.ParseJSONInt(p); ok && v >= 0 {
				return &evaluator{
					k: int(v),
				}
			}
		}
	}

	log.Debugf("default parameter %d used for rule %s", DefaultK, TemplateID)

	return &evaluator{
		k: DefaultK,
	}
}

func activeTime(c *selector.Candidate) int64 {
	if c.PulledTime > c.PushedTime {
		return c.PulledTime
	}

	return c.PushedTime
}
