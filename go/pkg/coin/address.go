package coin

// Placeholder identifier types. Real implementation would use cryptographic hashes.
type IDKey [20]byte
type IDScript [20]byte

type None struct{}

// DestinationTx is a simplified union of possible address types.
type DestinationTx interface{}

const (
	TypePubKey     = 71
	TypeScript     = 8
	TypePubKeyTest = 111
	TypeScriptTest = 196
)

// Address represents a base58-encoded address.
type Address struct {
	Base58
}

// SetString decodes a base58check string into the address.
func (a *Address) SetString(val string) bool {
	return a.Base58.SetString(val)
}

// String returns the base58check encoding of the address including the version byte.
func (a Address) String() string {
	return a.Base58.ToString(true)
}

func (a *Address) SetIDKey(id IDKey) bool {
	if TestNet {
		a.SetData(TypePubKeyTest, id[:])
	} else {
		a.SetData(TypePubKey, id[:])
	}
	return true
}

func (a *Address) SetIDScript(id IDScript) bool {
	if TestNet {
		a.SetData(TypeScriptTest, id[:])
	} else {
		a.SetData(TypeScript, id[:])
	}
	return true
}

func (a *Address) SetDestinationTx(v DestinationTx) bool {
	switch val := v.(type) {
	case IDKey:
		return a.SetIDKey(val)
	case IDScript:
		return a.SetIDScript(val)
	default:
		return false
	}
}

// IsValid performs basic validation of the address contents.
func (a Address) IsValid() bool {
	switch a.Version {
	case TypePubKey, TypePubKeyTest, TypeScript, TypeScriptTest:
		return len(a.Data) == 20
	default:
		return false
	}
}

// Get returns the destination represented by the address.
func (a Address) Get() DestinationTx {
	if !a.IsValid() {
		return None{}
	}
	var id [20]byte
	copy(id[:], a.Data)
	switch a.Version {
	case TypePubKey, TypePubKeyTest:
		return IDKey(id)
	case TypeScript, TypeScriptTest:
		return IDScript(id)
	default:
		return None{}
	}
}

// GetIDKey extracts the IDKey if present.
func (a Address) GetIDKey() (IDKey, bool) {
	if a.Version == TypePubKey || a.Version == TypePubKeyTest {
		var id IDKey
		copy(id[:], a.Data)
		return id, true
	}
	return IDKey{}, false
}

// IsScript indicates whether the address encodes a script hash.
func (a Address) IsScript() bool {
	return a.Version == TypeScript || a.Version == TypeScriptTest
}
