#!/usr/bin/env bash

biggest_lists() {
	find ${1:-.} -maxdepth 1 -type f -name "*.*set" -printf "%s,%p\n" | sort -nr | head -50 | cut -d',' -f2 | tr '\n' ' ';
}

find_lists() {
	find ${2:-.} -type f -name "*${1}*.*set" -printf "%p ";
}

sanitize_list() {
	sed -i '/^0/d' ${1} && printf "\\nRemoved lines/IPs starting with 0\\n";
}

# ensure we got the basic lists
FIREHOL_FOLDER=${1:-$(mktemp -d)}
FIREHOL_REPO="https://github.com/firehol/blocklist-ipsets"
cd "$FIREHOL_FOLDER";
printf "\\nFetching/Updating firehol lists..\\n";
set -x
git remote -v || (rm -fr *; git clone --progress "$FIREHOL_REPO" "$FIREHOL_FOLDER";);
git pull;
cd -;

# fetch extras
EXTRA_TMP=$(mktemp)
for extra_list in $(wget -qO- https://raw.githubusercontent.com/trick77/ipset-blacklist/master/ipset-blacklist.conf | grep -oE '\"http.*\"' | sed 's/"//g'); do
	wget -qO- "${extra_list}" >> "$EXTRA_TMP";
done

# check for trim4
command -v trimv4 || go get -v github.com/nexus166/trimv4 || exit 127;

set +x
ALL_LISTS="";

# banning entire continents actually helps your CPUs
for continent in continent_af continent_as continent_na continent_oc continent_sa; do
        ALL_LISTS+=$(find_lists "$continent" ${FIREHOL_FOLDER});
done

# or just countries
ALL_LISTS+=$(find_lists "country_ru" ${FIREHOL_FOLDER});

# include top 50 lists from firehol
ALL_LISTS+=$(biggest_lists ${FIREHOL_FOLDER});

# look for other lists we care about in the firehol folder
ALL_LISTS+=$(find_lists "tor_exits" ${FIREHOL_FOLDER});

# include trick77/ipset-blacklist list of lists
ALL_LISTS+="$EXTRA_TMP";

# done. start processing final output.
printf "\\nLists that will be considered:\\n%s\\n\\nStatus:\\t%s\\n" "${ALL_LISTS}" "$(cat $ALL_LISTS | wc -l)";

# compute final list
FINAL_LIST=$(mktemp)
printf "\\nRunning trimv4 against complete list.. \\r";
cat ${ALL_LISTS} | ${GOPATH}/bin/trimv4 - > "$FINAL_LIST" && printf "Running trimv4 against complete list.. Done.\\r\\n";

# sanitize it
sanitize_list "$FINAL_LIST"

printf "\\nStatus: %s\\n" "$(wc -l $FINAL_LIST)"

# create ipset restore file
printf "create blacklist-tmp -exist hash:net family inet\\ncreate blacklist -exist hash:net family inet\\n" > ./blacklist.restore
for _line in $(< ${FINAL_LIST}); do
	printf "add blacklist-tmp %s\\n" "${_line}" >> ./blacklist.restore;
done
printf "swap blacklist blacklist-tmp\\ndestroy blacklist-tmp\\n" >> ./blacklist.restore

printf "\\n%s is ready.\\n" "$(realpath blacklist.restore)"
mv -fv "$FINAL_LIST" "$(dirname blacklist.restore)/blacklist.txt"
