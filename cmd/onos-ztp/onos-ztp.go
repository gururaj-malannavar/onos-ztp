// Copyright 2019-present Open Networking Foundation.
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

/*
Package onos-ztp is the main entry point to the ONOS zero-touch provisioning subsystem.

Arguments

-caPath <the location of a CA certificate>

-keyPath <the location of a client private key>

-certPath <the location of a client certificate>


See ../../docs/run.md for how to run the application.
*/
package main

import (
	"flag"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/northbound"
	"github.com/onosproject/onos-ztp/pkg/manager"
	"github.com/onosproject/onos-ztp/pkg/northbound/admin"
	"github.com/onosproject/onos-ztp/pkg/northbound/roles"
)

var log = logging.GetLogger("main")

// The main entry point
func main() {
	caPath := flag.String("caPath", "", "path to CA certificate")
	keyPath := flag.String("keyPath", "", "path to client private key")
	certPath := flag.String("certPath", "", "path to client certificate")
	flag.Parse()

	log.Info("Starting onos-ztp")

	opts, err := certs.HandleCertArgs(keyPath, certPath)
	if err != nil {
		log.Fatal(err)
	}

	mgr, err := manager.LoadManager("roledb", opts...)
	if err != nil {
		log.Fatal("Unable to load onos-ztp ", err)
	} else {
		mgr.Run()
		err = startServer(*caPath, *keyPath, *certPath)
		if err != nil {
			log.Fatal("Unable to start onos-ztp ", err)
		}
	}
}

// Creates gRPC server and registers various services; then serves.
func startServer(caPath string, keyPath string, certPath string) error {
	s := northbound.NewServer(northbound.NewServerConfig(caPath, keyPath, certPath, 5150, true))
	s.AddService(admin.Service{})
	s.AddService(roles.Service{})
	s.AddService(logging.Service{})

	return s.Serve(func(started string) {
		log.Info("Started NBI on ", started)
	})
}
