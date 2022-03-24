.SILENT: status update_rules build status apply deploy destroy test run_test clean_test

.PHONY: status update_rules build status apply deploy destroy test run_test clean_test

# DEFINE CONSTANTS
MAKE_PROJ_NAME = sigma-engine
MAKE_IMG_TAG = latest
MAKE_BUILD_DIR = images/
MAKE_HELM_ENVIRONMENT = default
MAKE_K8S_NAMESPACE = se
MAKE_TEST_JOB_NAME = integration-test
MAKE_TEST_BUILD_DIR = images/integration_test/

# COMPUTED CONSTANTS
MAKE_IMG_NAME = $(MAKE_K8S_NAMESPACE)/$(MAKE_PROJ_NAME)
MAKE_HELM_OPTIONS = --namespace $(MAKE_K8S_NAMESPACE) --environment $(MAKE_HELM_ENVIRONMENT)
MAKE_KUBE_OPTIONS = --namespace $(MAKE_K8S_NAMESPACE)

status:
	echo " > Make Status"
	helmfile $(MAKE_HELM_OPTIONS) status

update_rules:
	echo " > Make Update_Rules"
	# Pull rules required for building of Sigma-Engine
	cd helpers && python3 update_rules.py

build: update_rules
	echo " > Make Build"
	# Build via Docker-Compose
	echo " > Building via docker-compose: $(MAKE_BUILD_DIR)docker-compose.yaml"
	docker-compose --file $(MAKE_BUILD_DIR)docker-compose.yaml build --no-cache
	echo " > Build complete"

create_namespace:
	echo " > Creating se namespace"
	- kubectl create namespace se

apply:
	echo " > Make Apply"
	helmfile --file deploy/helmfile.kafka.yaml $(MAKE_HELM_OPTIONS) apply

deploy: create_namespace apply

destroy: clean_test
	echo " > Make Destroy"
	helmfile --file deploy/helmfile.kafka.yaml $(MAKE_HELM_OPTIONS) destroy

run_test:
	echo " > Make Test"
	# Execute helper manually for now
	@read -p "Type (grpc/kafka/http)? : " JOB_TYPE \
	&& helpers/run_integration_test.sh $(MAKE_K8S_NAMESPACE) $(MAKE_TEST_JOB_NAME) \
	   deploy/resources/$${JOB_TYPE}_integration_test_job.yaml

clean_test:
	echo " > Clean Test"
	if kubectl $(MAKE_KUBE_OPTIONS) get jobs | grep -q $(MAKE_TEST_JOB_NAME) && True; then \
	  kubectl $(MAKE_KUBE_OPTIONS) delete job $(MAKE_TEST_JOB_NAME); \
	fi

test: clean_test run_test
