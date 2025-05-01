# Design principles of the actor model

- [] It's a graph: no supervisors;
- [x] Each node has an unique URI;
- [x] Node URIs can support multiple schemas, local to stay within the proces or remote, supporting different protocols, e.g. HTTP, TCP, etc;
- [X] Each node has a mailbox;
- [] Each node can emit messages;
- [] Each node can send a message to itself;
- [] Each node can receive messages;
- [] Each message has the sender and receiver address;
- [] Each node consumes messages from its mailbox;
- [] Messages can be just unicast;
- [] Each node can be a sender or receiver;
- [] Each node keeps track of the neighbors;
- [] Each node keeps track of the addresses of the entire graph;
- [] Each node has a companion which handles routing and node lifecycle;
- [X] Each node has a lifecycle: starting, running, stopping, idle;
- [] Each node can be started, stopped, restarted;
- [] Each node can fail, the companion can handle the failure: restart, stop, etc;
