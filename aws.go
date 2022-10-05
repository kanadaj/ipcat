package ipcat

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type AWSPrefix struct {
    IPPrefix string `json:"ip_prefix"`
    Region   string `json:"region"`
    Service  string `json:"service"`
}

type AWS struct {
    SyncToken  string      `json:"syncToken"`
    CreateDate string      `json:"createDate"`
    Prefixes   []AWSPrefix `json:"prefixes"`
}

// Downloads the latest AWS IP ranges list
func DownloadAWS() ([]byte, error) {
    //  Ref: AWS IP address ranges
    //  https://docs.aws.amazon.com/general/latest/gr/aws-ip-ranges.html
    const awsDownload = "https://ip-ranges.amazonaws.com/ip-ranges.json"

    resp, err := http.Get(awsDownload)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("Failed to download AWS ranges: status code %s", resp.Status)
    }
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    resp.Body.Close()
    return body, nil
}

// Parses the AWS IP json file and updates the interval set
func UpdateAWS(ipmap *IntervalSet, body []byte) error {
    const (
        dcName = "Amazon AWS"
        dcURL  = "http://www.amazon.com/aws/"
    )

    aws := AWS{}
    err := json.Unmarshal(body, &aws)
    if err != nil {
        return err
    }

    // delete all existing records
    ipmap.DeleteByName(dcName)

    // and add back
    for _, block := range aws.Prefixes {
        if block.Service == "EC2" {
            err := ipmap.AddCIDR(block.IPPrefix, dcName, dcURL)
            if err != nil {
                return err
            }
        }
    }

    return nil
}
