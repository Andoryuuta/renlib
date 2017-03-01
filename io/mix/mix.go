package mix

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
)

type FileInfo struct {
	Name string
	ID   uint32 // ID is actually a CRC32 of strings.ToUpper(filename)
	Data []byte
}

type MixFile struct {
	Files []FileInfo
}

const MIX_FILE_MAGIC string = "MIX1"

func ParseFile(filename string) (*MixFile, error) {
	log.Printf("Opening file %s\n", filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return ParseFromReader(f)
}

func ParseBytes(data []byte) (*MixFile, error) {
	r := bytes.NewReader(data)
	return ParseFromReader(r)
}

func ParseFromReader(f ReadAtSeeker) (*MixFile, error) {
	// Read in file magic and verify.
	// It should be equal to "MIX1"
	magic, err := ReadFixedUTF8String(f, 4)
	if err != nil {
		return nil, err
	}

	if magic != MIX_FILE_MAGIC {
		return nil, errors.New(fmt.Sprintf("Invalid file magic. Expected: (\"%s\" aka [% X]) Got: (\"%s\" aka [% X])", MIX_FILE_MAGIC, []byte(MIX_FILE_MAGIC), magic, []byte(magic)))
	}

	// Get the info table offset
	var infoTableOffset uint32
	err = binary.Read(f, binary.LittleEndian, &infoTableOffset)
	if err != nil {
		return nil, err
	}

	// Get the name table offset
	var nameTableOffset uint32
	err = binary.Read(f, binary.LittleEndian, &nameTableOffset)
	if err != nil {
		return nil, err
	}

	// Read in the file name table
	fileNameMap, err := readNameTable(f, nameTableOffset)
	if err != nil {
		return nil, err
	}

	// Go to the info name table
	f.Seek(int64(infoTableOffset), 0)

	// Read in file count
	var count uint32
	err = binary.Read(f, binary.LittleEndian, &count)
	if err != nil {
		return nil, err
	}

	output := &MixFile{}

	// Read in the file(s) information (id, offset, size)
	for i := 0; i < int(count); i++ {
		// Read in the file info
		var id, offset, size uint32

		err = binary.Read(f, binary.LittleEndian, &id)
		if err != nil {
			return nil, err
		}
		err = binary.Read(f, binary.LittleEndian, &offset)
		if err != nil {
			return nil, err
		}
		err = binary.Read(f, binary.LittleEndian, &size)
		if err != nil {
			return nil, err
		}

		// Read in the file data
		data := make([]byte, size)
		f.ReadAt(data, int64(offset))

		// Append file to files list
		output.Files = append(output.Files, FileInfo{
			Name: fileNameMap[i],
			ID:   id,
			Data: data,
		})
	}

	return output, nil
}

func readNameTable(f ReadAtSeeker, offset uint32) (map[int]string, error) {
	// Save original file offset
	originalOffset, err := f.Seek(0, 1)
	if err != nil {
		return nil, err
	}

	// Go to the file name table
	f.Seek(int64(offset), 0)

	// Read in filename entry count
	var count uint32
	err = binary.Read(f, binary.LittleEndian, &count)
	if err != nil {
		return nil, err
	}

	// Create a map of file ID->file name
	fileNameMap := make(map[int]string)

	// Read in file names
	for i := 0; i < int(count); i++ {

		// The file name is length prefixed
		var nameSize uint8
		err = binary.Read(f, binary.LittleEndian, &nameSize)
		if err != nil {
			return nil, err
		}

		// Read in file name
		fileName, err := ReadFixedUTF8String(f, int(nameSize))
		if err != nil {
			return nil, err
		}

		// The string are null terminated, so strip off the null byte
		fileName = fileName[:len(fileName)-1]

		// Add it to the map
		fileNameMap[i] = fileName
	}

	// Restore original file offset
	f.Seek(originalOffset, 0)

	return fileNameMap, nil
}
