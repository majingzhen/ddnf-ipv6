package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go-intl-en/tencentcloud/dnspod/v20210323"

	"ddns-ipv6/config"
	"ddns-ipv6/dns"
	"ddns-ipv6/health"
	"ddns-ipv6/iputil"
	"ddns-ipv6/notification"
)

func checkIPv6Connectivity() bool {
	// 测试一些常用的 IPv6 地址
	hosts := []string{
		"2400:3200:baba::1",    // 阿里云 IPv6
		"2400:da00:2::29",      // 腾讯云 IPv6
		"2606:4700:4700::1111", // Cloudflare IPv6
	}

	for _, host := range hosts {
		conn, err := net.DialTimeout("tcp6", fmt.Sprintf("[%s]:443", host), 5*time.Second)
		if err == nil {
			conn.Close()
			return true
		}
	}
	return false
}

func main() {
	// 初始化组件
	cache := dns.NewDNSCache()
	healthCheck := health.NewHealthCheck()

	// 读取配置文件
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// 创建腾讯云客户端
	log.Println("Creating Tencent Cloud client...")
	credential := common.NewCredential(
		cfg.Tencent.SecretId,
		cfg.Tencent.SecretKey,
	)
	cpf := profile.NewClientProfile()
	client, err := dnspod.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		log.Fatalf("Failed to create DNSPod client: %v", err)
	}
	log.Println("Tencent Cloud client created successfully.")

	log.Printf("Starting IPv6 DDNS service...")

	// 检查IPv6连接
	if !checkIPv6Connectivity() {
		log.Println("IPv6 connectivity check failed, sending notification...")
		notification.SendNotification(cfg.Email,
			"IPv6 DDNS 更新失败",
			"无法连接到公共 IPv6 地址")
	}
	// 定期检查并更新IP
	for {
		log.Println("Checking local IPv6 address...")
		ipv6, err := iputil.GetLocalIPv6()
		if err != nil {
			log.Printf("Failed to get IPv6 address: %v", err)
			if healthCheck.RecordError() >= 3 {
				log.Println("Error threshold reached, sending notification...")
				notification.SendNotification(cfg.Email,
					"IPv6 DDNS 更新失败",
					fmt.Sprintf("获取IPv6地址失败: %v", err))
			}
			time.Sleep(time.Duration(cfg.CheckInterval) * time.Second)
			continue
		}

		log.Printf("Local IPv6 address: %s", ipv6)

		// 检查缓存，避免重复更新
		cachedIP, _ := cache.GetIP()
		if cachedIP == ipv6 {
			log.Printf("IP未变化，跳过更新")
			time.Sleep(time.Duration(cfg.CheckInterval) * time.Second)
			continue
		}

		log.Println("Updating DNS record...")
		// 使用重试机制更新DNS记录
		err = dns.UpdateDNSRecordWithRetry(client, *cfg, ipv6)
		if err != nil {
			log.Printf("Failed to update DNS record: %v", err)
			if healthCheck.RecordError() >= 3 {
				log.Println("Error threshold reached, sending notification...")
				notification.SendNotification(cfg.Email,
					"IPv6 DDNS 更新失败",
					fmt.Sprintf("更新DNS记录失败: %v", err))
			}
		} else {
			log.Printf("Successfully updated DNS record: %s.%s -> %s",
				cfg.Domain.SubDomain, cfg.Domain.Domain, ipv6)
			cache.UpdateIP(ipv6)
			healthCheck.RecordSuccess()
		}

		time.Sleep(time.Duration(cfg.CheckInterval) * time.Second)
	}
}
