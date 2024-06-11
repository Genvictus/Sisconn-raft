
# Raft Protocol

 [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

We tried making [raft](https://raft.github.io/raft.pdf). "Raft is a consensus algorithm for managing a replicated log. It produces a result equivalent to (multi-)Paxos, and it is as efficient as Paxos, but its structure is different from Paxos; this makes Raft more understandable than Paxos and also provides a better foundation for building practical systems."

![Raft Protocol](./img/image.png "Raft Protocol")

Our raft system focuses on CP of CAP theorem (as it should be), but also has an increased degree of availability (the A in CAP) by accepting requests on bad condition, compensating it through delayed replication (if possible). The system allows pending request, but might get lost due to reasons (crashing, new elected leader overwrites the unreplicated request). The system also does not store its logs and states in stable storage (yet). Membership change is also not implemented, and might require a considerable amount of effort (emphasize on might).

## 🏹 Minimum Requirements

- Go 1.22
- Optional : Make

## ✨ How to Run

1. Clone this repository
2. Execute `go mod tidy` and then `go mod download` on terminal
3. Prepare the `.env` file in the project root directory. Use the [example](./.env.example) for reference
4. Run Makefile with `make server` to compile server and `make client` to compile client or using command `go build`
5. Run server with `./server` and client with `./client` with corresponding ip and port
6. Or just run with `go run ./cmd/server/` for server and `go run ./cmd/client/` for client.
7. The server takes `--host & --port args` to specify its address, and client takes `--host & --port args` for target address

Alternatively can run scripts in scripts folder to run server and client.

## 📚 How to Use

Here are some of the functionalities of the client:

- ping : to check if the server is alive.
- get `<key>` : to get the value of the key.
- set `<key>` `<value>` : to set the value of the key.
- strln `<key>` : to get the length of the string value of the key.
- del `<key>` : to delete the key.
- append `<key>` `<value>` : to append the value to the key.
- request-log : to get the log entries of the server.
- change-server `<host>` `<port>` : to change the server to connect to.
- help : to get the list of commands.
- start : to start batch transaction.
- exit : to exit the client.
- add-node `<host>` `<port>` : to add a node to the cluster.
- remove-node `<host>` `<port>` : to remove a node from the cluster.
- cluster-info : to get the information of the raft cluster.

---

## 📝 Contributors

| Nama | NIM |
|------|-----|
| [Yanuar Sano Nur Rasyid](https://github.com/yansans) | [13521110](https://github.com/noarotem) |
| [Ulung Adi Putra](https://github.com/Ulung32) | 13521122 |
| [Michael Utama](https://github.com/Michaelu670) | 13521137 |
| [Johann Christian Kandani](https://github.com/Genvictus) | 13521138 |
| [Dewana Gustavus Haraka Otang](https://github.com/DewanaGustavus) | 13521173 |
