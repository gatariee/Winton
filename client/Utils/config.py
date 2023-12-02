import yaml
import os

def load(config='./config.yml'):
    try:
        with open(config, 'r') as f:
            config_data = yaml.safe_load(f)
            return config_data
    except FileNotFoundError:
        print("[!] Error: Config file not found, looking for 'config.yml'")
        os._exit(1)