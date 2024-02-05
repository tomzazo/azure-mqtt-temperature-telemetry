@build:
    docker compose up build
    docker compose down

@deploy:
	ansible-playbook -k -i ./ansible/inventory.yml ./ansible/playbooks/setup.yml

@setup:
    just build deploy
