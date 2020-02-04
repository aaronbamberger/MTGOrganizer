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
)

const (
    jsonExt = ".json"
    bz2Ext = ".bz2"
    gzExt = ".gz"
)

func DownloadAllPrintings(useCachedIfAvailable bool) (map[string]MTGSet, error) {
    result, err := downloadData(
        useCachedIfAvailable,
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

func DownloadAllPrices(useCachedIfAvailable bool) (map[string]MTGCardPrices, error) {
    result, err := downloadData(
        useCachedIfAvailable,
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

func downloadData(
        useCachedIfAvailable bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    result, err := tryGz(useCachedIfAvailable, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    result, err = tryBz2(useCachedIfAvailable, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    result, err = tryRaw(useCachedIfAvailable, fileUrl, decoderFn)
    if err == nil {
        return result, nil
    }
    log.Print(err)

    return nil, fmt.Errorf("Unable to get %s from any sources", fileUrl)
}

func tryGz(
        useCachedIfAvailable bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt + gzExt, useCachedIfAvailable)
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
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt + bz2Ext, useCachedIfAvailable)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    decompressor := bzip2.NewReader(reader)

    return decoderFn(decompressor)
}

func tryRaw(
        useCachedIfAvailable bool,
        fileUrl string,
        decoderFn func(io.Reader) (interface{}, error)) (interface{}, error) {
    reader, err := tryDownload(fileUrl + jsonExt, useCachedIfAvailable)
    if err != nil {
        return nil, err
    }
    defer reader.Close()

    return decoderFn(reader)
}

func tryDownload(filename string, useCachedIfAvailable bool) (io.ReadCloser, error) {
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
