package ipcat

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type Fastly struct {
    IpV6Addresses  []string `json:"ipv6_addresses"`
    Addresses      []string `json:"addresses"`
}

// Downloads the latest IP ranges list
func DownloadFastly() ([]byte, error) {
    //  Ref: Fastly IP address ranges
    //  https://developer.fastly.com/reference/api/utils/public-ip-list/
    const url = "https://api.fastly.com/public-ip-list"

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to download IP ranges: status code %s", resp.Status)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    resp.Body.Close()
    return body, nil
}

// Parses the IP ranges json file and updates the interval set
func UpdateFastly(ipmap *IntervalSet, body []byte) error {
    const (
        dcName = "Fastly"
        dcURL  = "https://www.fastly.com/"
    )

    fastly := Fastly{}
    err := json.Unmarshal(body, &fastly)
    if err != nil {
        return err
    }

    // Delete all existing records
    ipmap.DeleteByName(dcName)

    // Now add the IP ranges to the map
    for _, cidr := range fastly.Addresses {
        err := ipmap.AddCIDR(cidr, dcName, dcURL)
        if err != nil {
            return err
        }
    }

    return nil
}
