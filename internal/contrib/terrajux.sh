#!/usr/bin/env bash

# this was the original PoC script that terrajux replaced

# usage: terrajux <giturl> <subpath> <ref> <ref>
# ex: terrajux https://github.com/terraform-aws-modules/terraform-aws-iam.git /modules/iam-user 4.1.0 4.2.0
# clean cache with: TJUX_CLEAN=1 terrajux <args...>

set -o errexit -o nounset #-o xtrace

if [[ "$#" -ne 4 ]] ; then
    printf '\nterrajux diffs the source code of a terraform project and its transitive dependencies.\n'
    printf 'usage: [TJUX_CLEAN=1] %s giturl subpath ref1 ref2\n\n' "$(basename "$0")"
    exit 1
fi

giturl="${1}"
subpath="${2}"
v1="${3}"
v2="${4}"

tdir="${HOME}/.terrajux"
tcache="${tdir}/cache"
gdir=$(basename "${giturl}")
cleanup="${TJUX_CLEAN:-}"


# begin

[[ -n $cleanup ]] && rm -rf "${tcache}"

mkdir -p "${tcache}"
pushd "${tcache}"

for v in "${v1}" "${v2}" ; do
    if [[ ! -d "${gdir}-${v}" ]] ; then
        git clone --quiet --depth=1 --shallow-submodules -b "${v}" "${giturl}" "${gdir}-${v}"
        pushd "${gdir}-${v}${subpath}"
        terraform init -backend=false
        popd
    fi
done

#diff -rybBNEW "${maxw}" --suppress-common-lines "${v1}${subpath}" "${v2}${subpath}" | less

git \
    # apparently this is git's way of saying false
    -c advice.detachedHead= \
    diff \
        --color=auto \
        --no-ext-diff \
        --no-index \
        --ignore-all-space \
        "${gdir}-${v1}${subpath}" \
        "${gdir}-${v2}${subpath}"

popd

# end
