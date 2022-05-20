// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package tracing

import (
	"context"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"testing"

	"github.com/cilium/tetragon/api/v1/tetragon"
	ec "github.com/cilium/tetragon/api/v1/tetragon/codegen/eventchecker"
	bc "github.com/cilium/tetragon/api/v1/tetragon/codegen/eventchecker/matchers/bytesmatcher"
	lc "github.com/cilium/tetragon/api/v1/tetragon/codegen/eventchecker/matchers/listmatcher"
	sm "github.com/cilium/tetragon/api/v1/tetragon/codegen/eventchecker/matchers/stringmatcher"
	"github.com/cilium/tetragon/pkg/kernels"
	"github.com/cilium/tetragon/pkg/observer"
	"github.com/cilium/tetragon/pkg/testutils"
	"github.com/stretchr/testify/assert"

	_ "github.com/cilium/tetragon/pkg/sensors/exec"
)

func TestKprobeNSChanges(t *testing.T) {
	if !kernels.MinKernelVersion("5.3.0") {
		t.Skip("matchNamespaceChanges requires at least 5.3.0 version")
	}

	var doneWG, readyWG sync.WaitGroup
	defer doneWG.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), cmdWaitTime)
	defer cancel()

	testBin := testutils.ContribPath("tester-progs/namespace-tester")
	testCmd := exec.CommandContext(ctx, testBin)
	testPipes, err := testutils.NewCmdBufferedPipes(testCmd)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		logOut(t, "stdout> ", testPipes.StdoutRd)
	}()

	go func() {
		logOut(t, "stderr> ", testPipes.StderrRd)
	}()

	// makeSpecFile creates a new spec file bsed on the template, and the provided arguments
	makeSpecFile := func(pid string) string {
		data := map[string]string{
			"MatchedPID":   pid,
			"NamespacePID": "false",
		}
		specName, err := testutils.GetSpecFromTemplate("nschanges.yaml.tmpl", data)
		if err != nil {
			t.Fatal(err)
		}
		return specName
	}

	pidStr := strconv.Itoa(os.Getpid())
	specFname := makeSpecFile(pidStr)
	t.Logf("pid is %s and spec file is %s", pidStr, specFname)

	obs, err := observer.GetDefaultObserverWithFile(t, specFname, tetragonLib)
	if err != nil {
		t.Fatalf("GetDefaultObserverWithFile error: %s", err)
	}
	observer.LoopEvents(ctx, t, &doneWG, &readyWG, obs)
	readyWG.Wait()

	if err := testCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := testCmd.Wait(); err != nil {
		t.Fatalf("command failed with %s. Context error: %s", err, ctx.Err())
	}

	writeArg0 := ec.NewKprobeArgumentChecker().
		WithFileArg(ec.NewKprobeFileChecker().
			WithPath(sm.Suffix("strange.txt")).
			WithFlags(sm.Full("")),
		)
	writeArg1 := ec.NewKprobeArgumentChecker().WithBytesArg(bc.Full([]byte("testdata\x00")))
	writeArg2 := ec.NewKprobeArgumentChecker().WithSizeArg(9)

	kprobeChecker := ec.NewProcessKprobeChecker().
		WithFunctionName(sm.Full("__x64_sys_write")).
		WithArgs(ec.NewKprobeArgumentListMatcher().
			WithOperator(lc.Ordered).
			WithValues(
				writeArg0,
				writeArg1,
				writeArg2,
			))

	checker := ec.NewUnorderedEventChecker(
		kprobeChecker,
	)
	err = observer.JsonTestCheck(t, checker)
	assert.NoError(t, err)
}

func TestKprobeCapChanges(t *testing.T) {
	if !kernels.MinKernelVersion("5.3.0") {
		t.Skip("matchCapabilityChanges requires at least 5.3.0 version")
	}

	var doneWG, readyWG sync.WaitGroup
	defer doneWG.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), cmdWaitTime)
	defer cancel()

	testBin := testutils.ContribPath("tester-progs/capabilities-tester")
	testCmd := exec.CommandContext(ctx, testBin)
	testPipes, err := testutils.NewCmdBufferedPipes(testCmd)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		logOut(t, "stdout> ", testPipes.StdoutRd)
	}()

	go func() {
		logOut(t, "stderr> ", testPipes.StderrRd)
	}()

	// makeSpecFile creates a new spec file bsed on the template, and the provided arguments
	makeSpecFile := func(pid string) string {
		data := map[string]string{
			"MatchedPID":   pid,
			"NamespacePID": "false",
		}
		specName, err := testutils.GetSpecFromTemplate("capchanges.yaml.tmpl", data)
		if err != nil {
			t.Fatal(err)
		}
		return specName
	}

	pidStr := strconv.Itoa(os.Getpid())
	specFname := makeSpecFile(pidStr)
	t.Logf("pid is %s and spec file is %s", pidStr, specFname)

	obs, err := observer.GetDefaultObserverWithFile(t, specFname, tetragonLib)
	if err != nil {
		t.Fatalf("GetDefaultObserverWithFile error: %s", err)
	}
	observer.LoopEvents(ctx, t, &doneWG, &readyWG, obs)
	readyWG.Wait()

	if err := testCmd.Start(); err != nil {
		t.Fatal(err)
	}

	if err := testCmd.Wait(); err != nil {
		t.Fatalf("command failed with %s. Context error: %s", err, ctx.Err())
	}

	writeArg0 := ec.NewKprobeArgumentChecker().
		WithFileArg(ec.NewKprobeFileChecker().
			WithPath(sm.Suffix("strange.txt")).
			WithFlags(sm.Full("")),
		)
	writeArg1 := ec.NewKprobeArgumentChecker().WithBytesArg(bc.Full([]byte("testdata\x00")))
	writeArg2 := ec.NewKprobeArgumentChecker().WithSizeArg(9)

	kprobeChecker := ec.NewProcessKprobeChecker().
		WithFunctionName(sm.Full("__x64_sys_write")).
		WithArgs(ec.NewKprobeArgumentListMatcher().
			WithOperator(lc.Ordered).
			WithValues(
				writeArg0,
				writeArg1,
				writeArg2,
			)).
		WithAction(tetragon.KprobeAction_KPROBE_ACTION_POST)

	checker := ec.NewUnorderedEventChecker(
		kprobeChecker,
	)

	err = observer.JsonTestCheck(t, checker)
	assert.NoError(t, err)
}
