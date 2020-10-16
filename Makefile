.PHONY: generate
generate:
	go run generate.go

.PHONY: badger
badger:
	go test -bench=. -storage-type=badger -count=10 > badger

.PHONY: redbull
redbull:
	go test -bench=. -storage-type=redbull -count=10

.PHONY: questdb
questdb:
	go test -bench=. -storage-type=questdb -count=10 > questdb

.PHONY: druid
druid:
	go test -bench=. -storage-type=druid -count=10 > druid

.PHONY: blob_write
blob_write:
	go test -bench=BlobStorageWrite -count=10 -timeout=0

.PHONY: blob_read
blob_read:
	go test -bench=BlobStorageRead -count=20 -timeout=0
