"""Get Sigma rules from various git repositories"""

import json
import logging
import sys
from subprocess import check_output, call

logging.basicConfig(format=r"%(asctime)s %(message)s", level=logging.INFO)


def check(command):
    """Wrapper around subprocess.check_ouptut"""
    return check_output(command, shell=True)


def get_rules():
    """Load in rule sources config file"""

    try:
        rule_sources = json.loads(sys.argv[1])
    except IndexError:
        with open("rule_sources.json") as rule_sources:
            rule_sources = json.load(rule_sources)
    return rule_sources


def main():
    """Main"""
    sync_dirs = []
    cleanup_dirs = []
    for source in get_rules():
        try:
            repo = source["repo"]
        except IndexError:
            logging.error("Missing field 'repo'")
            continue

        branch = source.get("branch", "master")
        path = source.get("path", "")

        # Clone
        call(f"git clone {repo} > /dev/null 2>&1", shell=True)
        repo_dir = repo.split("/")[-1].split(".")[0]
        check(
            f"cd {repo_dir} && git checkout {branch} > /dev/null 2>&1 && cd {'../'*len(repo_dir.split('/'))}")

        # Save for sync and cleanup outside of for loop
        cleanup_dirs.append(repo_dir)
        sync_dirs.append(f"{repo_dir}/{path}")

    # Sync all sync_dirs together into ./rules
    check("rm -rf rules && mkdir rules")
    check(f"rsync -rtvu {' '.join(sync_dirs)} rules")

    # Move rules into build dir
    logging.info("Moving rules to engine build directory")
    build_dir = "../images/rules/"
    check(f"rm -rf {build_dir}")
    check(f"mv rules {build_dir}")

    # Remove directories
    logging.info("Cleaning up")
    cleanup_dirs.append("rules")
    for directory in cleanup_dirs:
        logging.info(f"Removing {directory}")
        check(f"rm -rf {directory}")


if __name__ == "__main__":
    main()
