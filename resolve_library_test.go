/*
 * This file is part of Arduino Builder.
 *
 * Arduino Builder is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
 *
 * As a special exception, you may use this file as part of a free software
 * library without restriction.  Specifically, if other files instantiate
 * templates or use macros or inline functions from this file, or you compile
 * this file and link it with other files to produce an executable, this
 * file does not by itself cause the resulting executable to be covered by
 * the GNU General Public License.  This exception does not however
 * invalidate any other reasons why the executable file might be covered by
 * the GNU General Public License.
 *
 * Copyright 2017 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"testing"

	"github.com/arduino/arduino-builder/types"

	"github.com/stretchr/testify/require"
)

func testSelection(include string, arch *types.Platform, libs ...*types.Library) *types.Library {
	ctx := &types.Context{
		HeaderToLibraries:          map[string][]*types.Library{include: libs},
		ImportedLibraries:          []*types.Library{},
		ActualPlatform:             arch,
		TargetPlatform:             arch,
		LibrariesResolutionResults: map[string]types.LibraryResolutionResult{},
	}
	return ResolveLibrary(ctx, include)
}

func TestArchPriority(t *testing.T) {
	var bundleServo = &types.Library{Name: "Servo", Archs: []string{"avr", "sam", "samd"}}
	var userServo = &types.Library{Name: "Servo", Folder: "sketchbook", Archs: []string{"avr", "sam", "samd"}}
	var userServoNonAvr = &types.Library{Name: "Servo", Folder: "sketchbook", Archs: []string{"sam", "samd"}}
	var userAnotherServo = &types.Library{Name: "AnotherServo", Folder: "sketchbook", Archs: []string{"avr", "sam", "samd", "esp32"}}
	var userServoAllArch = &types.Library{Name: "Servo", Folder: "sketchbook", Archs: []string{"*"}}

	avr := &types.Platform{PlatformId: "avr"}
	esp32 := &types.Platform{PlatformId: "esp32"}
	esp8266 := &types.Platform{PlatformId: "esp8266"}

	res := testSelection("Servo.h", avr, bundleServo, userServo)
	require.NotNil(t, res)
	require.Equal(t, userServo, res)

	res = testSelection("Servo.h", avr, bundleServo, userServoNonAvr)
	require.NotNil(t, res)
	require.Equal(t, bundleServo, res)

	res = testSelection("Servo.h", avr, bundleServo, userAnotherServo)
	require.NotNil(t, res)
	require.Equal(t, bundleServo, res)

	res = testSelection("Servo.h", esp32, bundleServo, userAnotherServo)
	require.NotNil(t, res)
	require.Equal(t, userAnotherServo, res)

	res = testSelection("Servo.h", esp32, userAnotherServo, userServoAllArch)
	require.NotNil(t, res)
	require.Equal(t, userServoAllArch, res)

	var bundleSDesp = &types.Library{Name: "SD", Archs: []string{"esp8266"}}
	var userSDAllArch = &types.Library{Name: "SD", Folder: "sketchbook", Archs: []string{"*"}}
	res = testSelection("SD.h", esp8266, userSDAllArch, bundleSDesp)
	require.NotNil(t, res)
	require.Equal(t, bundleSDesp, res)
}

func TestFindBestLibraryWithHeader(t *testing.T) {
	l1 := &types.Library{Name: "Calculus Lib"}
	l2 := &types.Library{Name: "Calculus Lib-master"}
	l3 := &types.Library{Name: "Calculus Lib Improved"}
	l4 := &types.Library{Name: "Another Calculus Lib"}
	l5 := &types.Library{Name: "Yet Another Calculus Lib Improved"}
	l6 := &types.Library{Name: "Calculus Unified Lib"}
	l7 := &types.Library{Name: "AnotherLib"}

	// Test exact name matching
	res := findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6, l5, l4, l3, l2, l1})
	require.Equal(t, l1.Name, res.Name)

	// Test exact name with "-master" postfix matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6, l5, l4, l3, l2})
	require.Equal(t, l2.Name, res.Name)

	// Test prefix matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6, l5, l4, l3})
	require.Equal(t, l3.Name, res.Name)

	// Test postfix matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6, l5, l4})
	require.Equal(t, l4.Name, res.Name)

	// Test "contains"" matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6, l5})
	require.Equal(t, l5.Name, res.Name)

	// Test lexicographic similarity matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7, l6})
	require.Equal(t, l6.Name, res.Name)

	// Test lexicographic similarity matching (2)
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l6, l7})
	require.Equal(t, l6.Name, res.Name)

	// Test none matching
	res = findBestLibraryWithHeader("calculus_lib.h", []*types.Library{l7})
	require.Equal(t, l7.Name, res.Name)
}
