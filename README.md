Go httputils
============

[![Build Status](https://travis-ci.org/enr/go-httputils.png?branch=master)](https://travis-ci.org/enr/go-httputils)

Go library to easily download files via HTTP.

Import:

```Go
    import (
        "github.com/enr/go-httputils/httputils"
    )
```

Download only if destination file does not exist:

```Go
    err := httputils.DownloadIfNotExists(url, localFilePath)
    if err != nil {
        // ...
    }
```

Download overwriting destination file, if it exists:

```Go
    err := httputils.DownloadOverwriting(url, localFilePath)
    if err != nil {
        // ...
    }
```

Download preserving destination file if it exists; copy `localFilePath` in `backupFilePath` and download `url` in `localFilePath`:

```Go
    err := httputils.DownloadPreservingOld(url, localFilePath, backupFilePath)
    if err != nil {
        // ...
    }
```


License
-------

Apache 2.0 - see LICENSE file.
