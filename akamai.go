package ipcat

import (
    "encoding/json"
)

type Akamai struct {
    Addresses  []string `json:"addresses"`
}

// The IP ranges list
func DownloadAkamai() ([]byte, error) {
    //  Ref: Origin IP Access Control List
    //  https://techdocs.akamai.com/property-mgr/docs/origin-ip-access-control
    const body = `{"addresses": ["23.32.0.0/11","23.192.0.0/11","2.16.0.0/13","104.64.0.0/10","184.24.0.0/13","23.0.0.0/12","95.100.0.0/15","92.122.0.0/15","172.232.0.0/13","184.50.0.0/15","88.221.0.0/16","23.64.0.0/14","72.246.0.0/15","96.16.0.0/15","96.6.0.0/15","69.192.0.0/16","23.72.0.0/13","173.222.0.0/15","118.214.0.0/16","184.84.0.0/14"]}`
    return []byte(body), nil
}

// Parses the IP ranges json file and updates the interval set
func UpdateAkamai(ipmap *IntervalSet, body []byte) error {
    const (
        dcName = "Akamai"
        dcURL  = "http://akamai.com/"
    )

    akamai := Akamai{}
    err := json.Unmarshal(body, &akamai)
    if err != nil {
        return err
    }

    // Delete all existing records
    ipmap.DeleteByName(dcName)

    // Now add the IP ranges to the map
    for _, cidr := range akamai.Addresses {
        err := ipmap.AddCIDR(cidr, dcName, dcURL)
        if err != nil {
            return err
        }
    }

    return nil
}
