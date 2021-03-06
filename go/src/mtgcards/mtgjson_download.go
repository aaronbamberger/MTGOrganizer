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
    allPrintingsUrl = "AllPrintings"
    allPricesUrl = "AllPrices"
    versionUrl = "version"
)

const (
    jsonExt = ".json"
    bz2Ext = ".bz2"
    gzExt = ".gz"
)

const (
    downloadLocation = "/var/card-importer/card-data/"
    debugDownloadLocation = "./"
)

func DownloadAllPrintings(useCachedIfAvailable bool,
        useDebugDownloadLocation bool) (map[string]MTGSet, error) {
    result, err := downloadData(
        useCachedIfAvailable,
        useDebugDownloadLocation,
        allPrintingsUrl,
        decodeSets)
    if err != nil {
        return nil, err
    }

    if resultCast, ok := result.(map[string]MTGSet); !ok {
        return nil, fmt.Errorf("Unable to convert all printings result to correct type")
    } else {
        return resultCast, nil
    }
}

func DebugParseAllPrintingsGz(filename string) (map[string]MTGSet, error) {
    file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

    decompressor, err := gzip.NewReader(file)
    if err != nil {
        return nil, err
    }
    defer decompressor.Close()

    result, err := decodeSets(decompressor)
    if err != nil {
        return nil, err
    }

    if resultCast, ok := result.(map[string]MTGSet); !ok {
        return nil, fmt.Errorf("Unable to conver all printings result to correct type")
    } else {
        return resultCast, nil
    }
}

func DownloadAllPrices(useCachedIfAvailable bool) (map[string]MTGCardPrices, error) {
    result, err := downloadData(
        useCachedIfAvailable,
        false,
        allPricesUrl,
        decodePrices)
    if err != nil {
        return nil, err
    }

    if resultCast, ok := result.(map[string]MTGCardPrices); !ok {
        return nil, fmt.Errorf("Unable to convert all prices result to correct type")
    } else {
        return resultCast, nil
    }
}

func DownloadVersion() (MTGJSONVersion, error) {
    result, err := downloadData(
        false,
        false,
        versionUrl,
        decodeVersion)
    if err != nil {
        return MTGJSONVersion{}, err
    }

    if resultCast, ok := result.(MTGJSONVersion); !ok {
        return MTGJSONVersion{}, fmt.Errorf("Unable to convert version to correct type")
    } else {
        return resultCast, nil
    }
}

func decodeSets(input io.Reader) (interface{}, error) {
    decoder := json.NewDecoder(input)
    var result map[string]MTGSet
    if err := decoder.Decode(&result); err != nil {
        return nil, err
    }
    return result, nil
}

func decodePrices(input io.Reader) (interface{}, error) {
    decoder := json.NewDecoder(input)
    var result map[string]MTGCardPrices
    if err := decoder.Decode(&result); err != nil {
        return nil, err
    }
    return result, nil
}

func decodeVersion(input io.Reader) (interface{}, error) {
    decoder := json.NewDecoder(input)
    var result MTGJSONVersion
    if err := decoder.Decode(&result); err != nil {
        return nil, err
    }
    return result, nil
}

func downloadData(
        useCachedIfAvailable bool,
        useDebugDownloadLocation bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    result, err := tryGz(useCachedIfAvailable, useDebugDownloadLocation, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    result, err = tryBz2(useCachedIfAvailable, useDebugDownloadLocation, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    result, err = tryRaw(useCachedIfAvailable, useDebugDownloadLocation, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    return nil, fmt.Errorf("Unable to get %s from any sources", fileUrl)
}

func tryGz(
        useCachedIfAvailable bool,
        useDebugDownloadLocation bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt + gzExt, useCachedIfAvailable, useDebugDownloadLocation)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    decompressor, err := gzip.NewReader(reader)
    if err != nil {
        return nil, err
    }
    defer decompressor.Close()

    return decoderFn(decompressor)
}

func tryBz2(
        useCachedIfAvailable bool,
        useDebugDownloadLocation bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt + bz2Ext, useCachedIfAvailable, useDebugDownloadLocation)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    decompressor := bzip2.NewReader(reader)

    return decoderFn(decompressor)
}

func tryRaw(
        useCachedIfAvailable bool,
        useDebugDownloadLocation bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt, useCachedIfAvailable, useDebugDownloadLocation)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    return decoderFn(reader)
}

func tryDownload(filename string, useCachedIfAvailable bool, useDebugDownloadLocation bool) (io.ReadCloser, error) {
	// If we've either been asked to not use a local cached file, or
	// we have, but the file hasn't been downloaded, download the file
    var fileLocation string
    if useDebugDownloadLocation {
        fileLocation = debugDownloadLocation + filename
    } else {
        fileLocation = downloadLocation + filename
    }
	_, err := os.Stat(fileLocation)
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

		file, err := os.Create(fileLocation)
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
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}
	return file, err
}


