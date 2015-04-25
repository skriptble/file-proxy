# file-proxy
A Reverse File Proxy

##Building

##Architecture
There is a top level interface that is used by implementers. It should take the
file name and the destination and return a file. This interface is used between
the HTTP proxy server and the backend.

The backend implementations implement the interface mentioned above. There
should be a single interface for both local and remote files. This allows
additional file systems to be added without needing to change much code, if any,
on the HTTP proxy server side.

Below this top level interface should be the interfaces that handle actually
retrieving local or remote files. The implementations of these interfaces
register with the implementation of the top level interface, filling a
particular type of file system.
