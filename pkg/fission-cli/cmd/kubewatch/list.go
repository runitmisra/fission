/*
Copyright 2019 The Fission Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubewatch

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/pkg/errors"

	v1 "github.com/fission/fission/pkg/apis/core/v1"
	"github.com/fission/fission/pkg/fission-cli/cliwrapper/cli"
	"github.com/fission/fission/pkg/fission-cli/cmd"
	flagkey "github.com/fission/fission/pkg/fission-cli/flag/key"
	"github.com/fission/fission/pkg/fission-cli/util"
)

type ListSubCommand struct {
	cmd.CommandActioner
	namespace string
}

func List(input cli.Input) error {
	return (&ListSubCommand{}).do(input)
}

func (opts *ListSubCommand) do(input cli.Input) error {
	err := opts.complete(input)
	if err != nil {
		return err
	}
	return opts.run(input)
}

func (opts *ListSubCommand) complete(input cli.Input) (err error) {
	_, opts.namespace, err = util.GetResourceNamespace(input, flagkey.NamespaceTrigger)
	if err != nil {
		return errors.Wrap(err, "error listing kubewatchers")
	}
	return nil
}

func (opts *ListSubCommand) run(input cli.Input) (err error) {
	var ws []v1.KubernetesWatchTrigger
	if input.Bool(flagkey.AllNamespaces) {
		ws, err = opts.Client().V1().KubeWatcher().List("")
	} else {
		ws, err = opts.Client().V1().KubeWatcher().List(opts.namespace)
	}
	if err != nil {
		return errors.Wrap(err, "error listing kubewatchers")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
		"NAME", "NAMESPACE", "OBJTYPE", "LABELS", "FUNCTION_NAME")
	for _, wa := range ws {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n",
			wa.ObjectMeta.Name, wa.Spec.Namespace, wa.Spec.Type, wa.Spec.LabelSelector, wa.Spec.FunctionReference.Name)
	}
	w.Flush()

	return nil
}
