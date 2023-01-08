/*
GloFlow application and media management/publishing platform
Copyright (C) 2019 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/
	
package gf_core

import (
	"fmt"
	"time"
	"crypto/sha256"
	"math/rand"
)

//-------------------------------------------------------------
// STRING
//-------------------------------------------------------------

func StrRandom() string {

	randWithSeed := rand.New(rand.NewSource(time.Now().UnixNano()))
	randInt := randWithSeed.Int()
	randStr := fmt.Sprintf("%d", randInt)
	return randStr
}

func CastToStr(pElement interface{}) (bool, string) {
    var isStringBool bool
    var elementStr string
    if valStr, ok := pElement.(string); ok {
        elementStr = valStr
        isStringBool = true
    }
    return isStringBool, elementStr
}

//-------------------------------------------------------------

func HashValSha256(pVal interface{}) string {

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", pVal)))
	hashHexStr := fmt.Sprintf("%x", h.Sum(nil))
	return hashHexStr
}

//-------------------------------------------------------------
// MAP
//-------------------------------------------------------------

func MapHasKey[K string, V any](pMap map[K]V, pKeyStr string) bool {
	keysLst := make([]string, 0, len(pMap))
    for k := range pMap {
        keysLst = append(keysLst, string(k))
    }

	return ListContainsStr(pKeyStr, []string(keysLst))
}

//-------------------------------------------------------------
// LIST
//-------------------------------------------------------------

func ListContainsStr(pStr string, pLst []string) bool {
	for _, s := range pLst {
		if pStr == s {
			return true
		}
	}
	return false
}

func ListRemoveElementAtIndex(pLst []interface{}, pIndex int) []interface{} {
	newLst := make([]interface{}, 0)
	newLst = append(newLst, pLst[:pIndex]...)
	newLst = append(newLst, pLst[pIndex+1:]...)
	return newLst
}

func ListPop[T any](pLst []T) (T, []T) {

	var lastVal T
	if len(pLst) == 0 {
		return lastVal, pLst
	}

	lastVal       = pLst[len(pLst)-1]
    remainingLst := pLst[:len(pLst)-1]
    return lastVal, remainingLst
}