Go httputils
============

[![Build Status](https://travis-ci.org/enr/go-httputils.png?branch=master)](https://travis-ci.org/enr/go-httputils)
[![Build status](https://ci.appveyor.com/api/projects/status/moywc16vawynrvsx?svg=true)](https://ci.appveyor.com/project/enr/go-httputils)

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

Test if it is a valid URL:

```Go
    isValidUrl := httputils.IsValidUrl(maybeUrl)
```


License
-------

Apache 2.0 - see LICENSE file.

   Copyright 2014 go-httputils contributors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
