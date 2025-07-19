GO := go

BIN_DIR := bin

APP := cmd/app
MIGRATE := cmd/migrate
SEEDER := cmd/seeder

APP_BIN := $(BIN_DIR)/app
MIGRATE_BIN := $(BIN_DIR)/migrate
SEEDER_BIN := $(BIN_DIR)/seeder

# LDFLAGS for binary size optimization
LDFLAGS := -s -w -extldflags=-O2

.PHONY: all
all: build

.PHONY: build
build: $(APP_BIN) $(MIGRATE_BIN) $(SEEDER_BIN)

#BUILD APP backend
$(APP_BIN): $(APP)/*.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -tags=release -ldflags="$(LDFLAGS)" -trimpath -o $@ $<
	
#BUILD migrate
$(MIGRATE_BIN): $(MIGRATE)/*.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -tags=release -ldflags="$(LDFLAGS)" -trimpath -o $@ $<
	
#BUILD seeders
$(SEEDER_BIN): $(SEEDER)/*.go
	@mkdir -p $(BIN_DIR)
	$(GO) build -tags=release -ldflags="$(LDFLAGS)" -trimpath -o $@ $<

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
