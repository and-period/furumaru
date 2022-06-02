.PHONY: init plan apply

init:
	cd ${DIR}; terraform init

plan:
	cd ${DIR}; terraform validate
	cd ${DIR}; terraform plan -var-file='./../99_common.tfvars' -parallelism=120

apply:
	cd ${DIR}; terraform apply -var-file='./../99_common.tfvars' -parallelism=120
