package utils

import (
    "path"
    "runtime"
)

func ProjectRoot() string {
    var _, filename, _, _ = runtime.Caller(0)
    return path.Join(path.Dir(filename), "..")
}