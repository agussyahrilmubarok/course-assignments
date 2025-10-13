APP_NAME=point-service-batch
VERSION=0.0.1-SNAPSHOT
BUILD_DIR=build/libs
JAR_FILE=$(APP_NAME)-$(VERSION).jar
TARGET_DIR=/opt/batch-jobs
JAR_PATH=$(TARGET_DIR)/$(JAR_FILE)
SCRIPT_PATH=$(TARGET_DIR)/run-point-batch.sh
LOG_DIR=$(TARGET_DIR)/logs
CRON_CMD="0 1 * * * $(SCRIPT_PATH)"

.PHONY: all build install setup-cron clean

# Build the JAR using Gradle
build:
	./gradlew bootJar

# Copy JAR and setup directories
install: build
	sudo mkdir -p $(TARGET_DIR)
	sudo mkdir -p $(LOG_DIR)
	sudo cp $(BUILD_DIR)/$(JAR_FILE) $(JAR_PATH)
	sudo chmod +x $(JAR_PATH)

# Create the runner shell script
script:
	echo '#!/bin/bash' | sudo tee $(SCRIPT_PATH) > /dev/null
	echo '' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'LOG_FILE="$(LOG_DIR)/point-batch-$$(date +\%Y-\%m-\%d).log"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'JAR_PATH="$(JAR_PATH)"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'MAX_RETRIES=10' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'RETRY_DELAY=60' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'attempt=1' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'while [ $$attempt -le $$MAX_RETRIES ]; do' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  echo "Attempt $$attempt at $$(date)" >> "$$LOG_FILE"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  java -jar "$$JAR_PATH" >> "$$LOG_FILE" 2>&1' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  STATUS=$$?' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  if [ $$STATUS -eq 0 ]; then' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '    echo "Batch job succeeded on attempt $$attempt" >> "$$LOG_FILE"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '    break' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  else' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '    echo "Batch job failed on attempt $$attempt with status $$STATUS" >> "$$LOG_FILE"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '    ((attempt++))' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '    sleep $$RETRY_DELAY' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  fi' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'done' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'if [ $$attempt -gt $$MAX_RETRIES ]; then' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo '  echo "Batch job failed after $$MAX_RETRIES attempts" >> "$$LOG_FILE"' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	echo 'fi' | sudo tee -a $(SCRIPT_PATH) > /dev/null
	sudo chmod +x $(SCRIPT_PATH)

# Register cron job
setup-cron:
	# Hapus baris cron sebelumnya (jika ada)
	(crontab -l 2>/dev/null | grep -v "$(SCRIPT_PATH)"; echo $(CRON_CMD)) | crontab -

# Full setup
all: install script setup-cron

# Remove installed files
clean:
	sudo rm -rf $(TARGET_DIR)
	crontab -l | grep -v "$(SCRIPT_PATH)" | crontab -
