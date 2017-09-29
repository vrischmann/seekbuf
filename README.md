seekbuf
=======

This is a seekable read-write buffer.

It isn't as optimized as `bytes.Buffer` is so only use this if you really need a buffer which implements `io.Seeker`.
