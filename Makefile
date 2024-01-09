BIN_FILE=1brc
AMOUNT=1k
MEASUREMENTS_FILE=measurements_$(AMOUNT).txt
ACTUAL_RESULTS_FILE=actual_results_$(AMOUNT).txt
EXPECTED_RESULTS_FILE=expected_results_$(AMOUNT).txt

.PHONY: build
build:
	go build -o $(BIN_FILE)

.PHONY: run
time: build
	export MEASUREMENTS_FILE=$(MEASUREMENTS_FILE) \
		&& time ./$(BIN_FILE) > $(ACTUAL_RESULTS_FILE)

.PHONY: run
run: build
	export MEASUREMENTS_FILE=$(MEASUREMENTS_FILE) \
		&& ./$(BIN_FILE) > $(ACTUAL_RESULTS_FILE)

.PHONY: test
test: clean run
	diff -b -B -i $(ACTUAL_RESULTS_FILE) $(EXPECTED_RESULTS_FILE)

.PHONY: clean
clean:
	rm -rf $(ACTUAL_RESULTS_FILE) $(BIN_FILE)
