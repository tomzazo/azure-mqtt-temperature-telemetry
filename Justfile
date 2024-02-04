@build:
    just ./src/build

@setup:
	ansible-playbook -k -i ./ansible/inventory.yml ./ansible/playbooks/setup.yml
