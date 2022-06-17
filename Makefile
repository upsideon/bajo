MOCKS_PATH=mocks
.PHONY: mocks

mocks: database.go
	@echo "Generating mocks..."
	@rm -rf $(MOCKS_PATH)
	@mockgen -source=database.go -destination=$(MOCKS_PATH)/database_mock.go
