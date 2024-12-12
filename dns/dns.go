package dns

import (
	"fmt"
	"sync"
	"time"

	"ddns-ipv6/config"

	"github.com/cenkalti/backoff/v4"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dnspod/v20210323"
)

type DNSCache struct {
	CurrentIP  string
	LastUpdate time.Time
	sync.RWMutex
}

func NewDNSCache() *DNSCache {
	return &DNSCache{}
}

func (c *DNSCache) UpdateIP(ip string) {
	c.Lock()
	defer c.Unlock()
	c.CurrentIP = ip
	c.LastUpdate = time.Now()
}

func (c *DNSCache) GetIP() (string, time.Time) {
	c.RLock()
	defer c.RUnlock()
	return c.CurrentIP, c.LastUpdate
}

// UpdateDNSRecord 更新域名解析记录
func UpdateDNSRecord(client *dnspod.Client, config config.Config, ipv6 string) error {
	// 获取记录列表以找到需要更新的记录ID
	listRequest := dnspod.NewDescribeRecordListRequest()
	listRequest.Domain = common.StringPtr(config.Domain.Domain)
	listRequest.Subdomain = common.StringPtr(config.Domain.SubDomain)

	listResponse, err := client.DescribeRecordList(listRequest)
	if err != nil {
		return err
	}

	var recordID *uint64
	for _, record := range listResponse.Response.RecordList {
		if *record.Type == "AAAA" && *record.Name == config.Domain.SubDomain {
			recordID = record.RecordId
			break
		}
	}

	if recordID == nil {
		return fmt.Errorf("no matching AAAA record found for subdomain %s", config.Domain.SubDomain)
	}

	// 更新记录
	modifyRequest := dnspod.NewModifyRecordRequest()
	modifyRequest.Domain = common.StringPtr(config.Domain.Domain)
	modifyRequest.RecordId = recordID
	modifyRequest.SubDomain = common.StringPtr(config.Domain.SubDomain)
	modifyRequest.RecordType = common.StringPtr("AAAA")
	modifyRequest.RecordLine = common.StringPtr("默认")
	modifyRequest.Value = common.StringPtr(ipv6)

	_, err = client.ModifyRecord(modifyRequest)
	return err
}

// UpdateDNSRecordWithRetry 添加重试机制的更新函数
func UpdateDNSRecordWithRetry(client *dnspod.Client, config config.Config, ipv6 string) error {
	operation := func() error {
		return UpdateDNSRecord(client, config, ipv6)
	}

	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.MaxElapsedTime = 5 * time.Minute

	return backoff.Retry(operation, backoffConfig)
}
