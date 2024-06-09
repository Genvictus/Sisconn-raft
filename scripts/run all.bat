start "Server - localhost:6964" go run ./cmd/server/ --port 6964
start "Server - localhost:6965" go run ./cmd/server/ --port 6965
start "Server - localhost:6966" go run ./cmd/server/ --port 6966
start "Server - localhost:6967" go run ./cmd/server/ --port 6967
start "Server - localhost:6968" go run ./cmd/server/ --port 6968
start "Server - localhost:6969" go run ./cmd/server/ --port 6969

start "Client" go run ./cmd/client/

cd client-web
start "Client Web" bun dev

cd ..
cd dashboard-web
start "Dashboard Web" bun dev
