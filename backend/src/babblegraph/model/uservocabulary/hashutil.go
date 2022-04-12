package uservocabulary

import (
	"babblegraph/util/ptr"
	"babblegraph/wordsmith"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

type uniqueHashable interface {
	getType() VocabularyType
	getHashableComponent() string
	getDisplay() string
	getID() *string
}

type HashableLemma struct {
	LemmaID   wordsmith.LemmaID `json:"lemma_id"`
	LemmaText string            `json:"lemma_text"`
}

func (h *HashableLemma) getType() VocabularyType {
	return VocabularyTypeLemma
}

func (h *HashableLemma) getHashableComponent() string {
	return string(h.LemmaID)
}

func (h *HashableLemma) getDisplay() string {
	return h.LemmaText
}

func (h *HashableLemma) getID() *string {
	return ptr.String(h.LemmaID.Str())
}

type HashablePhrase struct {
	DisplayText  string
	DefinitionID *wordsmith.PhraseDefinitionID
}

func (h *HashablePhrase) getType() VocabularyType {
	return VocabularyTypePhrase
}

func (h *HashablePhrase) getHashableComponent() string {
	if h.DefinitionID != nil {
		return string(*h.DefinitionID)
	}
	return strings.Join(strings.Split(h.DisplayText, " "), "_")
}

func (h *HashablePhrase) getDisplay() string {
	return h.DisplayText
}

func (h *HashablePhrase) getID() *string {
	if h.DefinitionID != nil {
		return ptr.String(string(*h.DefinitionID))
	}
	return nil
}

func GetUniqueHash(u uniqueHashable) UniqueHash {
	hashInput := fmt.Sprintf("%s_%s", u.getType(), u.getHashableComponent())
	md5Hash := md5.Sum([]byte(hashInput))
	return UniqueHash(hex.EncodeToString(md5Hash[:]))
}
