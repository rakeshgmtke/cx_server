// hss_data.go
package main

import (
	"fmt"
	"sync"

	"github.com/rakeshgmtke/go-diameter/diam"
    "github.com/rakeshgmtke/go-diameter/diam/datatype"
	_"github.com/rakeshgmtke/go-diameter/diam/avp"	
    _ "github.com/rakeshgmtke/go-diameter/diam/dict"

)

// To store HSS related data
// HSSData struct to hold the data.
type HSSData struct {
	SCSCFName          string
	IMPI               string
	IMSRestorationInfo *diam.AVP
}

// A global map to store the data.
var HSSDataMap map[string]HSSData
var mapMutex sync.RWMutex // Use a mutex for concurrent access


// populateHSS populates the HSS data map with initial data.
func populateHSS(size int) {
	for i := 0; i < size; i++ {
		impu := fmt.Sprintf("IMPU-%d", i)
		imsRestorationInfo := diam.NewAVP(0, 0, 0, datatype.Unsigned32(i))
		data := HSSData{
			IMPI:               string(i),
			SCSCFName:          string(i),
			IMSRestorationInfo: imsRestorationInfo,
		}
		HSSDataMap[impu] = data // No mutex needed here; single-threaded during initialization.
	}
}

// addOrModifyIMPI adds or modifies the IMPI field of an HSSData entry.
func addOrModifyIMPI(impu string, impi string) bool {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	existingData, exists := HSSDataMap[impu]
	if exists {
		// Modify existing data
		existingData.IMPI = impi
		HSSDataMap[impu] = existingData // Update the map
		//log.Printf("Modified IMPI for IMPU '%s' to '%s'\n", impu, impi)
		return true
	} else {
		// Add new data
		HSSDataMap[impu] = HSSData{IMPI: impi} // Only set the IMPI
		//log.Printf("Added IMPI for IMPU '%s' to '%s'\n", impu, impi)
		return true
	}
	return false //Unreachable
}

// addOrModifySCSCFName adds or modifies the SCSCFName field of an HSSData entry.
func addOrModifySCSCFName(impu string, scscfName string) bool {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	existingData, exists := HSSDataMap[impu]
	if exists {
		// Modify existing data
		existingData.SCSCFName = scscfName
		HSSDataMap[impu] = existingData // Update the map
		//log.Printf("Modified SCSCFName for IMPU '%s' to '%s'\n", impu, scscfName)
		return true
	} else {
		// Add new data
		HSSDataMap[impu] = HSSData{SCSCFName: scscfName} // Only set SCSCFName
		//log.Printf("Added SCSCFName for IMPU '%s' to '%s'\n", impu, scscfName)
		return true
	}
	return false //Unreachable
}

// addOrModifyIMSRestorationInfo adds or modifies the IMSRestorationInfo field of an HSSData entry.
func addOrModifyIMSRestorationInfo(impu string, RestorationInfo *diam.AVP) bool {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	existingData, exists := HSSDataMap[impu]
	if exists {
		// Modify existing data
		existingData.IMSRestorationInfo = RestorationInfo
		HSSDataMap[impu] = existingData // Update the map
		//log.Printf("Modified IMSRestorationInfo for IMPU '%s' to '%s'\n", impu, RestorationInfo)
		return true
	} else {
		// Add new data
		HSSDataMap[impu] = HSSData{IMSRestorationInfo: RestorationInfo} // Only set RestorationInfo
		//log.Printf("Added IMSRestorationInfo for IMPU '%s' to '%s'\n", impu, RestorationInfo)
		return true
	}
	return false //Unreachable.
}

// deleteIMPUData deletes an entry from the HSS data map.
func deleteIMPUData(impu string) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	_, exists := HSSDataMap[impu]
	if !exists {
		//log.Printf("Error: IMPU '%s' not found\n", impu)
		return // IMPORTANT: Return on error!
	}
	delete(HSSDataMap, impu)
	//log.Printf("Deleted IMPU '%s'\n", impu)
}

// readIMPUData reads an entry from the HSS data map.
func readIMPUData(impu string) (string, string, *diam.AVP) {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	data, exists := HSSDataMap[impu]
	if !exists {
		//log.Printf("Error: IMPU '%s' not found\n", impu)
		return "", "", nil
	}

	return data.IMPI, data.SCSCFName, data.IMSRestorationInfo
}

// readIMPIData reads IMPI data for a given IMPU.
func readIMPIData(impu string) string {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	data, exists := HSSDataMap[impu]
	if !exists {
		//log.Printf("Error: IMPU '%s' not found\n", impu)
		return ""
	}
	return data.IMPI
}

// readSCSCFNameData reads SCSCFName data for a given IMPU.
func readSCSCFNameData(impu string) string {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	data, exists := HSSDataMap[impu]
	if !exists {
		//log.Printf("Error: IMPU '%s' not found\n", impu)
		return ""
	}
	return data.SCSCFName
}

// readIMSRestorationInfoData reads IMSRestorationInfo data for a given IMPU.
func readIMSRestorationInfoData(impu string) *diam.AVP {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	data, exists := HSSDataMap[impu]
	if !exists {
		//log.Printf("Error: IMPU '%s' not found\n", impu)
		return nil
	}
	return data.IMSRestorationInfo
}