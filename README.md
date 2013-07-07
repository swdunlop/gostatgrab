gostatgrab
==========

Gostatgrab is a Go wrapper for the [statgrab]() library from the i-scream group.  It closely reproduces the API provided by libstatgrab, providing CPU usage, VM usage, disk statistics, process information, network information and other resource information across a wide range of UNIX platforms.  (But, like many current projects, we really only expect it will run on Linux.)

[statgrab]: https://github.com/i-scream/libstatgrab

## License

Statgrab itself is LGPL; Gostatgrab is BSD licensed, which is as close as we come to "we don't care, it's just wrapper code, man."

## Building Gostatgrab

Gostatgrab requires libstatgrab, and uses pkg-config to determine how to compile against libstatgrab.  Most Linux distributions have statgrab packages in their community repositories.

Gostatgrab uses "go build" as a build system; with a weekly Go installation, as of 2/2012, this is as easy as "go get github.com/swdunlop/gostatgrab"

## Using Gostatgrab

An example, jsongrab, demonstrates how to use the API.  On most modern platforms, elevated privileges are not required; on other platforms, Gostatgrab will drop privileges as soon as possible -- once libstatgrab has kicked off its pipe.

Once jsongrab has been read and digested, the rest follows naturally after reading the statgrab(3) manpages.

## Improving Gostatgrab

Please submit any improvements as github pull requests, or if religiously opposed, diffs.  Pull requests tend to be faster, as they are quicker off the block for review time.

With bug reports, be sure to include OS distribution, Go version and Statgrab version.

## Why We Did It?

We use Go for weird little widgets in appliances, and wanted a system health monitor that integrated into our system.  Since statgrab thoroughly wrings the weirdness out of many of these API's, we decided to wrap that instead of rolling yet another POSIX abstraction layer.

I suppose we could have used Python..

-- Scott.
