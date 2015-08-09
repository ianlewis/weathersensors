# Aggre Mod

The aggregation module runs on a Raspberry Pi. The Raspberry Pi also runs an
instance of the spark-server which is used to register the devices and push
firmware updates.

## Provision the Raspberry Pi

The Raspberry Pi is provisioned using Ansible.

    $ ansible-playbook --inventory provisioning/inventory.ini provisioning/site.yml
