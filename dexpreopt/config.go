// Copyright 2018 Google Inc. All rights reserved.
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

package dexpreopt

import (
	"encoding/json"
	"io/ioutil"

	"android/soong/android"
)

// GlobalConfig stores the configuration for dex preopting set by the product
type GlobalConfig struct {
	DefaultNoStripping bool // don't strip dex files by default

	DisablePreoptModules []string // modules with preopt disabled by product-specific config

	OnlyPreoptBootImageAndSystemServer bool // only preopt jars in the boot image or system server

	HasSystemOther        bool     // store odex files that match PatternsOnSystemOther on the system_other partition
	PatternsOnSystemOther []string // patterns (using '%' to denote a prefix match) to put odex on the system_other partition

	DisableGenerateProfile bool // don't generate profiles

	BootJars []string // modules for jars that form the boot class path

	RuntimeApexJars               []string // modules for jars that are in the runtime apex
	ProductUpdatableBootModules   []string
	ProductUpdatableBootLocations []string

	SystemServerJars []string // jars that form the system server
	SystemServerApps []string // apps that are loaded into system server
	SpeedApps        []string // apps that should be speed optimized

	PreoptFlags []string // global dex2oat flags that should be used if no module-specific dex2oat flags are specified

	DefaultCompilerFilter      string // default compiler filter to pass to dex2oat, overridden by --compiler-filter= in module-specific dex2oat flags
	SystemServerCompilerFilter string // default compiler filter to pass to dex2oat for system server jars

	GenerateDMFiles     bool // generate Dex Metadata files
	NeverAllowStripping bool // whether stripping should not be done - used as build time check to make sure dex files are always available

	NoDebugInfo                 bool // don't generate debug info by default
	AlwaysSystemServerDebugInfo bool // always generate mini debug info for system server modules (overrides NoDebugInfo=true)
	NeverSystemServerDebugInfo  bool // never generate mini debug info for system server modules (overrides NoDebugInfo=false)
	AlwaysOtherDebugInfo        bool // always generate mini debug info for non-system server modules (overrides NoDebugInfo=true)
	NeverOtherDebugInfo         bool // never generate mini debug info for non-system server modules (overrides NoDebugInfo=true)

	MissingUsesLibraries []string // libraries that may be listed in OptionalUsesLibraries but will not be installed by the product

	IsEng        bool // build is a eng variant
	SanitizeLite bool // build is the second phase of a SANITIZE_LITE build

	DefaultAppImages bool // build app images (TODO: .art files?) by default

	Dex2oatXmx string // max heap size for dex2oat
	Dex2oatXms string // initial heap size for dex2oat

	EmptyDirectory string // path to an empty directory

	CpuVariant             map[android.ArchType]string // cpu variant for each architecture
	InstructionSetFeatures map[android.ArchType]string // instruction set for each architecture

	// Only used for boot image
	DirtyImageObjects string   // path to a dirty-image-objects file
	PreloadedClasses  string   // path to a preloaded-classes file
	BootImageProfiles []string // path to a boot-image-profile.txt file
	BootFlags         string   // extra flags to pass to dex2oat for the boot image
	Dex2oatImageXmx   string   // max heap size for dex2oat for the boot image
	Dex2oatImageXms   string   // initial heap size for dex2oat for the boot image

	Tools Tools // paths to tools possibly used by the generated commands
}

// Tools contains paths to tools possibly used by the generated commands.  If you add a new tool here you MUST add it
// to the order-only dependency list in DEXPREOPT_GEN_DEPS.
type Tools struct {
	Profman  string
	Dex2oat  string
	Aapt     string
	SoongZip string
	Zip2zip  string

	VerifyUsesLibraries string
	ConstructContext    string
}

type ModuleConfig struct {
	Name            string
	DexLocation     string // dex location on device
	BuildPath       string
	DexPath         string
	UncompressedDex bool
	HasApkLibraries bool
	PreoptFlags     []string

	ProfileClassListing  string
	ProfileIsTextListing bool

	EnforceUsesLibraries  bool
	OptionalUsesLibraries []string
	UsesLibraries         []string
	LibraryPaths          map[string]string

	Archs           []android.ArchType
	DexPreoptImages []string

	PreoptBootClassPathDexFiles     []string // file paths of boot class path files
	PreoptBootClassPathDexLocations []string // virtual locations of boot class path files

	PreoptExtractedApk bool // Overrides OnlyPreoptModules

	NoCreateAppImage    bool
	ForceCreateAppImage bool

	PresignedPrebuilt bool

	NoStripping     bool
	StripInputPath  string
	StripOutputPath string
}

func LoadGlobalConfig(path string) (GlobalConfig, error) {
	config := GlobalConfig{}
	err := loadConfig(path, &config)
	return config, err
}

func LoadModuleConfig(path string) (ModuleConfig, error) {
	config := ModuleConfig{}
	err := loadConfig(path, &config)
	return config, err
}

func loadConfig(path string, config interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
