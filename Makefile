BIN_FILE=1brc

TEST_AMOUNT=1k
TEST_MEASUREMENTS_FILE=measurements_$(TEST_AMOUNT).txt
TEST_ACTUAL_RESULTS_FILE=actual_results_$(TEST_AMOUNT).txt
TEST_EXPECTED_RESULTS_FILE=expected_results_$(TEST_AMOUNT).txt

AMOUNT=50m
MEASUREMENTS_FILE=measurements_$(AMOUNT).txt
ACTUAL_RESULTS_FILE=actual_results_$(AMOUNT).txt
EXPECTED_RESULTS_FILE=expected_results_$(AMOUNT).txt

.PHONY: build
build:
	go build -o $(BIN_FILE)

.PHONY: time
time: build	
	MEASUREMENTS_FILE=$(MEASUREMENTS_FILE) /usr/bin/time -v ./$(BIN_FILE) 

.PHONY: run
run: build
	MEASUREMENTS_FILE=$(MEASUREMENTS_FILE) ./$(BIN_FILE) > $(ACTUAL_RESULTS_FILE)

.PHONY: profile
profile: build
	MEASUREMENTS_FILE=$(MEASUREMENTS_FILE) ./$(BIN_FILE) -cpuprofile=1brc.prof
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
	rm -rf $(ACTUAL_RESULTS_FILE) $(BIN_FILE)
