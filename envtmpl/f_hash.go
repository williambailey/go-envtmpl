package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"sort"
	"strings"
)

func init() {
	algo := map[string]func() hash.Hash{
		"adler32": func() hash.Hash {
			h, _ := adler32.New().(hash.Hash)
			return h
		},
		"crc32": func() hash.Hash {
			h, _ := crc32.NewIEEE().(hash.Hash)
			return h
		},
		"crc64iso": func() hash.Hash {
			h, _ := crc64.New(crc64.MakeTable(crc64.ISO)).(hash.Hash)
			return h
		},
		"crc64ecma": func() hash.Hash {
			h, _ := crc64.New(crc64.MakeTable(crc64.ECMA)).(hash.Hash)
			return h
		},
		"fnv1-32": func() hash.Hash {
			h, _ := fnv.New32().(hash.Hash)
			return h
		},
		"fnv1a-32": func() hash.Hash {
			h, _ := fnv.New32a().(hash.Hash)
			return h
		},
		"fnv1-64": func() hash.Hash {
			h, _ := fnv.New64().(hash.Hash)
			return h
		},
		"fnv1a-64": func() hash.Hash {
			h, _ := fnv.New64a().(hash.Hash)
			return h
		},
		"md5":    md5.New,
		"sha1":   sha1.New,
		"sha224": sha256.New224,
		"sha256": sha256.New,
		"sha384": sha512.New384,
		"sha512": sha512.New,
	}
	var idList []string
	var exList []string
	for a := range algo {
		idList = append(idList, a)
		exList = append(
			exList,
			fmt.Sprintf("{{ \"Hello World!\" | %%[1]s \"%[1]s\" }}\n{{ \"Hello World!\" | %%[1]s \"%[1]s\" \"a key\" }}", a),
		)
	}
	sort.Strings(idList)
	funcMap["hash"] = &tmplFuncStruct{
		short: fmt.Sprintf(
			"Calculate the hex encoded hash of a string. You can optionally specify a key to produce a HMAC string. The following hash algorithms are supported: %s.",
			strings.Join(idList, ", "),
		),
		examples: exList,
		fn: func(in ...string) (string, error) {
			var data, hashName, macKey string
			var h hash.Hash
			switch len(in) {
			case 2:
				data = in[1]
				hashName = in[0]
			case 3:
				data = in[2]
				hashName = in[0]
				macKey = in[1]
			default:
				return "", errors.New("Expecting 2 or 3 arguments.")
			}
			hFn, validAlgo := algo[strings.ToLower(hashName)]
			if !validAlgo {
				return "", fmt.Errorf("Unknown hash algorithm '%s'.", hashName)
			}
			if macKey == "" {
				h = hFn()
			} else {
				h = hmac.New(hFn, []byte(macKey))
			}
			_, err := h.Write([]byte(data))
			if err != nil {
				return "", err
			}
			return hex.EncodeToString(h.Sum(nil)), nil
		},
	}
}
