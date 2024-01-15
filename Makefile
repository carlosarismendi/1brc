BIN_FILE=1brc

MEASUREMENTS_FOLDER=measurements
TEST_AMOUNT=1k
TEST_MEASUREMENTS_FILE=$(MEASUREMENTS_FOLDER)/measurements_$(TEST_AMOUNT).txt
TEST_ACTUAL_RESULTS_FILE=$(MEASUREMENTS_FOLDER)/actual_results_$(TEST_AMOUNT).txt
TEST_EXPECTED_RESULTS_FILE=$(MEASUREMENTS_FOLDER)/expected_results_$(TEST_AMOUNT)_mine.txt

AMOUNT=200m
MEASUREMENTS_FILE=$(MEASUREMENTS_FOLDER)/measurements_$(AMOUNT).txt
ACTUAL_RESULTS_FILE=$(MEASUREMENTS_FOLDER)/actual_results_$(AMOUNT).txt
EXPECTED_RESULTS_FILE=$(MEASUREMENTS_FOLDER)/expected_results_$(AMOUNT).txt

ARGS=-measurements-file=$(MEASUREMENTS_FILE) -max-workers=8 -max-ram=6

.PHONY: build
build: clean
	go build -o $(BIN_FILE)

.PHONY: time
time: build	
	/usr/bin/time -v ./$(BIN_FILE) $(ARGS)

.PHONY: run
run: build
	./$(BIN_FILE) $(ARGS) > $(ACTUAL_RESULTS_FILE)

.PHONY: profile
profile: build
	./$(BIN_FILE) $(ARGS) -cpuprofile=1brc.prof
	echo "top15" | go tool pprof 1brc 1brc.prof

.PHONY: test
test: AMOUNT=$(TEST_AMOUNT)
test: MEASUREMENTS_FILE=$(TEST_MEASUREMENTS_FILE)
test: ACTUAL_RESULTS_FILE=$(TEST_ACTUAL_RESULTS_FILE)
test: EXPECTED_RESULTS_FILE=$(TEST_EXPECTED_RESULTS_FILE)
test: clean run
	diff -b -B -i $(ACTUAL_RESULTS_FILE) $(EXPECTED_RESULTS_FILE)

.PHONY: clean
clean:
	rm -rf $(MEASUREMENTS_FOLDER)/actual* $(BIN_FILE)
