# Design principles of the actor model

## Actors

- [x] Each node has an unique URI;
- [] Node URIs can support multiple schemas, 'actor' to stay within the proces or remote, supporting different protocols, e.g. HTTP, TCP, etc;
- [X] Each node has a mailbox;
- [X] Each node can emit messages;
- [X] Each node can send a message to itself;
- [X] Each node can receive messages;
- [X] Each message has the sender and receiver address;
- [X] Each node consumes messages from its mailbox;
- [X] Each node can be configure for backpressure policies;
- [X] Messages can be just unicast;
- [X] Each node can be a sender or receiver;
- [X] Each node has a lifecycle: running, idle;
- [] Each node can be restarted;
- [] Handling failures and policies retry, backoff, drop and restart the node...
- [] Improve the Message model;
- [] Logging.
