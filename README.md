# dist_system_practices

1- basic json parser

Message-Based Communication

Distributed systems communicate through messages because nodes cannot share memory. This is a fundamental constraint that shapes how we design distributed algorithms.
Why Messages?

In a single-machine program, threads can share memory directly. But in a distributed system:

    Nodes are on different machines - they have separate memory spaces
    Networks are unreliable - messages can be delayed, duplicated, or lost
    Failures are partial - some nodes may crash while others continue

Each message must be self-contained with enough information for the recipient to process it independently.
The Maelstrom Protocol

Maelstrom uses a simple JSON-based protocol with three required fields:

    src - identifies who sent the message
    dest - identifies the intended recipient
    body - contains the actual payload with a type field

Why stdin/stdout?

Using standard streams makes the protocol language-agnostic. Maelstrom can spawn your binary and communicate with it regardless of what language you wrote it in. This same pattern is used by many real systems for inter-process communication (IPC).
