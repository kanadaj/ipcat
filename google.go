package ipcat

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

//  Ref: Google Cloud
//  https://cloud.google.com/compute/docs/faq#find_ip_range
var (
    googleDownload = "https://www.gstatic.com/ipranges/cloud.json"
)

type GooglePrefix struct {
    IPPrefix string `json:"ipv4Prefix"`
    Service  string `json:"service"`
    Scope    string `json:"scope"`
}

type GoogleCloud struct {
    SyncToken  string         `json:"syncToken"`
    CreateDate string         `json:"creationTime"`
    Prefixes   []GooglePrefix `json:"prefixes"`
}

// Downloads the latest Google Cloud public IP ranges list
func DownloadGoogle() ([]byte, error) {
    resp, err := http.Get(googleDownload)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to download Google Cloud ranges: status code %s", resp.Status)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    resp.Body.Close()
    return body, nil
}

// Parses the Google Cloud IP json file and updates the interval set
func UpdateGoogle(ipmap *IntervalSet, body []byte) error {
    const (
        dcName = "Google Cloud"
        dcURL  = "https://cloud.google.com/"
    )

    googleCloud := GoogleCloud{}
    err := json.Unmarshal(body, &googleCloud)
    if err != nil {
        return err
    }

    // delete all existing records
    ipmap.DeleteByName(dcName)

    // and add back
    for _, block := range googleCloud.Prefixes {
        if len(block.IPPrefix) > 0 {
            err := ipmap.AddCIDR(block.IPPrefix, dcName, dcURL)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
