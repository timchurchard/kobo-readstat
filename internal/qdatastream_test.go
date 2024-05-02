package internal

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	t.Run("run example as test for coverage", func(t *testing.T) {
		// An EventData blob from the Events table in the Kobo eReader firmware.
		b, err := hex.DecodeString("0000000400000010005600690065007700540079007000650000000a00000000060054004f00430000003000450078007400720061004400610074006100520065006100640069006e006700530065007300730069006f006e00730000000200000000030000002e00450078007400720061004400610074006100520065006100640069006e0067005300650063006f006e0064007300000002000000000a00000028004500780074007200610044006100740061004400610074006500430072006500610074006500640000000a00000000280032003000310039002d00310031002d00320035005400300031003a00310037003a00310034005a")
		if err != nil {
			panic(err)
		}

		r := bytes.NewBuffer(b)
		v, err := (&QDataStreamReader{
			Reader:    r,
			ByteOrder: binary.BigEndian,
		}).ReadQStringQVariantAssociative()
		if err != nil {
			panic(err)
		}

		err = json.NewEncoder(os.Stdout).Encode(v)
		assert.NoError(t, err)

		if r.Len() != 0 {
			panic("not all read")
		}
	})
}
