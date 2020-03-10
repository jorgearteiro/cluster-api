/*
Copyright 2019 The Kubernetes Authors.

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

package cloudinit

import (
	"testing"

	. "github.com/onsi/gomega"

	infrav1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/certs"
	"sigs.k8s.io/cluster-api/util/secret"
)

func TestNewInitControlPlaneAdditionalFileEncodings(t *testing.T) {
	g := NewWithT(t)

	cpinput := &ControlPlaneInput{
		BaseUserData: BaseUserData{
			Header:              "test",
			PreKubeadmCommands:  nil,
			PostKubeadmCommands: nil,
			AdditionalFiles: []infrav1.File{
				{
					Path:     "/tmp/my-path",
					Encoding: infrav1.Base64,
					Content:  "aGk=",
				},
				{
					Path:    "/tmp/my-other-path",
					Content: "hi",
				},
			},
			WriteFiles: nil,
			Users:      nil,
			NTP:        nil,
		},
		Certificates:         secret.Certificates{},
		ClusterConfiguration: "my-cluster-config",
		InitConfiguration:    "my-init-config",
	}

	for _, certificate := range cpinput.Certificates {
		certificate.KeyPair = &certs.KeyPair{
			Cert: []byte("some certificate"),
			Key:  []byte("some key"),
		}
	}

	out, err := NewInitControlPlane(cpinput)
	g.Expect(err).NotTo(HaveOccurred())

	expectedFiles := []string{
		`-   path: /tmp/my-path
    encoding: "base64"
    content: |
      aGk=`,
		`-   path: /tmp/my-other-path
    content: |
      hi`,
	}
	for _, f := range expectedFiles {
		g.Expect(out).To(ContainSubstring(f))
	}
}

func TestNewInitControlPlaneCommands(t *testing.T) {
	g := NewWithT(t)

	cpinput := &ControlPlaneInput{
		BaseUserData: BaseUserData{
			Header:              "test",
			PreKubeadmCommands:  []string{`"echo $(date) ': hello world!'"`},
			PostKubeadmCommands: []string{"echo $(date) ': hello world!'"},
			AdditionalFiles:     nil,
			WriteFiles:          nil,
			Users:               nil,
			NTP:                 nil,
		},
		Certificates:         secret.Certificates{},
		ClusterConfiguration: "my-cluster-config",
		InitConfiguration:    "my-init-config",
	}

	for _, certificate := range cpinput.Certificates {
		certificate.KeyPair = &certs.KeyPair{
			Cert: []byte("some certificate"),
			Key:  []byte("some key"),
		}
	}

	out, err := NewInitControlPlane(cpinput)
	g.Expect(err).NotTo(HaveOccurred())

	expectedCommands := []string{
		`"\"echo $(date) ': hello world!'\""`,
		`"echo $(date) ': hello world!'"`,
	}
	for _, f := range expectedCommands {
		g.Expect(out).To(ContainSubstring(f))
	}
}
