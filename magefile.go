//go:build mage
// +build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Hello is our default mage target which we also call
// by default within our Docker container
func Hello() error {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("hello %s\n", now)
	return nil
}

// Goodbye is an alternative mage target we can call
func Goodbye() error {
	now := time.Now().UTC().Format(time.RFC3339)
	fmt.Printf("goodbye %s\n", now)
	return nil
}

type Docker mg.Namespace

// BuildDev builds the container image, "jobs", with --no-cache
// and Dockerfile which builds a static binary and
// multi-stage builds to utilize a distroless image
func (Docker) Build() error {
	cmd1 := []string{
		"docker",
		"build",
		"--no-cache",
		"-t",
		"jobs",
		".",
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// BuildDev builds the container image, "jobs", with --no-cache
// and dev.Dockerfile which uses the golang:latest image,
// installs mage and vim, for more interactive development
func (Docker) BuildDev() error {
	cmd1 := []string{
		"docker",
		"build",
		"--no-cache",
		"-f",
		"dev.Dockerfile",
		"-t",
		"jobs",
		".",
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}

// Run runs the jobs container with the mage target
func (Docker) Run(target string) error {
	cmd1 := []string{
		"docker",
		"run",
		"-it",
		"jobs",
		"mage",
		target,
	}
	return sh.RunV(cmd1[0], cmd1[1:]...)
}
