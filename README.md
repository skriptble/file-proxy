# file-proxy
A Reverse File Proxy

##Dependencies
The Froxy library depends only on the Go standard library. Go version 1.4 and
above is suggested. The Froxy command uses spf13's Cobra and Viper, both of
which have been vendored into the cmd/foxy directory.

##Installing and Building

##Testing
Testing from a third party service can be achieved by starting this service, and
pointing the remote endpoint at a server the third party controls or knows the
contents of. The third party may then execute several requests to froxy,
evaluating the responses are correct.

Testing the local proxy functionality requires the service to be running in an
environment under the control of the tester. The tester should place files
somewhere in the local environment, start the froxy service, pointing the local
endpoint at the local environment. The tester may then execute several tests
from another service or machine and validate the responses are correct.

Unit and Functional tests have been written for each component of the library.
These tests cover most of the paths through each individual component.
Integration tests between components have not been written.

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
