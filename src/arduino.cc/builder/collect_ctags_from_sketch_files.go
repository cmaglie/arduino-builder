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
 * Copyright 2015 Arduino LLC (http://www.arduino.cc/)
 */

package builder

import (
	"arduino.cc/builder/constants"
	"arduino.cc/builder/types"
	"arduino.cc/builder/utils"
	"strings"
)

type CollectCTagsFromSketchFiles struct{}

func (s *CollectCTagsFromSketchFiles) Run(context map[string]interface{}) error {
	sketch := context[constants.CTX_SKETCH].(*types.Sketch)
	sketchFileNames := collectSketchFileNamesFrom(sketch)

	ctags := context[constants.CTX_CTAGS_OF_PREPROC_SOURCE].([]map[string]string)
	ctagsOfSketch := []map[string]string{}
	for _, ctag := range ctags {
		if utils.SliceContains(sketchFileNames, strings.Replace(ctag[FIELD_FILENAME], "\\\\", "\\", -1)) {
			ctagsOfSketch = append(ctagsOfSketch, ctag)
		}
	}

	context[constants.CTX_COLLECTED_CTAGS] = ctagsOfSketch

	return nil
}

func collectSketchFileNamesFrom(sketch *types.Sketch) []string {
	fileNames := []string{sketch.MainFile.Name}
	for _, file := range sketch.OtherSketchFiles {
		fileNames = append(fileNames, file.Name)
	}
	return fileNames
}
