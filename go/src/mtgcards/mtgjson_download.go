package mtgcards

import "compress/bzip2"
import "compress/gzip"
import "encoding/json"
import "fmt"
import "io"
import "log"
import "net/http"
import "os"

const (
	mtgjsonBaseUrl = "https://www.mtgjson.com/files/"
	allPrintingsRawJson = "AllPrintings.json"
	allPrintingsBz2 = "AllPrintings.json.bz2"
	allPrintingsGz = "AllPrintings.json.gz"
)

func DownloadAllPrintings(useCachedIfAvailable bool) (map[string]MTGSet, error) {
	// Go down the priority list of file types (gz, bz2, zip, raw),
	// optionally trying to load them from files already downloaded
	// to disk

	// Try gzip
	allSets, err := TryGz(useCachedIfAvailable)
	if err == nil {
		return allSets, nil
	}
	log.Print(err)

	// Try bz2
	allSets, err = TryBz2(useCachedIfAvailable)
	if err == nil {
		return allSets, nil
	}
	log.Print(err)

	// Try raw json
	allSets, err = TryRaw(useCachedIfAvailable)
	if err == nil {
		return allSets, nil
	}
	log.Print(err)

	return nil, fmt.Errorf("Unable to get card info from any sources")
}

func TryGz(useCachedIfAvailable bool) (map[string]MTGSet, error) {
	reader, err := TryDownload(allPrintingsGz, useCachedIfAvailable)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressor, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer decompressor.Close()

	decoder := json.NewDecoder(decompressor)
	var allSets map[string]MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		return nil, err
	}
	return allSets, nil
}

func TryBz2(useCachedIfAvailable bool) (map[string]MTGSet, error) {
	reader, err := TryDownload(allPrintingsBz2, useCachedIfAvailable)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressor := bzip2.NewReader(reader)

	decoder := json.NewDecoder(decompressor)
	var allSets map[string]MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		return nil, err
	}
	return allSets, nil
}

func TryRaw(useCachedIfAvailable bool) (map[string]MTGSet, error) {
	reader, err := TryDownload(allPrintingsRawJson, useCachedIfAvailable)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	var allSets map[string]MTGSet
	if err := decoder.Decode(&allSets); err != nil {
		return nil, err
	}
	return allSets, nil
}

func TryDownload(filename string, useCachedIfAvailable bool) (io.ReadCloser, error) {
	// If we've either been asked to not use a local cached file, or
	// we have, but the file hasn't been downloaded, download the file
	_, err := os.Stat(filename)
	if !useCachedIfAvailable || os.IsNotExist(err) {
		fullUrl := mtgjsonBaseUrl + filename
		resp, err := http.Get(fullUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Error %s while fetching %s", resp.Status, fullUrl)
		}

		file, err := os.Create(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return nil, err
		}
	}

	// If we're here, we've either freshly downloaded the file, or have determined
	// there's an existing cached version we can use
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return file, err
}
