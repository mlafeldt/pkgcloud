package main

import (
	"errors"
	"strings"
)

type target struct {
	repo   string
	distro string
}

func newTarget(s string) (*target, error) {
	if s == "" {
		return nil, errors.New("empty target")
	}
	elems := strings.Split(s, "/")
	switch len(elems) {
	case 2:
		return &target{
			repo: strings.Join(elems[0:2], "/"),
		}, nil
	case 3:
		return &target{
			repo:   strings.Join(elems[0:2], "/"),
			distro: elems[2],
		}, nil
	case 4:
		return &target{
			repo:   strings.Join(elems[0:2], "/"),
			distro: strings.Join(elems[2:4], "/"),
		}, nil
	}
	return nil, errors.New("invalid target")
}

func (t target) String() string {
	if t.distro == "" {
		return t.repo
	}
	return strings.Join([]string{t.repo, t.distro}, "/")
}
