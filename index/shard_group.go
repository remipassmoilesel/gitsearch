package index

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/mapping"
	"github.com/pkg/errors"
	"os"
	"path"
	"strconv"
)

type ShardGroup struct {
	shards       []bleve.Index
	shardNumber  int
	dataRootPath string
	searchIndex  bleve.IndexAlias
}

func NewShardGroup(dataRootPath string, shardNumber int) ShardGroup {
	return ShardGroup{
		shardNumber:  shardNumber,
		dataRootPath: dataRootPath,
	}
}

type GroupInitResult struct {
	number int
}

type shardInitResult struct {
	shard *bleve.Index
	err   error
}

func (s *ShardGroup) Initialize() (GroupInitResult, error) {
	if len(s.shards) != 0 {
		return GroupInitResult{}, errors.New("shard array must be empty")
	}

	s.searchIndex = bleve.NewIndexAlias()

	ch := make(chan shardInitResult)
	for i := 0; i < s.shardNumber; i++ {
		shardPath := path.Join(s.dataRootPath, strconv.Itoa(i))
		go s.initializeShard(ch, shardPath)
	}

	for res := range ch {
		if res.err != nil {
			return GroupInitResult{}, errors.Wrap(res.err, "shard init failed")
		}
		shard := *res.shard
		s.shards = append(s.shards, shard)
		s.searchIndex.Add(shard)

		if len(s.shards) == s.shardNumber {
			break
		}
	}

	return GroupInitResult{number: len(s.shards)}, nil
}

func (s *ShardGroup) initializeShard(ch chan shardInitResult, shardPath string) {
	var err error
	var shard bleve.Index
	if _, ferr := os.Stat(shardPath); ferr == nil {
		shard, err = bleve.Open(shardPath)
	} else {
		shard, err = bleve.New(shardPath, indexMapping())
	}

	result := shardInitResult{
		shard: &shard,
		err:   err,
	}

	ch <- result
}

func (s *ShardGroup) GetShard(id int) *bleve.Index {
	return &s.shards[id]
}

func (s *ShardGroup) Size() int {
	return len(s.shards)
}

func (s *ShardGroup) Clean() error {
	var err error
	if _, ferr := os.Stat(s.dataRootPath); ferr == nil {
		err = os.RemoveAll(s.dataRootPath)
	}
	if err != nil {
		return errors.Wrap(err, "cannot clean shards")
	}

	s.shards = []bleve.Index{}
	s.searchIndex = bleve.NewIndexAlias()

	return nil
}

func (s *ShardGroup) Close() error {
	var err error
	for i, shard := range s.shards {
		err = shard.Close()
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot close shard %v", i))
		}
	}
	s.shards = []bleve.Index{}
	return nil
}

func indexMapping() *mapping.IndexMappingImpl {
	notIndexed := bleve.NewTextFieldMapping()
	notIndexed.Store = true
	notIndexed.IncludeInAll = true
	notIndexed.IncludeTermVectors = false
	notIndexed.Analyzer = "standard" // TODO: register then use "simple"

	indexed := bleve.NewTextFieldMapping()
	indexed.Store = true
	indexed.IncludeInAll = true
	indexed.IncludeTermVectors = true
	indexed.Analyzer = "standard"

	docMapping := bleve.NewDocumentMapping()

	docMapping.AddFieldMappingsAt("Hash", notIndexed)
	docMapping.AddFieldMappingsAt("Commit", notIndexed)
	docMapping.AddFieldMappingsAt("Content", indexed)
	docMapping.AddFieldMappingsAt("Path", indexed)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = docMapping
	indexMapping.DefaultAnalyzer = "standard"
	return indexMapping
}
