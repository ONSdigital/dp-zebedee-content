# Source directories
SRC_DIR          = src
SRC_CONTENT_DIR  = content
SRC_SERVICES_DIR = services

# Generated content directory
GENERATED_ROOT_DIR = generated

# Publishing dirs
# PUB_ROOT_DIR is $zebedee_root + the extra zebedee subdir that zebedee requires
PUB_ROOT_DIR       = $(GENERATED_ROOT_DIR)/publishing/zebedee
PUB_MASTER_DIR     = master
PUB_SERVICES_DIR   = services
PUB_DIRS_TO_CREATE = $(PUB_MASTER_DIR) $(PUB_SERVICES_DIR) application-keys collections keyring launchpad permissions publishing-log transactions

# Web dirs
WEB_ROOT_DIR       = $(GENERATED_ROOT_DIR)/web
WEB_SITE_DIR       = site
WEB_DIRS_TO_CREATE = $(WEB_SITE_DIR) transactions

.PHONY: init
init: clean
	# Init publishing
	for dir in $(PUB_DIRS_TO_CREATE); do mkdir -p $(PUB_ROOT_DIR)/$$dir; done
	cp -r $(SRC_DIR)/$(SRC_CONTENT_DIR)/* $(PUB_ROOT_DIR)/$(PUB_MASTER_DIR)/
	cp -r $(SRC_DIR)/$(SRC_SERVICES_DIR)/* $(PUB_ROOT_DIR)/$(PUB_SERVICES_DIR)/

	# Init web
	for dir in $(WEB_DIRS_TO_CREATE); do mkdir -p $(WEB_ROOT_DIR)/$$dir; done
	cp -r $(SRC_DIR)/$(SRC_CONTENT_DIR)/* $(WEB_ROOT_DIR)/$(WEB_SITE_DIR)/

.PHONY: clean
clean:
	rm -rf $(GENERATED_ROOT_DIR)
