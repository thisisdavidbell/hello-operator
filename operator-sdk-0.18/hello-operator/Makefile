# Makefile for operator

version=v0.0.5

.PHONY: check-registry-env-vars
check-registry-env-vars:
	@if [[ -z "${IRHOSTNAME}" ]] || [[ -z "${IRNAMESPACE}" ]] || [[ -z "${IRUSER}" ]] || [[ -z "${IRPASSWORD}" ]]; then \
		printf "\nERROR: Image Registry env vars not set. \nRequired:\n"; \
		printf "  - IRHOSTNAME - full hostname of image registry\n"; \
		printf "  - IRNAMESPACE - namespace in image registry\n"; \
		printf "  - IRUSER - user for image registry\n"; \
		printf "  - IRPASSWORD - password for image registry\n"; \
		exit 1; \
	fi

.PHONY: set-controller-registry
set-controller-registry: check-registry-env-vars
	sed -i.tmp 's/SET_TO_IRHOSTNAME/${IRHOSTNAME}/' pkg/controller/hello/hello_controller.go
	sed -i.tmp 's/SET_TO_IRNAMESPACE/${IRNAMESPACE}/' pkg/controller/hello/hello_controller.go
	rm pkg/controller/hello/hello_controller.go.tmp

.PHONY: reset-controller-registry
reset-controller-registry: check-registry-env-vars
	sed -i.tmp 's/${IRHOSTNAME}/SET_TO_IRHOSTNAME/' pkg/controller/hello/hello_controller.go
	sed -i.tmp 's/${IRNAMESPACE}/SET_TO_IRNAMESPACE/' pkg/controller/hello/hello_controller.go
	rm pkg/controller/hello/hello_controller.go.tmp

.PHONY: set-operatoryaml-registry
set-operatoryaml-registry: check-registry-env-vars
	sed -i.tmp 's/SET_TO_IRHOSTNAME/${IRHOSTNAME}/' deploy/operator.yaml
	sed -i.tmp 's/SET_TO_IRNAMESPACE/${IRNAMESPACE}/' deploy/operator.yaml
	sed -i.tmp 's/SET_TO_OPERATOR_VERSION/${version}/' deploy/operator.yaml
	rm deploy/operator.yaml.tmp

.PHONY: reset-operatoryaml-registry
reset-operatoryaml-registry: check-registry-env-vars
	sed -i.tmp 's/${IRHOSTNAME}/SET_TO_IRHOSTNAME/' deploy/operator.yaml
	sed -i.tmp 's/${IRNAMESPACE}/SET_TO_IRNAMESPACE/' deploy/operator.yaml
	sed -i.tmp 's/${version}/SET_TO_OPERATOR_VERSION/' deploy/operator.yaml
	rm deploy/operator.yaml.tmp

.PHONY: build-and-push-operator-as-is
build-and-push-operator-as-is: check-registry-env-vars
	docker login -u ${IRUSER} -p ${IRPASSWORD} ${IRHOSTNAME}
	operator-sdk build ${IRHOSTNAME}/${IRNAMESPACE}/hello-operator:${version}
	docker push ${IRHOSTNAME}/${IRNAMESPACE}/hello-operator:${version}

.PHONY: build-and-push-operator
build-and-push-operator: check-registry-env-vars set-controller-registry build-and-push-operator-as-is reset-controller-registry

.PHONY: deploy-operator-as-is
deploy-operator-as-is: 
	oc create -f deploy/service_account.yaml
	oc create -f deploy/role.yaml
	oc create -f deploy/role_binding.yaml
	oc create -f deploy/operator.yaml

.PHONY: redeploy-operator-as-is
redeploy-operator-as-is: 
	oc create -f deploy/operator.yaml

.PHONY: deploy-operator
deploy-operator: check-registry-env-vars set-operatoryaml-registry deploy-operator-as-is reset-operatoryaml-registry

.PHONY: clean-up
clean-up:
	-oc delete hello example-hello
	-oc delete deployment hello-operator

.PHONY: redeploy-operator
redeploy-operator: check-registry-env-vars set-operatoryaml-registry redeploy-operator-as-is reset-operatoryaml-registry

.PHONY: create-cr
create-cr: 
	oc create -f deploy/crds/thisisdavidbell.example.com_v1alpha1_hello_cr.yaml  