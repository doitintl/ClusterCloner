set -x
locs=(eastus eastus2 westus centralus northcentralus southcentralus northeurope westeurope eastasia southeastasia japaneast japanwest australiaeast australiasoutheast australiacentral brazilsouth southindia centralindia westindia canadacentral canadaeast westus2 westcentralus uksouth ukwest koreacentral koreasouth francecentral southafricanorth uaenorth switzerlandnorth germanywestcentral norwayeast)
for i in ${!locs[@]}; do
  az vm list-sizes --location ${locs[$i]} | jq -r '(map(keys) | add | unique) as $cols | map(. as $row | $cols | map($row[.])) as $rows | $cols, $rows[] | @csv' >sizes_${locs[$i]}.csv
done
