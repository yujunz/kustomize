// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package loader has a data loading interface and various implementations.
package loader

import (
	"context"
	"log"
	"os"

	getter "github.com/yujunz/go-getter"
	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/api/ifc"
)

type targetSpec struct {
	// Raw is the original resource in kustomization.yaml
	Raw string

	// Dir is where the resource is saved
	Dir filesys.ConfirmedDir

	// TempDir is the directory created to hold all resources, including Dir
	TempDir filesys.ConfirmedDir
}

// Getter is a function that can gets resource
type targetGetter func(rs *targetSpec) error

// NewLoader returns a Loader pointed at the given target.
// If the target is remote, the loader will be restricted
// to the root and below only.  If the target is local, the
// loader will have the restrictions passed in.  Regardless,
// if a local target attempts to transitively load remote bases,
// the remote bases will all be root-only restricted.
func NewLoader(
	lr LoadRestrictorFunc,
	target string, fSys filesys.FileSystem) (ifc.Loader, error) {

	rs := &targetSpec{
		Raw: raw,
	}

	cleaner := func() error {
		return fSys.RemoveAll(rs.TempDir.String())
	}

	if err := getter(rs); err != nil {
		cleaner()
		return nil, err
	}

	return &fileLoader{
		loadRestrictor: RestrictionRootOnly,
		// TODO(yujunz): limit to getter root
		root:     rs.Dir,
		referrer: nil,
		fSys:     fSys,
		rscSpec:  rs,
		getter:   getter,
		cleaner:  cleaner,
	}, nil
}

func getTarget(rs *targetSpec) error {
	var err error

	rs.TempDir, err = filesys.NewTmpConfirmedDir()
	if err != nil {
		return err
	}

	rs.Dir = filesys.ConfirmedDir(rs.TempDir.Join("repo"))

	// Get the pwd
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}

	opts := []getter.ClientOption{}
	client := &getter.Client{
		Ctx:  context.TODO(),
		Src:  rs.Raw,
		Dst:  rs.Dir.String(),
		Pwd:  pwd,
		Mode: getter.ClientModeAny,
		Detectors: []getter.Detector{
			new(getter.GitHubDetector),
			new(getter.GitLabDetector),
			new(getter.GitDetector),
			new(getter.BitBucketDetector),
			new(getter.FileDetector),
		},
		Options: opts,
	}
	return client.Get()
}

func getNothing(rs *targetSpec) error {
	var err error
	rs.Dir, err = filesys.NewTmpConfirmedDir()
	if err != nil {
		return err
	}

	// Get the pwd
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd: %s", err)
	}

	_, err = getter.Detect(rs.Raw, pwd, []getter.Detector{})
	return err
}
