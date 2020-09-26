package queuedefs

import "babblegraph/util/queue"

type queueTopicName string

func (q queueTopicName) Str() string {
	return string(q)
}

func RegisterQueues(errs chan error) error {
	linkHandlerQueue, err := initializeLinkHandlerQueue(errs)
	if err != nil {
		return err
	}
	queues := []queue.Queue{
		fetchQueue{},
		parseQueue{},
		normalizeTextQueue{},
		linkHandlerQueue,
		lemmatizeQueue{},
		indexQueue{},
	}
	for _, q := range queues {
		if err := queue.RegisterQueue(q); err != nil {
			return err
		}
	}
	queue.StartQueue(errs)
	return publishMessageToFetchQueue("https://cnnespanol.cnn.com/2020/08/29/la-lucha-de-europa-contra-el-covid-19-pasa-de-los-hospitales-a-las-calles/")
}
