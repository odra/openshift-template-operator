SHELL = /bin/bash
REG = docker.io
ORG = odranoel
IMAGE_NAME = openshift-template-operator
IMAGE_TAG = latest

dep/ensure:
	@dep ensure -v

image/build:
	@operator-sdk build ${REG}/${ORG}/${IMAGE_NAME}:${IMAGE_TAG}

image/push:
	@docker push ${REG}/${ORG}/${IMAGE_NAME}:${IMAGE_TAG}

cluster/prepare:
	@oc apply -f deploy/service_account.yaml
	@oc apply -f deploy/role.yaml
	@oc apply -f deploy/role_binding.yaml
	@oc apply -f deploy/crds/odra_v1alpha1_okdtemplate_crd.yaml

cluster/deploy:
	@oc apply -f deploy/crds/odra_v1alpha1_okdtemplate_cr.yaml
	sleep 5
	@oc apply -f deploy/operator.yaml