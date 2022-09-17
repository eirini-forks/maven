/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package maven_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/maven/v6/maven"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect maven.Detect
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "maven")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
	})

	it("only provides with pom.xml", func() {
		os.Setenv("BP_MAVEN_POM_FILE", "pom.xml")
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "maven"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
						{Name: "maven"},
					},
				},
			},
		}))
	})

	it("passes with pom.xml", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "pom.xml"), []byte{}, 0644))

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "maven"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "syft"},
						{Name: "jdk"},
						{Name: "maven"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
						{Name: "maven"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "syft"},
						{Name: "jdk"},
						{Name: "maven"},
					},
				},
			},
		}))
	})

	it("passes with a custom pom.xml", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "pom2.xml"), []byte{}, 0644))

		os.Setenv("BP_MAVEN_POM_FILE", "pom2.xml")

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "maven"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "jdk"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "syft"},
						{Name: "jdk"},
						{Name: "maven"},
					},
				},
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "jvm-application-package"},
						{Name: "maven"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "syft"},
						{Name: "jdk"},
						{Name: "maven"},
					},
				},
			},
		}))
	})
}
