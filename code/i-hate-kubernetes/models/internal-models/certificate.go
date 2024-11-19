package models

import (
	"strings"

	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/console"
	"github.com/ogdans3/i-hate-kubernetes/code/i-hate-kubernetes/models/util"
)

type CertificationService int

const (
	CERT_BOT CertificationService = iota
)

type CertificateHandler struct {
	Id         string
	Blocks     []CertificateBlock
	ServiceJob *CertificateJob
}

type CertificateBlock struct {
	Domains              []string
	Emails               []string
	CertificationService CertificationService
}

func ParseCertificateBlock(project Project, service *Service) *CertificateBlock {
	//TODO: Add email
	//TODO: Add command or type of certificate or service, only certbot for now
	return &CertificateBlock{
		Domains:              removeElements(service.Domain, "localhost", "127.0.0.1"),
		Emails:               nil,
		CertificationService: CERT_BOT,
	}
}

func ParseCertificateBlocks(project Project, services map[string]*Service) *CertificateHandler {
	blocks := make([]*CertificateBlock, 0)
	for _, service := range services {
		blocks = append(blocks, ParseCertificateBlock(project, service))
	}
	return &CertificateHandler{
		Id:         util.RandStringBytesMaskImpr(5),
		Blocks:     mergeBlocks(blocks),
		ServiceJob: ParseCertificateJob(project),
	}
}

func mergeBlocks(blocks []*CertificateBlock) []CertificateBlock {
	mergedBlocks := []CertificateBlock{}
	usedBlocks := map[*CertificateBlock]bool{} // To track used blocks

	for _, block1 := range blocks {
		// If block1 is already merged, skip it
		if usedBlocks[block1] {
			continue
		}

		mergedBlock := block1
		// Merge with other blocks if they have matching URLs
		for _, block2 := range blocks {
			if block1 == block2 || usedBlocks[block2] {
				continue
			}

			if anyUrlMatches(block1.Domains, block2.Domains) {
				mergedBlock.Domains = append(mergedBlock.Domains, block2.Domains...)
				mergedBlock.Emails = append(mergedBlock.Emails, block2.Emails...)
				usedBlocks[block2] = true
			}
		}

		// Remove duplicate domains and emails
		mergedBlock.Domains = uniqueElements(mergedBlock.Domains)
		mergedBlock.Emails = uniqueElements(mergedBlock.Emails)
		mergedBlocks = append(mergedBlocks, *mergedBlock)

		usedBlocks[block1] = true
	}
	return mergedBlocks
}

func getUrl(url string) string {
	parts := strings.Split(url, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1]

	}
	console.InfoLog.Fatal("Url does not have atleast 2 parts. Is the url correct?", url)
	panic("Should never get here")
}

func anyUrlMatches(urls1 []string, urls2 []string) bool {
	for _, url1 := range urls1 {
		for _, url2 := range urls2 {
			if getUrl(url1) == getUrl(url2) {
				return true
			}
		}
	}
	return false
}

func uniqueElements(slice []string) []string {
	// Create a map to store unique elements
	seen := make(map[string]bool)
	var result []string

	// Loop through the slice
	for _, item := range slice {
		if !seen[item] {
			// If the item is not in the map, add it to the result and mark it as seen
			result = append(result, item)
			seen[item] = true
		}
	}

	return result
}

func removeElements(slice []string, elements ...string) []string {
	// Create a map to track elements to remove for faster lookup
	toRemove := make(map[string]bool)
	for _, e := range elements {
		toRemove[e] = true
	}

	// Create a new slice to hold the result
	var result []string

	// Loop through the original slice
	for _, item := range slice {
		if !toRemove[item] {
			// Append item if it's not in the toRemove map
			result = append(result, item)
		}
	}

	return result
}
