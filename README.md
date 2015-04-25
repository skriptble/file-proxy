# file-proxy
A Reverse File Proxy

##Building

##Architecture
The main interface for this package is Proxy. Implementers should use that and a
ProxyBuilder to construct the proxy, then pass file names and sources in to
retrieve files from said sources. The Proxy will return a ReadCloser or an
error.

The sources for files implement the FileSource interface. This interface is
structured similarly to the http package's FileSystem, the main difference being
the Open method returns an io.ReadCloser instead of an http.File. This is done
so that an http.Response.Body does not need to wrapped in an implementation of
http.File. Since this interface only has a single Open method, it is fairly
trivial to add new types of file sources.

This package provides two FileSource implementations, one for the local file
system, and another for a remote file system. The local file system
implementation wraps the http.Dir implementation while the remote file system is
written from scratch. The path for the local and URL for the remote are
configurable.
