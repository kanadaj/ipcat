package ipcat

import (
    "bytes"
    "fmt"
    "strings"
    "io/ioutil"
    "net/http"
)

// Downloads the latest IP ranges list
func DownloadDigitalOcean() ([]byte, error) {
    //  Ref: Digital Ocean Platform Information
    //  https://docs.digitalocean.com/products/platform/#platform-information
    const url = "https://www.digitalocean.com/geo/google.csv"

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

    return bytes.TrimSpace(body), nil
}

// Parses the IP ranges json file and updates the interval set
func UpdateDigitalOcean(ipmap *IntervalSet, body []byte) error {
    const (
        dcName = "DigitalOcean"
        dcURL  = "https://www.digitalocean.com/"
    )

    // Delete all existing records
    ipmap.DeleteByName(dcName)

    // Body is CSV format with 5 columns per row
    // CIDR,Country,Region,City,PostalCode

    // Now add the IP ranges to the map
    for _, row := range bytes.Split(body, []byte("\n")) {
        cols := strings.Split(string(row), ",")
        if len(cols) != 5 {
            return fmt.Errorf("Row does not have expected columns: row = %s", string(row))
        }
        cidr := cols[0]
        // Only add IPv4 prefixes
        if strings.Count(cidr, ":") == 0 {
            err := ipmap.AddCIDR(cidr, dcName, dcURL)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
