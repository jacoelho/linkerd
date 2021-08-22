package linkerd

// This package implements the Linkerd propagator specification as defined
// at https://github.com/twitter/finagle/blob/develop/finagle-core/src/main/scala/com/twitter/finagle/tracing/TraceId.scala#L35

// Format as documented:
// For backward compatibility for TraceID: 40 bytes if 128bit, 32 bytes if 64bit
// val bytes = new Array[Byte](if (traceId.traceIdHigh.isDefined) 40 else 32)
// ByteArrays.put64be(bytes, 0, traceId.spanId.toLong)
// ByteArrays.put64be(bytes, 8, traceId.parentId.toLong)
// ByteArrays.put64be(bytes, 16, traceId.traceId.toLong)
// ByteArrays.put64be(bytes, 24, flags.toLong)
// if (traceId.traceIdHigh.isDefined) ByteArrays.put64be(bytes, 32, traceId.traceIdHigh.get.toLong)