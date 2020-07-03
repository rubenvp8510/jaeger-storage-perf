.PHONY: badger
badger:
	go test -bench=. -storage-type=badger -count=5 > badger

.PHONY: redbull
redbull:
	go test -bench=. -storage-type=redbull -count=5 > redbull


.PHONY: questdb
questdb:
	go test -bench=. -storage-type=questdb -count=5 > questdb

.PHONY: druid
druid:
	go test -bench=. -storage-type=druid -count=5 > druid