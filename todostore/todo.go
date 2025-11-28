package todostore

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"todo/filestorage"
	"todo/logger"
)

type TodoList struct {
	Name   string         `json:"Name"`
	LItems []TodoListItem `json:"lItems"`
}

type TodoListMutex struct {
	mutex sync.Mutex
}

type TodoListItem struct {
	Name        string
	Description string
	Status      string
}

var lists map[string]*TodoList = make(map[string]*TodoList)
var lists_mutex map[string]*TodoListMutex = make(map[string]*TodoListMutex)

var mutex sync.Mutex

const (
	StatusNotStarted string = "not started"
	StatusStarted    string = "started"
	StatusCompleted  string = "completed"
)

func ReadFromMap(ctx context.Context, todoListName string) (*TodoList, *TodoListMutex, error) {
	mutex.Lock()
	defer mutex.Unlock()
	list, ok := lists[todoListName]
	mutex := lists_mutex[todoListName]
	if ok {
		return list, mutex, nil
	} else {
		list, err := retrieveListFromFile(ctx, todoListName)
		if err != nil {
			return nil, nil, err
		}

		if list != nil {
			logger.GetCtxLogger(ctx).Info("Retrieved list from file: " + todoListName)
			lists[todoListName] = list
		} else {
			logger.GetCtxLogger(ctx).Info("Creating todo list and adding to todostore: " + todoListName)
			lists[todoListName] = &TodoList{Name: todoListName}
		}

		lists_mutex[todoListName] = &TodoListMutex{}
		return lists[todoListName], lists_mutex[todoListName], nil
	}
}

func GetList(ctx context.Context, todoListName string) (*TodoList, error) {
	list, mutex, err := ReadFromMap(ctx, todoListName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return nil, err
	}
	if list == nil {
		logger.GetCtxLogger(ctx).Info("Init TodoStore for todolist: " + todoListName)
		var err error
		list, err = retrieveListFromFile(ctx, todoListName)
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}

func CreateList(ctx context.Context, todoListName string) (*TodoList, error) {
	list, mutex, err := ReadFromMap(ctx, todoListName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func AddItemToList(ctx context.Context, listName string, itemName string, itemDescription string) error {
	list, mutex, err := ReadFromMap(ctx, listName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return err
	}
	var alreadyExists bool = false
	for _, lItem := range list.LItems {
		if lItem.Name == itemName {
			alreadyExists = true
			break
		}
	}

	if alreadyExists {
		return errors.New("Item already exists: " + itemName)
	}

	lItem := TodoListItem{Name: itemName, Description: itemDescription, Status: StatusNotStarted}
	list.LItems = append(list.LItems, lItem)
	logger.GetCtxLogger(ctx).Info("Added item: " + lItem.Name + " to List: " + list.Name)

	return nil
}

func UpdateListItemDescription(ctx context.Context, listName string, itemName string, itemDescription string) error {
	list, mutex, err := ReadFromMap(ctx, listName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return err
	}

	var itemFound bool = false
	var updateItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems[updateItemIndex].Description = itemDescription
		logger.GetCtxLogger(ctx).Info("Item Updated (Description): " + itemName + " in List: " + list.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func UpdateListItemStatus(ctx context.Context, listName string, itemName string, itemStatus string) error {
	list, mutex, err := ReadFromMap(ctx, listName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return err
	}

	if itemStatus != StatusNotStarted && itemStatus != StatusStarted && itemStatus != StatusCompleted {
		err := errors.New("Invalid status provided: " + itemStatus)
		return err
	}

	var itemFound bool = false
	var updateItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			updateItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems[updateItemIndex].Status = itemStatus
		logger.GetCtxLogger(ctx).Info("Item Updated (Status): " + itemName + " in List: " + list.Name)
	} else {
		return errors.New("Cannot find Item to update: " + itemName)
	}

	return nil
}

func DeleteItemFromList(ctx context.Context, listName string, itemName string) error {
	list, mutex, err := ReadFromMap(ctx, listName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return err
	}

	var itemFound bool = false
	var deleteItemIndex int
	for index, lItem := range list.LItems {
		if lItem.Name == itemName {
			deleteItemIndex = index
			itemFound = true
			break
		}
	}
	if itemFound {
		list.LItems = append(list.LItems[:deleteItemIndex], list.LItems[deleteItemIndex+1:]...)
		logger.GetCtxLogger(ctx).Info("Item Deleted: " + itemName + " from List: " + list.Name)
	} else {
		logger.GetCtxLogger(ctx).Info("Cannot find Item to delete: " + itemName)
	}

	return nil
}

func retrieveListFromFile(ctx context.Context, listName string) (*TodoList, error) {
	list_b, err := filestorage.LoadFileToByteSlice(ctx, listName+".json")
	if err != nil || list_b == nil {
		return nil, err
	}

	var list TodoList

	if list_b != nil {
		logger.GetCtxLogger(ctx).Info("Getting todo list: " + listName)
		logger.GetCtxLogger(ctx).Info(string(list_b))
		err := json.Unmarshal(list_b, &list)
		if err != nil {
			return nil, err
		}
	}

	return &list, nil
}

func SaveList(ctx context.Context, listName string) error {
	list, mutex, err := ReadFromMap(ctx, listName)
	mutex.mutex.Lock()
	defer mutex.mutex.Unlock()
	if err != nil {
		return err
	}
	list_bb, err := json.Marshal(list)
	if err != nil {
		return err
	}

	filestorage.SaveByteSliceToFile(list_bb, listName+".json")

	logger.GetCtxLogger(ctx).Info(string(list_bb))

	return nil
}
