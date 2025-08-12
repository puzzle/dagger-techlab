#!/bin/bash

set -euo pipefail

# Required environment variables
: "${LOGIN_TOKEN:?Missing LOGIN_TOKEN}"

# github.repository.name as <account>
#export ACCOUNT='${{ github.event.repository.name }}'

export PACKAGE_NAME='dagger-techlab'
#export PR_NUMBER='${{ github.event.pull_request.number }}'
export PACKAGE_TYPE='container'
export ORG='puzzle'

export VERSIONS_FILE=versions.json
export SIGNATURES_FILE=versions-with-sigs.json

# get all packages: images and signatures (versions)
curl -L -f \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${LOGIN_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  https://api.github.com/orgs/${ORG}/packages/${PACKAGE_TYPE}/${PACKAGE_NAME}/versions > ${VERSIONS_FILE}


# Extract all names (SHAs) from the file
all_names=$(jq -r '.[].name' versions.json)

echo "all_names= $all_names"

# Loop through each entry in the JSON file
jq -c '.[]' versions.json | while read -r entry; do
    # Extract the SHA tag from metadata.container.tags
 #  tag=$(echo "$entry" | jq -r '.metadata.container.tags[] | select(startswith("sha256:"))')
    full_tag=$(echo "$entry" | jq -r '.metadata.container.tags[] | select(test("sha256-.*sig"))')
    echo "full_tag: $full_tag"
    tag=$(sed 's/sha256-/sha256:/; s/\.sig$//' <<< "$full_tag")
    echo "tag: $tag"


    # If a SHA tag exists, check if it's in the list of all names
    if [[ -n "$tag" ]]; then
        # Check if the tag is not in the list of all names
        if ! grep -q "$tag" <<< "$all_names"; then
            # If it's not found, print the name of the entry
            echo "NOOOT FOUND"
            echo "Dangling Sig to delete: $full_tag"
            sig_sha=$(cat versions.json | jq ".[] | select(.metadata.container.tags[]? == \"${full_tag}\") | .name")
            echo "sig_sha: $sig_sha"
        else
            echo "FOUND"
            echo "$entry" | jq -r '.name'
        fi
    else
        echo "NO sig tag"
    fi
    echo ""
done

echo "🏁 Registry cleanup completed."