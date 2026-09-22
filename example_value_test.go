// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2018 The Go Authors. All rights reserved.

package zflag_test

import (
	"fmt"
	"net/url"

	"github.com/zulucmd/zflag/v2"
)

type URLValue struct {
	URL *url.URL
}

func (v URLValue) String() string {
	if v.URL != nil {
		return v.URL.String()
	}
	return ""
}

func (v URLValue) Set(s string) error {
	u, err := url.Parse(s)
	if err != nil {
		return err
	}

	*v.URL = *u

	return nil
}

var u = &url.URL{}

func ExampleValue() {
	fs := zflag.NewFlagSet("ExampleValue", zflag.ExitOnError)
	fs.Var(&URLValue{u}, "url", "URL to parse")

	_ = fs.Parse([]string{"--url", "https://golang.org/pkg/flag/"})
	fmt.Printf(`{scheme: %q, host: %q, path: %q}`, u.Scheme, u.Host, u.Path)

	// Output:
	// {scheme: "https", host: "golang.org", path: "/pkg/flag/"}
}
