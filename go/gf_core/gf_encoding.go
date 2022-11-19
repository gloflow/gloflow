/*
MIT License

Copyright (c) 2022 Ivan Trajkovic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package gf_core

import (
	"encoding/json"
)

//---------------------------------------------------

func EncodeJSONfromMap(pInputMap map[string]interface{}) []byte {
	outputLst, _ := json.Marshal(pInputMap)
	return outputLst
}

//---------------------------------------------------

func ParseJSONfromByteList(pBytesLst []byte,
	pRuntimeSys *RuntimeSys) (map[string]interface{}, *GFerror) {

	var outputMap map[string]interface{}
	err := json.Unmarshal(pBytesLst, &outputMap)

	if err != nil {
		gfErr := ErrorCreate("failed to parse json byte list",
			"json_decode_error",
			map[string]interface{}{
				"json_bytes_lst": pBytesLst,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	return outputMap, nil
}

//---------------------------------------------------

func ParseJSONfromString(pJSONstr string,
	pRuntimeSys *RuntimeSys) (map[string]interface{}, *GFerror) {

	var outputMap map[string]interface{}
	err := json.Unmarshal([]byte(pJSONstr), &outputMap)

	if err != nil {
		gfErr := ErrorCreate("failed to parse json string",
			"json_decode_error",
			map[string]interface{}{
				"json_str": pJSONstr,
			},
			err, "gf_core", pRuntimeSys)
		return nil, gfErr
	}

	return outputMap, nil
}